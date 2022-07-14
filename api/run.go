package api

import (
	"github.com/tikz/bcov/db"
)

func RunServer() {
	db.ConnectDB()

	engine := Endpoints()
	engine.Run()
}
