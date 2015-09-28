package main

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/go-errors/errors"
	"github.com/jesusrmoreno/nutrition-scraper/lib"
	"github.com/jesusrmoreno/nutrition-scraper/models"
	// "io/ioutil"
	"log"
	// "net/http"
	// "strconv"
	"time"
	// "unicode/utf8"
)

func main() {
	// We want to get all Available SIDS
	sids, err := lib.AvailableSIDS()
	if err != nil {
		log.Fatal(err)
	}

	// We can open up to 50 connections but only about 15 of them can be handled
	// at any one time so we cap the concurrency at 15
	concurrency := 15

	startTime := time.Now()
	fmt.Println("Starting scrape at: ", startTime)
	defer lib.TimeTrack(startTime, "Entire Scrape")

	// sem will keep track of the number of concurrent connections and server as
	// a semaphore
	sem := make(chan bool, concurrency)

	// venues is not used right now but will be used to work with venues as they
	// are filled in with their information.
	venues := make(chan models.VenueInfo)

	// Make sure that our channels get closed
	defer func() { close(sem) }()
	defer func() { close(venues) }()

	for key, venue := range sids {

		fmt.Println("Starting", key, "at:", time.Now())

		info := models.VenueInfo{}
		sid, err := lib.GetSID(key)
		if err != nil {
			panic(err)
		}

		info.Venue = venue
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

		fmt.Println("Getting recipes for ", key, "at", time.Now())

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

		// This section is the part that benefits the most from concurrency
		// the top parts finish in about 5 seconds but this will take up to
		// 15 minutes if done one by one.
		bar := pb.StartNew(len(info.Recipes))
		bar.ShowSpeed = true
		bar.SetMaxWidth(80)
		for index := range info.Recipes {

			// Start a new goroutine for each nutrition request
			go func(key string, index int, info *models.VenueInfo) {

				// Read from the semaphore after we are done to free up a space for the
				// next connection
				bar.Increment()
				defer func() {
					<-sem
				}()

				_, err := lib.GetNutrients(info.SID, &info.Recipes[index])

				if err != nil {
					log.Println(err.(*errors.Error).ErrorStack())
				}

			}(key, index, &info)

			sem <- true

		}
		bar.FinishPrint("Done getting recipes for: " + key)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

}
