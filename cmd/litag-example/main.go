package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/fwojciec/litag-example/generated/gqlgen" // remember to update your username
	"github.com/fwojciec/litag-example/postgres"         // remember to update your username
	"github.com/fwojciec/litag-example/resolvers"        // remember to update your username
)

func main() {
	// initialize the repo
	repo, err := postgres.NewRepo("dbname=litag_example_db sslmode=disable")
	if err != nil {
		panic(err)
	}

	// initialize the GraphQL handler
	gqlHandler := handler.GraphQL(gqlgen.NewExecutableSchema(gqlgen.Config{
		Resolvers: &resolvers.Resolver{
			Repo: repo,
		},
	}))

	// configure the server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Playground("GraphQL Playground", "/query"))
	mux.HandleFunc("/query", gqlHandler)

	// run the server
	port := ":8080"
	fmt.Printf("🚀 Server ready at http://localhost%s\n", port)
	log.Fatalln(http.ListenAndServe(port, mux))
}
