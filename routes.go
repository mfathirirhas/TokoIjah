package main

import (
	"net/http"
	"github.com/mfathirirhas/TokoIjah/inventory/index"
)

// Routes handle incoming request and map them to correspondent APIs
func Routes() {
	http.HandleFunc("/", index.Home)
}