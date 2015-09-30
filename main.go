package main

import (
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/go-errors/errors"
	"github.com/jesusrmoreno/nutrition-scraper/lib"
	"github.com/jesusrmoreno/nutrition-scraper/models"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func scrape(c *cli.Context) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working directory!")
	}
	if c.Bool("write-files") {
		fmt.Println()
		fmt.Println("Output files will be placed in", pwd)
	}
	// We want to get all Available SIDS
	sids, err := lib.AvailableSIDS()
	if err != nil {
		log.Fatal(err)
	}
	// Formula for number of concurrent goroutines
	// total := (nutritionRoutines * venueRoutines) + venueRoutines
	// How many nutrition routines we want to make at a time
	nutritionRoutines := 30
	// How many venue routines we want to open at a time
	venueRoutines := 1

	// Time tracking stuff to see how long it took to run the entire program..
	startTime := time.Now()
	defer lib.TimeTrack(startTime, "Scrape")

	// venues is used to grab the information from the request as it finishes
	venues := make(chan models.VenueInfo)
	defer close(venues)

	// venueThrottle is used to throttle the requests that can be sent at once.
	// Right now we set it to a single request because the Dartmouth API will only
	// allowaround 15 connections at any one time and we use them all up during
	// recipes/ It's slightly faster to do it one location at a time with each one
	// using up the 15 connections that it is to have them running at once
	// with each using 5 or so connections.
	// Removing this limit is simple enough, just remove all reveferences to
	// venueThrottle
	venueThrottle := make(chan bool, venueRoutines)
	defer close(venueThrottle)

	for key, value := range sids {
		// We can a goroutine for each location to run in the background and to
		// throttle our requests
		go func(key, value string) {

			// Because the channel is buffered to accept only one value, feeding it
			// with true will stop any other goroutines from running.
			venueThrottle <- true

			// We defer this so that when we're finished it removes one of the items
			// from the venueThrottle channel and allows the next one in
			defer func() { <-venueThrottle }()
			// throttleRequests much like venueThrottle stops too many requests from
			// firing. At this point we want to fire around 15 for so that we get
			// 15 nutrition objects at at time. We can set it to anything and 15
			// will be the max number of connections that will get a response but
			// setting it too high will cause their API to complain that there are
			// too many connections open.
			throttleRequests := make(chan bool, nutritionRoutines)
			defer close(throttleRequests)

			fmt.Println("\nGetting info for:", value)
			info := models.VenueInfo{}
			sid, err := lib.GetSID(key)
			if err != nil {
				panic(err)
			}
			info.Venue = value
			info.Key = key
			info.SID = sid

			info.Menus, err = lib.GetMenuList(sid)
			if err != nil {
				log.Fatal(err)
			}

			info.Meals, err = lib.GetMealList(sid)
			if err != nil {
				log.Fatal(err)
			}

			for _, menu := range info.Menus {
				for _, meal := range info.Meals {
					newRecipes, err := lib.GetRecipesMenuMealDate(sid, menu.ID, meal.ID)
					if err != nil {
						log.Println(err.(*errors.Error).ErrorStack())
						return
					}
					info.Recipes = append(info.Recipes, newRecipes...)
				}
			}

			// Pretty progress bar stuff...
			bar := pb.StartNew(len(info.Recipes))
			bar.ShowSpeed = true
			bar.SetMaxWidth(80)
			// This section is the part that benefits the most from concurrency
			// the top parts finish in about 5 seconds but this will take up to
			// 15 minutes if done one by one.
			for index := range info.Recipes {
				// Start a new goroutine for each nutrition request
				go func(key string, index int, info *models.VenueInfo) {
					// Read from the semaphore after we are done to free up a space for
					// the next connection.
					defer func() {
						<-throttleRequests
						// Make the progress bar go up.
						bar.Increment()
					}()

					// GetNutrients returns a pointer but we don't really care about it
					// simply ignore it. We pass &info.Recipes[index] so that the actual
					// pointer in the info object will be updated, otherwise a copy
					// will be worked on and we won't see the result
					_, err := lib.GetNutrients(info.SID, &info.Recipes[index])
					if err != nil {
						log.Println(err.(*errors.Error).ErrorStack())
					}

				}(key, index, &info)
				/// Add our request to the list of running requests.
				throttleRequests <- true
			}

			// We want to fill them up by default..
			for i := 0; i < cap(throttleRequests); i++ {
				throttleRequests <- true
			}

			// Place the final info object in the venues channel
			venues <- info
			bar.FinishPrint("Got info for: " + value)
		}(key, value)
	}

	// We know that there are len(sids) venues to look at so we wait until we have
	// received that many objects to quit the program.
	for venueIndex := 0; venueIndex < len(sids); {
		select {
		case venue := <-venues:
			// Write a file to the directory it is run under with the output
			if c.Bool("write-files") {
				fileName := fmt.Sprintf("output_%s.json", venue.Key)
				filePath := path.Join(pwd, fileName)
				b, err := json.MarshalIndent(venue, "", "  ")
				if err != nil {
					fmt.Println("error:", err)
				}
				err = ioutil.WriteFile(filePath, b, 0644)
			}
			venueIndex++
		}
	}
}

func main() {
	// CLI configuration
	app := cli.NewApp()
	app.Name = "nutrition-scraper"
	app.Usage = "A tool for scraping the Dartmouth Dining Services menu."
	app.Action = scrape
	// Add more flags at the end of this slice
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "write-files, wf",
			Usage: "If present will write the scraped information to json files.",
		},
	}
	app.Run(os.Args)
}
