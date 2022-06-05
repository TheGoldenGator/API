package api

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mahcks/TheGoldenGator/configure"
	"github.com/Mahcks/TheGoldenGator/websocket"
	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	Router *mux.Router
}

type ResponseBody struct {
	StatusCode int         `json:"status_code"`
	Timestamp  int         `json:"timestamp"`
	Data       interface{} `json:"data"`
}

type ResponseBodyError struct {
	StatusCode int    `json:"status_code"`
	Timestamp  int    `json:"timestamp"`
	Error      string `json:"error"`
}

const (
	stateCallbackKey = "oauth-state-callback"
	oauthSessionName = "oauth-oidc-session"
	oauthTokenKey    = "oauth-token"
)

var (
	// Consider storing the secret in an environment variable or a dedicated storage system.
	scopes       = []string{"user:read:email", "openid", "channel:read:redemptions"}
	claims       = oauth2.SetAuthURLParam("claims", `{"id_token":{"email":null}}`)
	oauth2Config *oauth2.Config
	oidcIssuer   = "https://id.twitch.tv/oauth2"
	oidcVerifier *oidc.IDTokenVerifier
	cookieSecret = []byte("Please use a more sensible secret than this one")
	cookieStore  = sessions.NewCookieStore(cookieSecret)
)

func (a *App) Initialize() {
	a.Router = mux.NewRouter()

	/* logger := log.New(os.Stdout, "", log.LstdFlags)
	logMiddleware := middleware.NewLogMiddleware(logger)
	a.Router.Use(logMiddleware.Func()) */

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	port := configure.Config.GetString("port")
	if port == "" {
		port = "8000"
	}

	fmt.Println("[INFO] Started running on http://localhost:" + port)
	//log.Fatal(http.ListenAndServe(":8000", a.Router))
	log.Fatal(
		http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}

func RespondWithError(w http.ResponseWriter, r *http.Request, code int, message string) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	responseBody := &ResponseBodyError{
		StatusCode: code,
		Timestamp:  int(time.Now().Unix()),
		Error:      message,
	}

	buf, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}

	w.Write(buf)
	return
}

func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	responseBody := &ResponseBody{
		StatusCode: code,
		Timestamp:  int(time.Now().Unix()),
		Data:       payload,
	}

	buf, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}

	w.Write(buf)
	return
}

func (a *App) initializeRoutes() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})

	/* Auth */
	provider, err := oidc.NewProvider(context.Background(), oidcIssuer)
	if err != nil {
		log.Fatal(err)
	}

	oauth2Config = &oauth2.Config{
		ClientID:     configure.Config.GetString("twitch_client_id"),
		ClientSecret: configure.Config.GetString("twitch_client_secret"),
		Scopes:       scopes,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8000/redirect",
	}
	oidcVerifier = provider.Verifier(&oidc.Config{ClientID: configure.Config.GetString("twitch_client_id")})

	a.Router.NotFoundHandler = http.HandlerFunc(a.NotFound)

	pool := websocket.WSPool
	go pool.Start()

	/* WebSockets */
	a.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		a.serveWs(pool, w, r)
	})

	/* REST ROUTES */
	a.Router.HandleFunc("/", a.Home).Methods("GET")
	a.Router.HandleFunc("/streams", a.Streams).Methods("GET")
	a.Router.HandleFunc("/streamers", a.Users).Methods("GET")
	a.Router.HandleFunc("/eventsub", a.EventsubRecievedNotification).Methods("POST")
	a.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	/* Auth */
	a.Router.HandleFunc("/login", a.HandleLogin).Methods("GET")
	a.Router.HandleFunc("/redirect", a.HandleOAuth2Callback).Methods("GET")

	a.Router.HandleFunc("/test", a.Test).Methods("GET")
	/* a.Router.HandleFunc("/teams", a.TeamData).Methods("GET") */
}
