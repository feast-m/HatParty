package main

import (
	"log"
	"net/http"

	"github.com/feastM/HatParty/app"
)

func main() {
	app.HandleRequests()
	log.Fatal(http.ListenAndServe(":10000", nil))
}
