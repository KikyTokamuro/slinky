package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collection *mongo.Collection
	client     *mongo.Client
	ctx        = context.TODO()

	indexTmpl *template.Template

	siteURL = "sliy.herokuapp.com"
)

func init() {
	// Mongo client
	clientOptions := options.Client().ApplyURI(
		"CONNECTION STRING TO MONGODB",
	)

	// Mongo client connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping mongo
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to collection
	collection = client.Database("slinky").Collection("links")

	fmt.Println("Connected to MongoDB!")

	// Index.html template
	indexTmpl = template.Must(template.New("index.html").ParseFiles(
		path.Join("views", "index.html"),
	))
}

func main() {
	// CORS
	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET"})

	// Router
	r := mux.NewRouter()

	// URLs handlers
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/api/", apiHandler).Methods("GET")
	r.HandleFunc("/s/{code}", slinkyHandler)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))),
	)

	// Run
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", getEnv("PORT", "3000")),
		handlers.CORS(originsOk, headersOk, methodsOk)(r),
	)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
