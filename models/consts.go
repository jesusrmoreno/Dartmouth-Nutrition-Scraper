package models

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5Hash ...
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// AvailableSIDSRequest is the Request String to be used in getting the
// Available sids
const AvailableSIDSRequest = `
  {
    "service": "",
    "method": "get_available_sids",
    "id": 1,
    "params": [
      null,
      "{ \"remoteProcedure\":\"get_available_sids\" }"
    ]
  }
`

// GetSIDSRequest ...
const GetSIDSRequest = `{
  "service": "",
  "method": "create_context",
  "id": 2,
  "params": ["%s"]
}`

// GetMenuListRequest ...
const GetMenuListRequest = `{
  "service": "",
  "method": "get_webmenu_list",
  "id": 5,
  "params": [{
    "sid":"%s"
  }, "{\"remoteProcedure\":\"get_webmenu_list\"}"]
}`

// GetMealListRequest ...
const GetMealListRequest = `{
  "service": "",
  "method": "get_webmenu_meals_list",
  "id": 6,
  "params": [{
  "sid":"%s"},
  "{\"remoteProcedure\":\"get_webmenu_meals_list\"}"]
}`

// RecipesForMenuMealDate ...
const RecipesMenuMealDate = `{
  "service": "",
  "method": "get_recipes_for_menumealdate",
  "id": 7,
  "params":[{
    "sid":"%s"},
    "{\"menu_id\":\"%d\",\"meal_id\":\"%d\",\"remoteProcedure\":\"get_recipes_for_menumealdate\",\"day\":%d,\"month\":%d,\"year\":%d,\"use_menu_query\":true,\"order_by\":\"pubgroup-alpha\",\"cache\":true}"]
}`

// GetNutrientsRequest ...
const GetNutrientsRequest = `{"service":"","method":"get_nutrient_label_items","id":8,"params":[{"sid":"%s"},"{\"remoteProcedure\":\"get_nutrient_label_items\",\"mm_id\":%d,\"recipe_id\":-%d,\"mmr_rank\":%d,\"rule\":\"fda|raw\",\"output\":\"dictionary\",\"options\":\"facts\",\"cache\":true,\"recdata\":null}"]}`
