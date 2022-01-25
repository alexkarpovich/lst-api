package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/infrastructure"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/repos"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/services"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"))

	db, err := infrastructure.NewPostgresDB(dataSourceName)
	if err != nil {
		fmt.Print(err)
		return
	}
	repos := repos.NewRepos(db)
	serverAddress := fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv("API_PORT"))
	fmt.Printf("User API server listening %s", serverAddress)

	services := services.NewServices(repos)

	srv, err := interfaces.NewHTTPServer(serverAddress, repos, services)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
