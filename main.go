package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/go-errors/errors"
	"github.com/jesusrmoreno/nutrition-scraper/lib"
	"github.com/jesusrmoreno/nutrition-scraper/models"
	"github.com/jesusrmoreno/parse"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.New()

func init() {
	log.Formatter = new(prefixed.TextFormatter)
}

// State ...
type State struct {
	DB        *parse.Client
	Recipes   map[int]models.ParseRecipe
	Nutrients map[int]bool
	Offerings map[string]bool
}

func getRecipesFromParse(s *State, limit int) []models.ParseRecipe {
	if limit > 1000 {
		log.Warn("Parse has a max return limit of 1000 objects.")
		log.Warn("Using 1000 as the limit")
		limit = 1000
	}
	returnRecipes := []models.ParseRecipe{}
	skipValue := 0
	for {
		rawRecipes := []models.ParseRecipe{}
		status, errs := s.DB.Get(parse.Params{
			Class: "Recipe",
			Limit: limit,
			Skip:  skipValue,
		}, &rawRecipes)
		if errs != nil {
			log.Fatal("Could not get recipes, status:", status)
		}
		skipValue += limit
		if len(rawRecipes) == 0 {
			break
		}
		returnRecipes = append(returnRecipes, rawRecipes...)
	}
	return returnRecipes
}

func getOfferingsFromParse(s *State, limit int) []models.ParseOffering {
	if limit > 1000 {
		log.Warn("Parse has a max return limit of 1000 objects.")
		log.Warn("Using 1000 as the limit")
		limit = 1000
	}
	returnOfferings := []models.ParseOffering{}
	skipValue := 0
	for {
		dbOfferings := []models.ParseOffering{}
		status, errs := s.DB.Get(parse.Params{
			Class: "Offering",
			Limit: limit,
			Skip:  skipValue,
		}, &dbOfferings)
		if errs != nil {
			log.Fatal("Could not get offerings, status:", status)
		}
		skipValue += limit
		if len(dbOfferings) == 0 {
			break
		}
		returnOfferings = append(returnOfferings, dbOfferings...)
	}
	return returnOfferings
}

// InitParse ...
func InitParse(s *State) {
	dbRecipes := getRecipesFromParse(s, 1000)
	for _, dbRecipe := range dbRecipes {
		s.Recipes[dbRecipe.DartmouthID] = dbRecipe
	}

	dbOfferings := getOfferingsFromParse(s, 1000)
	for _, dbOffering := range dbOfferings {
		s.Offerings[dbOffering.UUID] = true
	}

	log.WithFields(logrus.Fields{
		"Recipes":   len(dbRecipes),
		"Offerings": len(dbOfferings),
	}).Info("In Database")
}

func mmRecipes(s *State, meal, menu int, rs models.RecipeInfoSlice) []string {
	mmR := []string{}
	for _, r := range rs {
		if r.MealID == meal && r.MenuID == menu {
			mmR = append(mmR, s.Recipes[r.ID].ObjectID())
		}
	}
	return mmR
}

func uniqueRecipes(rs models.RecipeInfoSlice) []models.RecipeInfo {
	unique := []models.RecipeInfo{}
	tally := map[int]bool{}
	for _, r := range rs {
		if tally[r.ID] == false {
			unique = append(unique, r)
			tally[r.ID] = true
		}
	}
	return unique
}

func saveToParse(s *State, v models.VenueInfo) {
	u := uniqueRecipes(v.Recipes)
	var duplicates, new int
	for _, recipe := range u {
		if s.Recipes[recipe.ID].DartmouthID != recipe.ID {

			returnObj, status, errs := s.DB.Post(models.ParseRecipe{
				Name:        recipe.Name,
				Category:    recipe.Category,
				DartmouthID: recipe.ID,
				Rank:        recipe.Rank,
				UUID:        lib.GetMD5Hash(recipe.Name),
				Nutrients:   recipe.Nutrients,
				Class:       "Recipe",
			})
			if errs != nil {
				log.Error(status)
				log.Error(errors.Errorf("Unable to post recipe with ID: %d", recipe.ID))
				continue
			}
			returnedRecipe := returnObj.(models.ParseRecipe)
			s.Recipes[recipe.ID] = returnedRecipe
			log.Info("Created new recipe with objectId: ", returnedRecipe.ObjectID())
			new++
		} else {
			duplicates++
		}
	}

	log.WithFields(logrus.Fields{
		"Saved":     new,
		"Duplicate": duplicates,
	}).Info("Scraped Recipes")

	offers := []models.ParseOffering{}
	duplicates, new = 0, 0
	for _, item := range v.MealsList {
		meal := item.Meal
		for _, menu := range item.Menus {
			day := v.Date.Day()
			month := int(v.Date.Month())
			year := v.Date.Year()
			uuidStr := fmt.Sprintf("%d%d%d%s%s%s", day, month, year, menu.Name, meal.Name, v.Key)
			uuid := lib.GetMD5Hash(uuidStr)
			if s.Offerings[uuid] == false {
				rs := mmRecipes(s, meal.ID, menu.ID, v.Recipes)
				offer := models.ParseOffering{
					Venue:    v.Key,
					Day:      v.Date.Day(),
					Month:    int(v.Date.Month()),
					Year:     v.Date.Year(),
					MenuName: menu.Name,
					MealName: meal.Name,
					Class:    "Offering",
					UUID:     uuid,
				}
				for _, s := range rs {
					offer.AddRecipe(s)
				}
				new++
				offers = append(offers, offer)
			} else {
				duplicates++
			}
		}
	}
	for _, o := range offers {
		returnObj, status, errs := s.DB.Post(o)
		if errs != nil {
			log.Error(status)
			log.Error(errors.Errorf("Unable to post recipe with ID: %s", o.UUID))
			continue
		}
		offering := returnObj.(models.ParseOffering)
		s.Offerings[offering.UUID] = true
		log.Info("Created new offering with objectId: ", offering.ObjectID())
	}
	log.WithFields(logrus.Fields{
		"Saved":     new,
		"Duplicate": duplicates,
	}).Info("Scraped Offerings")
}

func scrape(c *cli.Context) {
	p := parse.Client{
		BaseURL:       "https://api.parse.com/1",
		ApplicationID: "BAihtNGpVTx4IJsuuFV5f9LibJGnD1ZBOsnXk9qp",
		Key:           "zJYR2d3dFN3bXL6vUANZyoVLZ3bcTF7fpXTCrU7s",
	}
	s := State{
		DB:        &p,
		Recipes:   make(map[int]models.ParseRecipe),
		Nutrients: make(map[int]bool),
		Offerings: make(map[string]bool),
	}

	if c.Bool("mock") {
		log.Info("Mocked Scrape")
		InitParse(&s)
		file, err := os.Open("output_DDS.json")
		if err != nil {
			log.Fatal(err)
		}
		info := models.VenueInfo{}
		if err := json.NewDecoder(file).Decode(&info); err != nil {
			log.Fatal(err)
		}
		saveToParse(&s, info)
		log.Info("End Mocked Scrape")
		return
	}
	InitParse(&s)
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
	if rDate = c.String("startDate"); rDate != "" {
		rDate = c.String("startDate")
	}

	date, err := time.Parse(template, rDate)
	if err != nil {
		log.Fatal("Unable to parse date make sure it looks like MM/dd/YY")
	}

	dateArray := []time.Time{}
	fmt.Println("Will try to scrape for:")
	for i := 0; i < 7; i++ {
		dateToAdd := date.AddDate(0, 0, i)
		fmt.Println("  ", dateToAdd.Format(template))
		dateArray = append(dateArray, dateToAdd)
	}
	fmt.Println()

	shouldPost := c.Bool("save")

	// Time tracking stuff to see how long it took to run the entire program..
	startTime := time.Now()
	defer lib.TimeTrack(startTime, "Scrape")

	for _, date := range dateArray {
		fmt.Println("Scraping info for:", date.Format(template))

		// We want to get all Available SIDS
		sids, err := lib.AvailableSIDS()
		if err != nil {
			log.Fatal(err)
		}

		// How many nutrition routines we want to make at a time
		nutritionRoutines := 50

		for key, value := range sids {
			throttleRequests := make(chan bool, nutritionRoutines)
			defer close(throttleRequests)

			fmt.Println("Getting info for:", value)

			info := models.VenueInfo{
				Date: date,
			}
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
				saveToParse(&s, info)
			}
			// Write a file to the directory it is run under with the output
			if c.Bool("write-files") {
				fileName := fmt.Sprintf("output_%s.json", info.Key)
				filePath := path.Join(pwd, fileName)
				b, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					fmt.Println("error:", err)
				}
				fmt.Println("Wrote to:", fileName)
				err = ioutil.WriteFile(filePath, b, 0644)
				if err != nil {
					log.Println(err)
					continue
				}
			}
			fmt.Println()
		}
	}
}

func main() {

	// CLI configuration
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "nutrition-scraper"
	app.Usage = "A tool for scraping the Dartmouth Dining Services menu."
	app.Action = scrape
	app.Version = "0.1.9"
	// Add more flags at the end of this slice
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "write-files, wf",
			Usage: "If present will write the scraped information to json files.",
		},
		cli.BoolFlag{
			Name:  "mock",
			Usage: "Will use output_CYC.json data to post. Must also use --url",
		},
		cli.StringFlag{
			Name:  "startDate, sd",
			Usage: "Set the date we want to scrape. MM/dd/YY",
		},
		cli.BoolFlag{
			Name:  "save",
			Usage: "Include to save to parse",
		},
	}
	app.Run(os.Args)
}
