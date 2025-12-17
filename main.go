package main

import (
	"fmt"
	"go-app/controllers"
	"go-app/db"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var LogPath = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: %s (%s)", r.Host, r.RequestURI, r.Method))
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Инициализация БД — теперь явная, берётся из переменных окружения.
	if err := db.InitDBFromEnv(); err != nil {
		log.Fatal("DB init failed: ", err)
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Warn("Failed to close DB: ", err)
		}
	}()

	router := mux.NewRouter()

	router.Use(controllers.MetricsMiddleware)
	router.Use(LogPath)

	router.HandleFunc("/notes", controllers.NoteQuery).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/notes", controllers.NoteCreate).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteRetrieve).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteUpdate).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteDelete).Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/analytics/summary", controllers.AnalyticsSummary).Methods(http.MethodGet)
	router.HandleFunc("/analytics/notes/count", controllers.AnalyticsNotesCount).Methods(http.MethodGet)
	router.HandleFunc("/analytics/notes/avg-length", controllers.AnalyticsAvgNoteLength).Methods(http.MethodGet)

	log.Info("Listening on 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
