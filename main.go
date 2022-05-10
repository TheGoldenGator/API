package main

import (
	"fmt"

	"github.com/Mahcks/golden-gator-api/api"
	"github.com/Mahcks/golden-gator-api/configure"
	"github.com/Mahcks/golden-gator-api/database"
)

func main() {
	err := database.Connect(configure.Config.GetString("mongo_uri"))
	if err != nil {
		fmt.Println("error connecting to database: ", err)
	}

	a := api.App{}
	a.Initialize()
	a.Run(":7500")
}
