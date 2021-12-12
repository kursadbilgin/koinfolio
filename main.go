package main

import (
	"koinfolio/Routers"
	"koinfolio/db"
)

func main() {
	db.InitMongoDB()
	r := Routers.SetupRouter()

	r.Run(":8080")
}
