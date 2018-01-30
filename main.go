package main

import (
	"github.com/mfathirirhas/TokoIjah/model"
	"github.com/mfathirirhas/TokoIjah/api"
)

var (
	port = ":8080"
)

func main() {
	db 		:= model.InitDB()
	server 	:= api.InitRouter(db)
	server.Run(port)
}