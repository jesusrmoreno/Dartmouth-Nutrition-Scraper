package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/go-errors/errors"
	"github.com/jesusrmoreno/nutrition-scraper/lib"
	"github.com/jesusrmoreno/nutrition-scraper/models"
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

	rDate := ""
	template := "01/02/06"
	if rDate = c.String("date"); rDate != "" {
		rDate = c.String("date")
	}

	date, err := time.Parse(template, rDate)
	if err != nil {
		log.Fatal("Unable to parse date make sure it looks like MM/dd/YY")
	}

	fmt.Println("Scraping info for:", rDate)

	url := ""
	shouldPost := false
	if url = c.String("post-url"); url != "" {
		url = c.String("post-url")
		shouldPost = true
	}

	// We want to get all Available SIDS
	sids, err := lib.AvailableSIDS()
	if err != nil {
		log.Fatal(err)
	}

	// How many nutrition routines we want to make at a time
	nutritionRoutines := 30

	// Time tracking stuff to see how long it took to run the entire program..
	startTime := time.Now()
	defer lib.TimeTrack(startTime, "Scrape")

	// venues is used to grab the information from the request as it finishes
	venues := make(chan models.VenueInfo)
	defer close(venues)

	for key, value := range sids {
		throttleRequests := make(chan bool, nutritionRoutines)
		defer close(throttleRequests)

		fmt.Println("Getting info for:", value)

		info := models.VenueInfo{}
		sid, err := lib.SID(key)
		if err != nil {
			log.Println("[ERROR]", err)
			continue
		}

		info.Venue = value
		info.Key = key
		info.SID = sid

		info.Menus, err = lib.MenuList(sid)
		if err != nil {
			log.Println("[ERROR]", err)
			continue
		}

		info.Meals, err = lib.MealList(sid)
		if err != nil {
			log.Println("[ERROR]", err)
		}

		for _, meal := range info.Meals {
			menuMeal := models.MenuMeal{
				Meal:  meal,
				Menus: models.MenuInfoSlice{},
			}
			for _, menu := range info.Menus {
				newRecipes, err := lib.
					RecipesMenuMealDate(sid, menu.ID, meal.ID, date)
				if err != nil {
					log.Println(err.(*errors.Error).ErrorStack())
					return
				}
				if len(newRecipes) > 0 {
					menuMeal.Menus = append(menuMeal.Menus, menu)
				}
				info.Recipes = append(info.Recipes, newRecipes...)
			}
			info.MealsList = append(info.MealsList, menuMeal)
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

		bar.FinishPrint("Got info for: " + value)
		if shouldPost {
			venueJSON, err := json.Marshal(info)
			if err != nil {
				log.Println(err)
			}
			postData := bytes.NewBuffer(venueJSON)
			log.Println("Making request!")
			resp, err := http.Post(url, "application/json", postData)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(resp)
				log.Println(resp.StatusCode)
			}
		}
		// Write a file to the directory it is run under with the output
		if c.Bool("write-files") {
			fileName := fmt.Sprintf("output_%s.json", info.Key)
			filePath := path.Join(pwd, fileName)
			b, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			err = ioutil.WriteFile(filePath, b, 0644)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		fmt.Println()
	}

}

func main() {
	// CLI configuration
	runtime.GOMAXPROCS(runtime.NumCPU())

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
		cli.StringFlag{
			Name:  "date",
			Usage: "Set the date we want to scrape. MM/dd/YY",
		},
		cli.StringFlag{
			Name:  "post-url, url",
			Usage: "If present will post to this url when finished.",
		},
	}
	app.Run(os.Args)
}
