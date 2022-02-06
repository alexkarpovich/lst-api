package interfaces

import (
	"net/http"
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/app/usecases"
	app_handlers "github.com/alexkarpovich/lst-api/src/internal/interfaces/handlers"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/middlewares"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/repos"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func configureRouter(repos *repos.Repos, services *services.Services) http.Handler {
	baseRouter := mux.NewRouter().StrictSlash(true)

	authInterector := usecases.NewAuthInteractor(repos.User, services.Email)
	app_handlers.ConfigureAuthHandler(authInterector, baseRouter)

	userInterector := usecases.NewUserInteractor(repos.User)
	app_handlers.ConfigureUserHandler(userInterector, baseRouter)

	groupInterector := usecases.NewGroupInteractor(repos.Group, repos.Node, repos.User, services.Email)
	app_handlers.ConfigureGroupHandler(groupInterector, baseRouter)

	nodeInterector := usecases.NewNodeInteractor(repos.Node, repos.Group, repos.Expression)
	app_handlers.ConfigureNodeHandler(nodeInterector, baseRouter)

	expressionInterector := usecases.NewExpressionInteractor(repos.Expression)
	app_handlers.ConfigureExpressionHandler(expressionInterector, baseRouter)

	translationInterector := usecases.NewTranslationInteractor(repos.Translation)
	app_handlers.ConfigureTranslationHandler(translationInterector, baseRouter)

	langInterector := usecases.NewLangInteractor(repos.Lang)
	app_handlers.ConfigureLangHandler(langInterector, baseRouter)

	trainingInterector := usecases.NewTrainingInteractor(repos.Training, repos.Node, services.Training)
	app_handlers.ConfigureTrainingHandler(trainingInterector, baseRouter)

	return baseRouter
}

// NewServer - initialize HTTP Server
func NewHTTPServer(address string, repos *repos.Repos, services *services.Services) (*http.Server, error) {
	// authEnforcer, err := casbin.NewEnforcer(
	// 	"config/auth_model.conf",
	// 	"config/policy.csv")

	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"})

	chain := alice.New(
		handlers.CORS(headersOk, originsOk, methodsOk),
		middlewares.CurrentUser(repos.User),
		//mw.Authorizer(authEnforcer),
	).Then(configureRouter(repos, services))

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      chain,
		Addr:         address,
	}

	return server, nil
}
