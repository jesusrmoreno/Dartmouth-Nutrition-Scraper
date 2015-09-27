package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jesusrmoreno/nutrition-scraper/constants"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// AvailableSIDSResponse is the structure of the JSON we're expecting to get
// back when we query for the AvailableSIDs ie: DDS, NOVACK, etc..
type AvailableSIDSResponse struct {
	Err    string `json:"error"`
	ID     int    `json:"id"`
	Result struct {
		CWPVersion string     `json:"cwp_version"`
		Result     [][]string `json:"result"`
	} `json:"result"`
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

	url := urlBuilder()

	// The JSON string copied from the Nutrition Website request
	params := constants.AvailableSIDSRequest

	// Params is a string above and must be turned into a byte array to be sent
	// with http.Post
	byteParams := []byte(params)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(byteParams))

	// If there is an error making the POST request return the error
	if err != nil {
		return map[string]string{}, err
	}

	// Read the body into b, b will be a byte array representation of the
	// response
	b, err := ioutil.ReadAll(res.Body)

	// If we can't read the response return err
	if err != nil {
		return map[string]string{}, err
	}

	// Create a struct to hold the response. This allows us to see whether the
	// response returned matches our expectations, if it doesn't then we want
	// to return early since we can't do anything with it anyway
	response := AvailableSIDSResponse{}
	if err := json.Unmarshal(b, &response); err != nil {
		return map[string]string{}, err
	}
	if response.Err != `` {
		return map[string]string{}, errors.New(response.Err)
	}

	sidMap := map[string]string{}
	for _, sidArray := range response.Result.Result {
		// sidArray[0] currently holds the sid and sidArray[1] holds the display
		// name. This might change in the future but such is the nature of scrapers
		sidMap[sidArray[0]] = sidArray[1]
		// We can't check if it indeed holds anything so this will panic if it
		// can't read in the array.
	}

	// If the map is empty we know something went wrong so we return an error.
	if len(sidMap) == 0 {
		return map[string]string{}, errors.New("No new possible sids")
	}

	// If we made it this far then our map should contain key:value pairs
	// with the sid:displayName
	return sidMap, nil
}

// GetSID ...
func GetSID(sid string) (string, error) {

	params := fmt.Sprintf(constants.GetSIDSRequest, sid)
	url := urlBuilder()

	// The JSON string copied from the Nutrition Website request
	// params := constants.AvailableSIDSRequest

	// Params is a string above and must be turned into a byte array to be sent
	// with http.Post
	byteParams := []byte(params)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(byteParams))

	// If there is an error making the POST request return the error
	if err != nil {
		return ``, err
	}

	// Read the body into b, b will be a byte array representation of the
	// response
	b, err := ioutil.ReadAll(res.Body)

	// If we can't read the response return err
	if err != nil {
		return ``, err
	}
	log.Println(string(b))

	return ``, nil
}

func main() {
	sids, err := AvailableSIDS()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(sids)
	for key, _ := range sids {
		GetSID(key)
		time.Sleep(2 * time.Second)
		// log.Println(key, value)
	}

	// Get's the different service id ie: Novack and shit
	// params := `{"service":"","method":"get_available_sids","id":1,"params":[null,"{\"remoteProcedure\":\"get_available_sids\"}"]}`

	// Gets the sID
	// params :=

	// Gets the menu for the sID
	// params := `{"service":"","method":"get_webmenu_list","id":5,"params":[{"sid":"DDS.4ef6cc52093e2da095f926af6a241154"},"{\"remoteProcedure\":\"get_webmenu_list\"}"]}`

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
