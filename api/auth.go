package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/twitch"
	"go.mongodb.org/mongo-driver/bson"
)

// HandleLogin is a handler that redirects the user to Twitch for login, and provides the 'state'
// parameter which protects against login CSRF
func (a *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("Corrupted session %s -- generated new", err)
		err = nil
	}

	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		RespondWithError(w, r, http.StatusBadRequest, "Couldn't generate a session!")
	}

	state := hex.EncodeToString(tokenBytes[:])

	session.AddFlash(state, stateCallbackKey)

	if err = session.Save(r, w); err != nil {
		return
	}

	http.Redirect(w, r, oauth2Config.AuthCodeURL(state, claims), http.StatusTemporaryRedirect)
}

// HandleOauth2Callback is a Handler for oauth's 'redirect_uri' endpoint;
// it validates the state token and retrieves an OAuth token from the request parameters.
func (a *App) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	// ensure we flush the csrf challenge even if the request is ultimately unsuccessful
	defer func() {
		if err := session.Save(r, w); err != nil {
			log.Printf("error saving session: %s", err)
		}
	}()

	switch stateChallenge, state := session.Flashes(stateCallbackKey), r.FormValue("state"); {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge[0]:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'", state, stateChallenge[0])
	}

	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, "Couldn't verify your confirmation, please try again.")
	}

	token, err := oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return
	}

	// add the oauth token to session
	session.Values[oauthTokenKey] = token

	fmt.Printf("Access token: %s\n", token.AccessToken)

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		fmt.Println("can't extract id token from access token")
	}

	idToken, err := oidcVerifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, "Couldn't verify your confirmation, please try again")
	}

	var claims struct {
		Iss   string `json:"iss"`
		Sub   string `json:"sub"`
		Aud   string `json:"aud"`
		Exp   int32  `json:"exp"`
		Iat   int32  `json:"iat"`
		Nonce string `json:"nonce"`
		Email string `json:"email"`
	}

	if err := idToken.Claims(&claims); err != nil {
		RespondWithError(w, r, http.StatusBadRequest, fmt.Sprintf("Couldn't verify your confirmation, please try again: %v", err))
	}

	// Fetch Twitch data by access token
	userData, err := twitch.GetTwitchUserByToken(token.AccessToken)
	if err != nil {
		RespondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	ud := userData.Users[0]

	toStore := twitch.AuthUser{
		Login:        ud.Login,
		ID:           ud.ID,
		Email:        claims.Email,
		Scopes:       scopes,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	// Checks if auth user exists or not
	var search twitch.AuthUser
	if err := database.AuthUsers.FindOne(context.Background(), bson.M{"id": toStore.ID}).Decode(&search); err != nil {
		if err.Error() == "mongo: no documents in result" {
			_, errStore := database.AuthUsers.InsertOne(context.Background(), toStore)
			if errStore != nil {
				RespondWithError(w, r, http.StatusInternalServerError, "Error storing user")
				return
			}
		} else {
			_, errUpdate := database.AuthUsers.UpdateOne(
				context.Background(),
				bson.M{"id": toStore.ID},
				bson.M{"$set": bson.M{"access_token": token.AccessToken, "refresh_token": token.RefreshToken}},
			)
			if errUpdate != nil {
				RespondWithError(w, r, http.StatusInternalServerError, "Error updating tokens")
				return
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
