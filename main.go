package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "github.com/cheggaaa/pb"
	"github.com/go-errors/errors"
	"github.com/jesusrmoreno/nutrition-scraper/constants"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"time"
	"unicode/utf8"
)

// makeRequest is a helper function that takes the parameters as a string and
// executes the http request returning any errors, or nil and the body as a
// byte array
func makeRequest(params string) ([]byte, error) {
	url := urlBuilder()
	// Params is a string above and must be turned into a byte array to be sent
	// with http.Post
	byteParams := []byte(params)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(byteParams))

	// If there is an error making the POST request return the error
	if err != nil {
		return []byte{}, errors.Wrap(err, 1)
	}

	// Read the body into b, b will be a byte array representation of the
	// response
	b, err := ioutil.ReadAll(res.Body)

	// If we can't read the response return err
	if err != nil {
		return []byte{}, errors.Wrap(err, 1)
	}

	return b, nil
}

// GetFoodInfo ...
func GetFoodInfo(mmID int) {
	return
}

// urlBuilder is abstracted so that we can change the base url easily and so
// that we don't have to remember to add the nocache at the end
func urlBuilder() string {
	// noCache just needs to be a unique int so that their server doesn't return
	// the same value every time
	noCache := strconv.FormatInt(time.Now().UnixNano(), 10)
	return "http://nutrition.dartmouth.edu:8088/cwp?nocache=" + noCache
}

// AvailableSIDS gets the AvailableSIDs and returns them as a map with the keys
// being the sids options and the values being the display name for the sid eg:
// 	DDS: 53 Commons
//  CYC: Courtyard Cafe
func AvailableSIDS() (map[string]string, error) {

	availablesIDs := map[string]string{}
	// The JSON string copied from the Nutrition Website request
	params := constants.AvailableSIDSRequest

	b, err := makeRequest(params)
	if err != nil {
		return availablesIDs, errors.Wrap(err, 1)
	}

	// Create a struct to hold the response. This allows us to see whether the
	// response returned matches our expectations, if it doesn't then we want
	// to return early since we can't do anything with it anyway
	response := constants.AvailableSIDSResponse{}
	if err := json.Unmarshal(b, &response); err != nil {
		return availablesIDs, errors.Wrap(err, 1)
	}

	for _, sidArray := range response.Result.Result {
		// sidArray[0] currently holds the sid and sidArray[1] holds the display
		// name. This might change in the future but such is the nature of scrapers
		availablesIDs[sidArray[0]] = sidArray[1]
		// We can't check if it indeed holds anything so this will panic if it
		// can't read in the array.
	}

	// If the map is empty we know something went wrong so we return an error.
	if len(availablesIDs) == 0 {
		return availablesIDs, errors.Errorf("No new possible sids")
	}

	// If we made it this far then our map should contain key:value pairs
	// with the sid:displayName
	return availablesIDs, nil
}

// GetSID ...
func GetSID(sidKey string) (string, error) {

	params := fmt.Sprintf(constants.GetSIDSRequest, sidKey)
	b, err := makeRequest(params)

	if err != nil {
		return ``, errors.Wrap(err, 1)
	}

	sidResponse := constants.SIDResponse{}
	if err := json.Unmarshal(b, &sidResponse); err != nil {
		return ``, errors.Wrap(err, 1)
	}

	sid := sidResponse.Result.Sid
	if utf8.RuneCountInString(sid) == 0 {
		return ``, errors.Errorf("No SID found")
	}

	return sid, nil
}

// GetMenuList ...
func GetMenuList(sid string) (constants.MenuInfoSlice, error) {

	menuInfos := constants.MenuInfoSlice{}
	params := fmt.Sprintf(constants.GetMenuListRequest, sid)
	b, err := makeRequest(params)
	// If we can't read the response return err
	if err != nil {
		return menuInfos, errors.Wrap(err, 1)
	}

	menuList := constants.MenuListResponse{}
	if err := json.Unmarshal(b, &menuList); err != nil {
		return menuInfos, errors.Wrap(err, 1)
	}

	for _, v := range menuList.Result.MenusList {
		menu := constants.MenuInfo{
			ID:   int(v[0].(float64)),
			Name: v[3].(string),
		}
		menuInfos = append(menuInfos, menu)
	}
	return menuInfos, nil
}

// GetMealList ...
func GetMealList(sid string) (constants.MealInfoSlice, error) {
	params := fmt.Sprintf(constants.GetMealListRequest, sid)
	mealsList := constants.MealsListResponse{}
	b, err := makeRequest(params)
	// Will contain all of our meal info's
	mealInfoList := constants.MealInfoSlice{}
	if err != nil {
		return mealInfoList, errors.Wrap(err, 1)
	}

	if err := json.Unmarshal(b, &mealsList); err != nil {
		return mealInfoList, errors.Wrap(err, 1)
	}

	// This is a hack to get around the formatting of the response that we get
	// from the Dartmouth API
	// Basically we unmarshal the json into a map because the response looks like
	// { "1": ..., "2": ..., "n": ... } where n is variable...
	// and instead of being a list their API returns it as object with int keys
	// because we don't care about order we can unmarshal it into a map and loop
	// through the keys so that we don't have to have switch statements for each
	// menu. All in all this makes it so that there is less cognative overhead
	// at the price of having to use interface and type casting..
	// If the type conversion fails we return an error to remind programmer to
	// check the format of the api response
	for _, value := range mealsList.Result.MealsList {
		// pretty.Println(value)
		id, ok := value.([]interface{})[0].(float64)
		if !ok {
			return mealInfoList, errors.Errorf("Format of MealsList is incorrect.")
		}
		intID := int(id)

		name, ok := value.([]interface{})[2].(string) // string
		if !ok {
			return mealInfoList, errors.Errorf("Format of MealsList is incorrect.")
		}

		code, ok := value.([]interface{})[4].(string) // string
		if !ok {
			return mealInfoList, errors.Errorf("Format of MealsList is incorrect.")
		}

		start, ok := value.([]interface{})[5].(float64)
		if !ok {
			return mealInfoList, errors.Errorf("Format of MealsList is incorrect.")
		}
		intStart := int(start)

		end, ok := value.([]interface{})[6].(float64)
		if !ok {
			return mealInfoList, errors.Errorf("Format of MealsList is incorrect.")
		}
		intEnd := int(end)

		mealInfo := constants.MealInfo{intID, intStart, intEnd, name, code}
		mealInfoList = append(mealInfoList, mealInfo)
	}

	// mealInfoList = append(mealInfoList, mealOne, mealTwo)

	return mealInfoList, nil
}

// GetRecipesMenuMealDate ...
func GetRecipesMenuMealDate(sid string, menuID, mealID int) (constants.RecipeInfoSlice, error) {
	params := fmt.Sprintf(constants.GetRecipesMenuMealDate, sid, menuID, mealID)
	recipes := constants.RecipeInfoSlice{}
	b, err := makeRequest(params)
	if err != nil {
		return recipes, errors.Wrap(err, 1)
	}

	response := constants.RecipeResponse{}

	if err := json.Unmarshal(b, &response); err != nil {
		return recipes, errors.Wrap(err, 1)
	}

	for _, recipeRaw := range response.Result.RecipeitemsList {

		name := recipeRaw[0].(string)
		category := recipeRaw[1].([]interface{})[0].(string)
		recipeID := int(recipeRaw[1].([]interface{})[3].(float64))
		recipeMMID := int(response.Result.MmID)
		recipeRank := int(recipeRaw[1].([]interface{})[4].(float64))
		recipe := constants.RecipeInfo{
			Name:     name,
			Category: category,
			ID:       recipeID,
			Rank:     recipeRank,
			MmID:     recipeMMID,
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// GetNutrients ...
func GetNutrients(sid string, r *constants.RecipeInfo) (*constants.RecipeInfo, error) {
	params := fmt.Sprintf(constants.GetNutrientsRequest,
		sid, r.MmID, r.ID, r.Rank)
	b, err := makeRequest(params)
	if err != nil {
		return r, errors.Wrap(err, 1)
	}
	response := constants.NutrientInfoResponse{}
	if err := json.Unmarshal(b, &response); err != nil {
		return r, errors.Errorf(string(b))
	}
	r.Nutrients = response
	return r, nil
}
func timeTrack(start time.Time, fn string) {
	elapsed := time.Since(start)
	fmt.Println(fn, "took", elapsed)
}
func main() {
	sids, err := AvailableSIDS()
	if err != nil {
		log.Fatal(err)
	}
	quit := make(chan bool)
	for key, value := range sids {
		fmt.Println("Starting", key, "at", time.Now())
		go func(k, v string) {
			defer timeTrack(time.Now(), k)
			defer func() {
				quit <- true
			}()
			info := constants.VenueInfo{}
			sid, err := GetSID(k)
			if err != nil {
				panic(err)
			}
			info.Venue = v
			info.Key = k
			info.SID = sid
			info.Menus, err = GetMenuList(sid)
			info.Meals, err = GetMealList(sid)
			for _, menu := range info.Menus {
				for _, meal := range info.Meals {
					newRecipes, err := GetRecipesMenuMealDate(sid, menu.ID, meal.ID)
					if err != nil {
						log.Println(err.(*errors.Error).ErrorStack())
						return
					}
					info.Recipes = append(info.Recipes, newRecipes...)
				}
			}

			for index := range info.Recipes {
				fmt.Println(k, index, len(info.Recipes))
				_, err := GetNutrients(info.SID, &info.Recipes[index])
				if err != nil {
					log.Println(err.(*errors.Error).ErrorStack())
				}
			}
			return
		}(key, value)
	}

	for c := 0; c < len(sids); {
		select {
		case <-quit:
			c++
		}
	}
	close(quit)

}
