package interfaces

import (
	"net/http"
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/app/usecases"
	app_handlers "github.com/alexkarpovich/lst-api/src/internal/interfaces/handlers"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/middlewares"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/repos"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func configureRouter(repos *repos.Repos) http.Handler {
	baseRouter := mux.NewRouter().StrictSlash(true)
	authInterector := usecases.NewAuthInteractor(repos.User)
	app_handlers.ConfigureAuthHandler(authInterector, baseRouter)
	groupInterector := usecases.NewGroupInteractor(repos.Group, repos.Slice)
	app_handlers.ConfigureGroupHandler(groupInterector, baseRouter)
	sliceInterector := usecases.NewSliceInteractor(repos.Slice, repos.Expression)
	app_handlers.ConfigureSliceHandler(sliceInterector, baseRouter)

	return baseRouter
}

// NewServer - initialize HTTP Server
func NewHTTPServer(address string, repos *repos.Repos) (*http.Server, error) {
	// authEnforcer, err := casbin.NewEnforcer(
	// 	"config/auth_model.conf",
	// 	"config/policy.csv")

	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	chain := alice.New(
		handlers.CORS(headersOk, originsOk, methodsOk),
		middlewares.CurrentUser(repos.User),
		//mw.Authorizer(authEnforcer),
	).Then(configureRouter(repos))

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      chain,
		Addr:         address,
	}

	return server, nil
}
