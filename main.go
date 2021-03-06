package main

import (
	"fmt"

	"github.com/Mahcks/TheGoldenGator/api"
	"github.com/Mahcks/TheGoldenGator/configure"
	"github.com/Mahcks/TheGoldenGator/database"
	"github.com/Mahcks/TheGoldenGator/queries"
	"github.com/jasonlvhit/gocron"

	_ "github.com/Mahcks/TheGoldenGator/docs"
)

// @title TheGoldenGator API
// @version 0.7.3
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

	s := gocron.NewScheduler()
	s.Every(5).Minutes().Do(func() {
		err := queries.ViewCountPoll()
		if err != nil {
			fmt.Println("Error updating viewcount", err)
		}

		errCache := queries.GetStreamsDeleteCache()
		if errCache != nil {
			fmt.Println("Error deleting stream cache", errCache)
		}
		fmt.Println("Deleted stream cache")

		errStreams := queries.UpdateStreams()
		if errStreams != nil {
			fmt.Println("Error updating streams", errStreams)
		}
		fmt.Println("Updated stream statuses")
	})

	s.Start()
	a := api.App{}
	a.Initialize()
	a.Run(":7500")
}
