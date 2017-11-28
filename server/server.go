package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/remimenke/mastermind/api/mastermind"
	"github.com/remimenke/mastermind/data"
	"github.com/remimenke/mastermind/params"
)

// Server owns all server methods and state
type Server struct{}

// New returns a new server
func New() *Server {
	s := &Server{}
	return s
}

// Run starts the server
func (s *Server) Run() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("JAWSDB_URL")
	}

	db, err := data.NewDatabase(dbURL)
	if err != nil {
		panic(err)
	}

	handler := gziphandler.GzipHandler(
		handler(db),
	)

	server := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		Handler:           handler,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 16,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	go func() {
		<-signalChannel

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		return
	}
	if err != nil {
		panic(err)
	}
}

func handler(db *gorm.DB) http.Handler {

	ctrl := mastermind.New(db)

	// API router
	api := httprouter.New()
	api.RedirectTrailingSlash = false

	api.POST("/v1/register", params.ParamHandler(ctrl.Register()))
	api.POST("/v1/create-game", params.ParamHandler(ctrl.CreateGame()))
	api.POST("/v1/join-game", params.ParamHandler(ctrl.JoinGame()))
	api.POST("/v1/try-code", params.ParamHandler(ctrl.TryCode()))
	api.POST("/v1/validate-code", params.ParamHandler(ctrl.ValidateCode()))

	// Main router
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*60)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Println(err)
				log.Println(string(debug.Stack()))
			}
		}()

		api.ServeHTTP(w, r.WithContext(ctx))
	}))

	return handler
}
