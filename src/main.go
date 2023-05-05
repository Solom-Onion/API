package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"solomon/api/routes"
	"solomon/api/util"

	"github.com/gorilla/mux"
)

type App struct {
	client *mongo.Client
}

var config map[string]string

func main() {
	config, err := util.ParseConfig()

	if err != nil {
		log.Fatalf("Error parsing config: %s", err)
		return
	}

	serverPort := fmt.Sprintf(":%s", config["server_port"])

	clientOptions := options.Client().ApplyURI(config["mongo_uri"])
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/view", func(w http.ResponseWriter, r *http.Request) {
		routes.ViewHandler(w, r, client, config)
	}).Methods("GET")

	router.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		routes.SearchHandler(w, r, client, config)
	}).Methods("GET")

	log.Printf("Server listening on port %s\n", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, util.MonkeyPatchMiddleMan(router)))
}
