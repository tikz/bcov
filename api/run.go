package api

import (
	"bcov/db"
)

func RunServer() {
	db.ConnectDB()

	engine := Endpoints()
	engine.Run()
}
