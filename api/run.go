package api

import (
	"github.com/tikz/bcov/db"
)

// RunServer starts serving Gin.
func RunServer() {
	db.ConnectDB()

	engine := Endpoints()
	engine.Run()
}
