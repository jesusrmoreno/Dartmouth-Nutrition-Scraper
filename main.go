package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jesusrmoreno/nutrition-scraper/constants"
	"github.com/kr/pretty"
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
		return []byte{}, err
	}

	// Read the body into b, b will be a byte array representation of the
	// response
	b, err := ioutil.ReadAll(res.Body)

	// If we can't read the response return err
	if err != nil {
		return []byte{}, err
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
		return availablesIDs, err
	}

	// Create a struct to hold the response. This allows us to see whether the
	// response returned matches our expectations, if it doesn't then we want
	// to return early since we can't do anything with it anyway
	response := constants.AvailableSIDSResponse{}
	if err := json.Unmarshal(b, &response); err != nil {
		return availablesIDs, err
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
		return availablesIDs, errors.New("No new possible sids")
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
		return ``, err
	}

	sidResponse := constants.SIDResponse{}
	if err := json.Unmarshal(b, &sidResponse); err != nil {
		return ``, err
	}

	sid := sidResponse.Result.Sid
	if utf8.RuneCountInString(sid) == 0 {
		return ``, errors.New("No SID found")
	}
	return sid, nil
}

// GetMenuList ...
func GetMenuList(sid string) (string, error) {
	params := fmt.Sprintf(constants.GetMenuListRequest, sid)
	b, err := makeRequest(params)
	// If we can't read the response return err
	if err != nil {
		return ``, err
	}

	menuList := constants.MenuListResponse{}
	if err := json.Unmarshal(b, &menuList); err != nil {
		return ``, err
	}
	pretty.Println(menuList)
	return ``, nil
}

func main() {
	sids, err := AvailableSIDS()
	if err != nil {
		log.Fatal(err)
	}

	for key := range sids {
		fmt.Println()
		sid, err := GetSID(key)
		// If there is an error getting the SID we want to continue getting them
		// for others in case it is a one off thing.
		if err != nil {
			log.Println(err)
			continue
		}
		GetMenuList(sid)
		// time.Sleep(2 * time.Second)
		// log.Println(key, value)
	}
	// GetMenuList("Hello")

	// Gets the meal times
	// params := `{"service":"","method":"get_webmenu_meals_list","id":6,"params":[{"sid":"DDS.4ef6cc52093e2da095f926af6a241154"},"{\"remoteProcedure\":\"get_webmenu_meals_list\"}"]}`

	//
	// params := `{"service":"","method":"get_recipes_for_menumealdate","id":7,"params":[{"sid":"DDS.4ef6cc52093e2da095f926af6a241154"},"{\"menu_id\":\"27\",\"meal_id\":\"1\",\"remoteProcedure\":\"get_recipes_for_menumealdate\",\"day\":26,\"month\":9,\"year\":2015,\"use_menu_query\":true,\"order_by\":\"pubgroup-alpha\",\"cache\":true}"]}`

	//
	// params := `{"service":"","method":"get_nutrient_label_items","id":8,"params":[{"sid":"DDS.4ef6cc52093e2da095f926af6a241154"},"{\"remoteProcedure\":\"get_nutrient_label_items\",\"mm_id\":22445,\"recipe_id\":-752,\"mmr_rank\":200,\"rule\":\"fda|raw\",\"output\":\"dictionary\",\"options\":\"facts\",\"cache\":true,\"recdata\":null}"]}`

	//
	// params := `{"service":"","method":"get_recipe_sub_ingredients","id":9,"params":[{"sid":"DDS.4ef6cc52093e2da095f926af6a241154"},"{\"remoteProcedure\":\"get_recipe_sub_ingredients\",\"recipeId\":752}"]}`

	//
	// params := `{"service":"","method":"get_recipe_allergen_list","id":10,"params":[{"sid":"DDS.4ef6cc52093e2da095f926af6a241154"},"{\"remoteProcedure\":\"get_recipe_allergen_list\",\"recipeId\":752}"]}`
	// jsonStr := []byte(params)
	//
	// res, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	//
	// // res, _ := http.Post(url, "empty")
	// bytes, _ := ioutil.ReadAll(res.Body)
	// log.Println(string(bytes))

}
