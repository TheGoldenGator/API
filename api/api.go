package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Mahcks/TheGoldenGator/configure"

	"github.com/Mahcks/TheGoldenGator/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

func (a *App) Initialize() {
	a.Router = mux.NewRouter()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logMiddleware := middleware.NewLogMiddleware(logger)
	a.Router.Use(logMiddleware.Func())

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
	a.Router.NotFoundHandler = http.HandlerFunc(a.NotFound)

	/* REST ROUTES */
	a.Router.HandleFunc("/", a.Home).Methods("GET")
	a.Router.HandleFunc("/streams", a.Streams).Methods("GET")
	a.Router.HandleFunc("/users", a.Users).Methods("GET")
	a.Router.HandleFunc("/eventsub", a.EventsubRecievedNotification).Methods("POST")
	a.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	a.Router.HandleFunc("/test", a.Test).Methods("GET")
	/* a.Router.HandleFunc("/teams", a.TeamData).Methods("GET") */
}
