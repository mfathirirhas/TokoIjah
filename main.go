package main

import (
	"github.com/mfathirirhas/TokoIjah/config"
)

func main() {
	// StartServer() start the server
	config.StartServer()
	
	// StopServer() detect interruption and stop the server
	config.StopServer()
}