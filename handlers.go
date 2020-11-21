package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Link struct
type Link struct {
	Link  string
	Code  string
	SLink string
}

// ErrLink for not valid link url
type ErrLink struct {
	Link string
	Err  string
}

// Index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, nil)
}

// API page
func apiHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	url := vars.Get("url")

	// Check URL address
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		errjson, err := json.Marshal(ErrLink{url, "Incorrect url format"})
		if err != nil {
			log.Fatal(err)
		}

		// Return error json message
		w.Header().Set("Content-Type", "application/json")
		w.Write(errjson)
		return
	}

	code := genCode(url)

	// Search for a link with the same code in the database
	findlink, err := collection.Find(ctx, bson.M{"code": code})
	if err != nil {
		log.Fatal(err)
	}

	// If there is one, then we give its link
	for findlink.Next(ctx) {
		var elem Link

		err := findlink.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		jsonresp, err := json.Marshal(elem)
		if err != nil {
			log.Fatal(err)
		}

		// Return json
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonresp)
		return
	}

	tinylink := fmt.Sprintf("%s/s/%s", siteURL, code)
	slink := Link{url, code, tinylink}

	// Insert new link to database
	_, err = collection.InsertOne(ctx, slink)
	if err != nil {
		log.Fatal(err)
	}

	jsonresp, err := json.Marshal(slink)
	if err != nil {
		log.Fatal(err)
	}

	// Return json
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonresp)
}

// Slinky redirect page
func slinkyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	var link Link

	// Find link by code
	slink, err := collection.Find(ctx, bson.M{"code": code})
	if err != nil {
		log.Fatal(err)
	}

	for slink.Next(ctx) {
		var elem Link
		err := slink.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		link = elem
	}

	fmt.Printf("Redirect %s%s -> %s\n", r.Host, r.URL.Path, link.Link)

	// Redirect client to link from database
	http.Redirect(w, r, link.Link, 301)
}
