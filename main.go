package main

import (
	"fmt"
	"time"

	"github.com/Mahcks/TheGoldenGator/api"
	"github.com/Mahcks/TheGoldenGator/configure"
	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/queries"

	_ "github.com/Mahcks/TheGoldenGator/docs"
)

// @title TheGoldenGator API
// @version 0.0.3
// @description Documentation for the public REST API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email mahcks@protonmail.com

// @license.name Apache
// @license.url https://insertthislater.com

// @BasePath /
func main() {
	err := database.Connect(configure.Config.GetString("mongo_uri"))
	if err != nil {
		fmt.Println("error connecting to database: ", err)
	}

	go queries.DoEvery(time.Minute*5, func(t time.Time) {
		err := queries.ViewCountPoll()
		if err != nil {
			panic(err)
		}
	})

	a := api.App{}
	a.Initialize()
	a.Run(":7500")
}
