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
	DB            *parse.Client
	Recipes       map[int]models.ParseRecipe
	Nutrients     map[int]bool
	Offerings     map[string]models.ParseOffering
	Subscriptions map[int][]string
	Notifications map[string]models.ParseNotification
}

func getNotificationsFromParse(s *State, limit int) []models.ParseNotification {
	if limit > 1000 {
		log.Warn("Parse has a max return limit of 1000 objects.")
		log.Warn("Using 1000 as the limit")
		limit = 1000
	}
	skipValue := 0
	notifications := []models.ParseNotification{}
	for {
		newNotifications := []models.ParseNotification{}
		status, errs := s.DB.Get(parse.Params{
			Class: "Notification",
			Limit: limit,
			Skip:  skipValue,
		}, &newNotifications)

		if errs != nil {
			log.Fatal("Could not get notifications, status:", status)
		}
		skipValue += limit

		if len(newNotifications) == 0 {
			break
		}
		notifications = append(notifications, newNotifications...)
	}
	return notifications
}

func getSubscriptionsFromParse(s *State, limit int) models.SubscriptionSlice {
	if limit > 1000 {
		log.Warn("Parse has a max return limit of 1000 objects.")
		log.Warn("Using 1000 as the limit")
		limit = 1000
	}
	skipValue := 0
	subscriptions := models.SubscriptionSlice{}
	for {
		newSubscriptions := models.SubscriptionSlice{}
		status, errs := s.DB.Get(parse.Params{
			Class: "Subscription",
			Limit: limit,
			Skip:  skipValue,
		}, &newSubscriptions)

		if errs != nil {
			log.Fatal("Could not get subscriptions, status:", status)
		}
		skipValue += limit

		if len(newSubscriptions) == 0 {
			break
		}
		subscriptions = append(subscriptions, newSubscriptions...)
	}
	return subscriptions
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

func offeringExists(s *State, vK string, m, ml string, d time.Time) bool {
	uuidStr := fmt.Sprintf("%d%d%d%s%s%s",
		d.Day(), int(d.Month()), d.Year(), m, ml, vK)
	uuid := lib.GetMD5Hash(uuidStr)
	return s.Offerings[uuid].ObjectID() != ""
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
		s.Offerings[dbOffering.UUID] = dbOffering
	}

	dbSubscriptions := getSubscriptionsFromParse(s, 1000)
	for _, sub := range dbSubscriptions {
		for _, recipe := range sub.Recipes {
			s.Subscriptions[recipe] = append(s.Subscriptions[recipe], sub.User.ObjectID)
		}
	}

	dbNotifications := getNotificationsFromParse(s, 1000)
	for _, not := range dbNotifications {
		s.Notifications[not.UUID] = not
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

// SetDietaryInfo ...
func SetDietaryInfo(n *models.NutrientInfoResponse, title string) *models.NutrientInfoResponse {
	info := models.TitleToProps(title)
	for _, item := range info {
		switch item {
		case "l/o":
			n.Result.Vegetarian = true
		case "gf":
			n.Result.Gluten = true
		case "r":
			n.Result.Local = true
		case "k":
			n.Result.Kosher = true
		case "h":
			n.Result.Halal = true
		case "v":
			n.Result.Vegan = true
		case "e":
			n.Result.Eggs = true
		case "f":
			n.Result.Fish = true
		case "d":
			n.Result.Dairy = true
		case "n":
			n.Result.TreeNuts = true
		case "p":
			n.Result.Peanuts = true
		case "pk":
			n.Result.Pork = true
		case "sb":
			n.Result.Soy = true
		case "sf":
			n.Result.ShellFish = true
		case "w":
			n.Result.Wheat = true
		}
	}
	return n
}

func saveRecipes(s *State, v models.VenueInfo) {
	u := uniqueRecipes(v.Recipes)
	var duplicates, new int
	for _, recipe := range u {
		if s.Recipes[recipe.ID].DartmouthID != recipe.ID {
			c := models.CreatedBy{
				Kind:      "Pointer",
				ClassName: "_User",
				ObjectID:  "95xfYTL7GG",
			}
			returnObj, status, errs := s.DB.Post(models.ParseRecipe{
				Name:        models.RemoveMetaData(recipe.Name),
				Category:    recipe.Category,
				DartmouthID: recipe.ID,
				Rank:        recipe.Rank,
				UUID:        lib.GetMD5Hash(models.RemoveMetaData(recipe.Name)),
				Nutrients:   *SetDietaryInfo(&recipe.Nutrients, recipe.Name),
				Class:       "Recipe",
				CreatedBy:   c,
			})
			if errs != nil || status == 400 {
				log.Error(status)
				log.Error(errors.Errorf("Unable to post recipe with ID: %d", recipe.ID))
				continue
			}
			returnedRecipe := returnObj.(models.ParseRecipe)
			s.Recipes[recipe.ID] = returnedRecipe
			log.Debug("Created new recipe with objectId: ", returnedRecipe.ObjectID())
			new++
		} else {
			duplicates++
		}
	}

	log.WithFields(logrus.Fields{
		"Saved":     new,
		"Duplicate": duplicates,
	}).Info("Scraped Recipes")
}

func saveOfferings(s *State, v models.VenueInfo) {
	offers := []models.ParseOffering{}
	duplicates, new := 0, 0
	for _, item := range v.MealsList {
		meal := item.Meal
		for _, menu := range item.Menus {
			day := v.Date.Day()
			month := int(v.Date.Month())
			year := v.Date.Year()
			uuidStr := fmt.Sprintf(
				"%d%d%d%s%s%s", day, month, year, menu.Name, meal.Name, v.Key)
			uuid := lib.GetMD5Hash(uuidStr)

			if s.Offerings[uuid].ObjectID() == "" {
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
		s.Offerings[offering.UUID] = offering
		log.Debug("Created new offering with objectId: ", offering.ObjectID())
	}
	log.WithFields(logrus.Fields{
		"Saved":     new,
		"Duplicate": duplicates,
	}).Info("Scraped Offerings")
}

func saveToParse(s *State, v models.VenueInfo) {
	saveRecipes(s, v)
	saveOfferings(s, v)
}

// NameToNutrientMigration ...
func NameToNutrientMigration(s *State) {
	for _, recipe := range s.Recipes {
		recipe.Nutrients = *SetDietaryInfo(&recipe.Nutrients, recipe.Name)
		fmt.Println(models.RemoveMetaData(recipe.Name))
		recipe.Name = models.RemoveMetaData(recipe.Name)

		x := struct {
			Name      string                      `json:"name"`
			Nutrients models.NutrientInfoResponse `json:"nutrients"`
			ID        string                      `json:"objectId"`
			UUID      string                      `json:"uuid"`
		}{
			Name:      models.RemoveMetaData(recipe.Name),
			Nutrients: *SetDietaryInfo(&recipe.Nutrients, recipe.Name),
			ID:        recipe.ObjectID(),
			UUID:      lib.GetMD5Hash(models.RemoveMetaData(recipe.Name)),
		}
		xString, _ := json.Marshal(x)
		fmt.Println(string(xString))
		_, status, errs := s.DB.Put(x, "Recipe", recipe.ID)

		if errs != nil || status == 400 {
			log.Error(status)
			log.Error(errs)
			log.Error(errors.Errorf("Unable to post recipe with ID: %s", recipe.ID))
			break
		}
		// time.Sleep(1 * time.Second)
	}
}

func scrape(c *cli.Context) {
	log.Info("Initializing Scraper")
	p := parse.Client{
		BaseURL:       "https://api.parse.com/1",
		ApplicationID: "BAihtNGpVTx4IJsuuFV5f9LibJGnD1ZBOsnXk9qp",
		Key:           "zJYR2d3dFN3bXL6vUANZyoVLZ3bcTF7fpXTCrU7s",
	}

	s := State{
		DB:            &p,
		Recipes:       make(map[int]models.ParseRecipe),
		Nutrients:     make(map[int]bool),
		Offerings:     make(map[string]models.ParseOffering),
		Subscriptions: make(map[int][]string),
		Notifications: make(map[string]models.ParseNotification),
	}

	if c.Bool("subscriptions") {
		// InitParse(&s)
		getSubscriptionsFromParse(&s, 100)
		return
	}

	if c.Bool("nameNutrientMigration") {
		fmt.Println("Running Migration...")
		InitParse(&s)
		NameToNutrientMigration(&s)
		return
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
	for i := 0; i < 7; i++ {
		dateToAdd := date.AddDate(0, 0, i)
		dateArray = append(dateArray, dateToAdd)
	}
	shouldPost := c.Bool("save")
	notificationsToCreate := []models.Notification{}
	for _, date := range dateArray {
		log.WithFields(logrus.Fields{
			"date": date.Format(template),
		}).Info("Start Scrape")

		// We want to get all Available SIDS
		sids, err := lib.AvailableSIDS()
		if err != nil {
			log.Fatal(err)
		}
		log.WithFields(logrus.Fields{
			"count": len(sids),
		}).Info("SIDS")

		// How many nutrition routines we want to make at a time
		nutritionRoutines := 50

		for key, value := range sids {
			throttleRequests := make(chan bool, nutritionRoutines)
			defer close(throttleRequests)
			log.WithFields(logrus.Fields{
				"venue": key,
			}).Info("Venue Scrape")
			info := models.VenueInfo{
				Date: date,
			}
			sid, err := lib.SID(key)
			if err != nil {
				log.Error(err)
				continue
			}

			info.Venue = value
			info.Key = key
			info.SID = sid

			info.Menus, err = lib.MenuList(sid)
			log.WithFields(logrus.Fields{
				"count": len(info.Menus),
			}).Info("Got Menus")
			if err != nil {
				log.Error(err)
				continue
			}

			info.Meals, err = lib.MealList(sid)
			if err != nil {
				log.Error(err)
			}
			log.WithFields(logrus.Fields{
				"count": len(info.Meals),
			}).Info("Got Meals")

			for _, meal := range info.Meals {
				menuMeal := models.MenuMeal{
					Meal:  meal,
					Menus: models.MenuInfoSlice{},
				}
				for _, menu := range info.Menus {
					newRecipes, err := lib.
						RecipesMenuMealDate(sid, menu.ID, meal.ID, date)
					if err != nil {
						log.Error(err)
						continue
					}
					for _, recipe := range newRecipes {
						if len(s.Subscriptions[recipe.ID]) > 0 {
							notificationsToCreate = append(notificationsToCreate, models.Notification{
								RecipeID: recipe.ID,
								Name:     models.RemoveMetaData(recipe.Name),
								Day:      date.Day(),
								Month:    int(date.Month()),
								Year:     date.Year(),
								OnDate:   date,
								MenuName: menu.Name,
								MealName: meal.Name,
								Venue:    info.Key,
							})
						}
					}
					// We need to scrape the recipes so that we can create notifications
					// but if the offering exists then we can just skip everything else
					if offeringExists(&s, info.Key, menu.Name, meal.Name, date) {
						log.WithFields(logrus.Fields{
							"meal":  meal.ID,
							"menu":  menu.ID,
							"venue": info.Key,
							"date":  date.Format(template),
						}).Info("Offering Exists")
						// newRecipes, err := lib.RecipesMenuMealDate(sid, menu.ID, meal.ID, date)
						continue
					}
					if len(newRecipes) > 0 {
						menuMeal.Menus = append(menuMeal.Menus, menu)
					}
					info.Recipes = append(info.Recipes, newRecipes...)
				}
				info.MealsList = append(info.MealsList, menuMeal)
			}

			// This section is the part that benefits the most from concurrency
			// the top parts finish in about 5 seconds but this will take up to
			// 15 minutes if done one by one.
			log.WithFields(logrus.Fields{
				"count": len(info.Recipes),
			}).Info("Start Recipe Scrape")
			for index := range info.Recipes {
				// Start a new goroutine for each nutrition request
				go func(key string, index int, info *models.VenueInfo) {
					// Read from the semaphore after we are done to free up a space for
					// the next connection.
					defer func() { <-throttleRequests }()

					// GetNutrients returns a pointer but we don't really care about it
					// simply ignore it. We pass &info.Recipes[index] so that the actual
					// pointer in the info object will be updated, otherwise a copy
					// will be worked on and we won't see the result
					_, err := lib.GetNutrients(info.SID, &info.Recipes[index])
					if err != nil {
						log.Error(err)
					}

				}(key, index, &info)
				/// Add our request to the list of running requests.
				throttleRequests <- true
			}

			// We want to fill them up by default..
			for i := 0; i < cap(throttleRequests); i++ {
				throttleRequests <- true
			}

			log.WithFields(logrus.Fields{
				"count": len(info.Recipes),
			}).Info("Finish Recipe Scrape")
			if shouldPost {
				saveToParse(&s, info)
			}
			log.WithFields(logrus.Fields{
				"venue": info.Key,
			}).Info("Finish Venue Scrape")
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
		}
	}
	ns := createNotifications(&s, notificationsToCreate)
	saveNotifications(&s, ns)
	removeOldNotifications(&s)
}

func saveNotifications(s *State, ns []models.ParseNotification) {
	throttleRequests := make(chan bool, 20)
	skipped := 0
	defer func(throttleRequests chan bool) {
		close(throttleRequests)
		fmt.Println("Duplicates: ", skipped)
	}(throttleRequests)
	// We want to fill them up by default..
	for _, n := range ns {
		if s.Notifications[n.UUID].UUID != n.UUID {
			go func(n models.ParseNotification) {
				defer func() {
					<-throttleRequests
				}()
				_, status, errs := s.DB.Post(n)
				if status == 200 || status == 201 {
					fmt.Println("Successfully created notification with ID:", n.UUID)
				}
				if errs != nil {
					fmt.Print(status)
					fmt.Print(errors.Errorf("Unable to post Notification with ID: %d", n.RecipeID))
				}
			}(n)
			throttleRequests <- true
		} else {
			skipped++
		}
	}
	for i := 0; i < cap(throttleRequests); i++ {
		throttleRequests <- false
	}
}

func removeOldNotifications(s *State) {
	throttleRequests := make(chan bool, 20)
	defer close(throttleRequests)
	toDelete := []models.ParseNotification{}
	for _, n := range s.Notifications {
		// Deletes things from the day before gives a day leway...
		if n.OnDate.ISO.Before(time.Now().AddDate(0, 0, -1)) {
			toDelete = append(toDelete, n)
		}
	}
	for _, n := range toDelete {
		go func(n models.ParseNotification) {
			defer func() {
				<-throttleRequests
			}()
			status, errs := s.DB.Delete(parse.Params{
				Class:    "Notification",
				ObjectID: n.ObjectID(),
			}, nil)
			if status == 200 || status == 201 {
				fmt.Println("Successfully deleted notification with ID:", n.UUID)
			}
			if errs != nil {
				fmt.Print(status)
				fmt.Print(errors.Errorf("Unable to delete Notification with ID: %s", n.UUID))
			}
		}(n)
		throttleRequests <- true
	}
	for i := 0; i < cap(throttleRequests); i++ {
		throttleRequests <- false
	}
}

func createNotifications(s *State, ns []models.Notification) []models.ParseNotification {
	// We do concurrency here because it takes a while for the uuid to be
	// calculated and we need to handle a lot of notification creation at a time
	notificationChan := make(chan models.ParseNotification)
	go func() {
		for _, n := range ns {
			for _, userID := range s.Subscriptions[n.RecipeID] {
				p := models.ParseNotification{
					Class:    "Notification",
					RecipeID: n.RecipeID,
					Name:     n.Name,
					Day:      n.Day,
					Month:    n.Month,
					Year:     n.Year,
					OnDate: models.DateObject{
						Type: "Date",
						ISO:  n.OnDate,
					},
					MenuName: n.MenuName,
					MealName: n.MealName,
					Venue:    n.Venue,
					For: models.CreatedBy{
						Kind:      "Pointer",
						ClassName: "_User",
						ObjectID:  userID,
					},
				}
				p.UUID = p.GenerateUUID()
				notificationChan <- p
			}
		}
		close(notificationChan)
	}()
	toPost := []models.ParseNotification{}
	for p := range notificationChan {
		toPost = append(toPost, p)
	}
	return toPost
}

func main() {
	// CLI configuration
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "nutrition-scraper"
	app.Usage = "A tool for scraping the Dartmouth Dining Services menu."
	app.Action = scrape
	app.Version = "0.1.10"
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
		cli.BoolFlag{
			Name:  "nameNutrientMigration",
			Usage: "Will migrate the recipes...",
		},
		cli.StringFlag{
			Name:  "startDate, sd",
			Usage: "Set the date we want to scrape. MM/dd/YY",
		},
		cli.BoolFlag{
			Name:  "save",
			Usage: "Include to save to parse",
		},
		cli.BoolFlag{
			Name:  "subscriptions",
			Usage: "Include to save to parse",
		},
	}
	app.Run(os.Args)
}
