package models

import (
	"time"
)

// We're making the errors a string because if there is an error we probably
// need a code rewrite anyway... No point in doing type assertions when all
// we need to know is whether or not an error exists.

// AvailableSIDSResponse is the structure of the JSON we're expecting to get
// back when we query for the AvailableSIDs ie: DDS, NOVACK, etc..
type AvailableSIDSResponse struct {
	Error  string `json:"error"`
	ID     int    `json:"id"`
	Result struct {
		CWPVersion string     `json:"cwp_version"`
		Result     [][]string `json:"result"`
	} `json:"result"`
}

// MenuListResponse ...
type MenuListResponse struct {
	Error  string `json:"error"`
	ID     int    `json:"id"`
	Result struct {
		MenusList [][]interface{} `json:"menus_list"`
	} `json:"result"`
}

// SIDResponse ...
type SIDResponse struct {
	Error  string `json:"error"`
	ID     int    `json:"id"`
	Result struct {
		Sid string `json:"sid"`
	} `json:"result"`
}

// MealsListResponse ...
type MealsListResponse struct {
	Error  string `json:"error"`
	ID     int    `json:"id"`
	Result struct {
		MealsList map[string]interface{} `json:"meals_list"`
	} `json:"result"`
}

// RecipeResponse ..
type RecipeResponse struct {
	Error  string `json:"error"`
	ID     int    `json:"id"`
	Result struct {
		MmID            int             `json:"mm_id"`
		RecipeitemsList [][]interface{} `json:"recipeitems_list"`
		CatList         [][]string      `json:"cat_list"`
	} `json:"result"`
}

// VenueInfo ...
type VenueInfo struct {
	SID     string          `json:"sid"`
	Venue   string          `json:"venue"`
	Key     string          `json:"key"`
	Menus   MenuInfoSlice   `json:"menus"`
	Meals   MealInfoSlice   `json:"meals"`
	Recipes RecipeInfoSlice `json:"recipes"`
}

// MenuInfo ...
type MenuInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MenuInfoSlice ...
type MenuInfoSlice []MenuInfo

// MealInfo ...
type MealInfo struct {
	ID        int    `json:"id"`
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
	Name      string `json:"name"`
	Code      string `json:"code"`
}

// MealInfoSlice ...
type MealInfoSlice []MealInfo

// RecipeInfo ...
type RecipeInfo struct {
	Name      string               `json:"name"`
	Category  string               `json:"category"`
	ID        int                  `json:"id"`
	Rank      int                  `json:"rank"`
	MmID      int                  `json:"mmId"`
	VenueSID  string               `json:"venueSID"`
	Date      time.Time            `json:"date"`
	MealID    int                  `json:"mealId"`
	MenuID    int                  `json:"menuId"`
	Nutrients NutrientInfoResponse `json:"nutrients"`
}

// RecipeInfoSlice ...
type RecipeInfoSlice []RecipeInfo

// NutrientInfoResponse ...
type NutrientInfoResponse struct {
	Error  interface{} `json:"error"`
	ID     int         `json:"id"`
	Result struct {
		VitaIu               string      `json:"vita_iu"`
		Vitb6P               string      `json:"vitb6_p"`
		Sfa                  string      `json:"sfa"`
		CalciumP             string      `json:"calcium_p"`
		Thiamin              string      `json:"thiamin"`
		MufaP                string      `json:"mufa_p"`
		Zinc                 string      `json:"zinc"`
		Vitc                 string      `json:"vitc"`
		Message              string      `json:"message"`
		FatP                 string      `json:"fat_p"`
		Fiberdtry            string      `json:"fiberdtry"`
		Title                string      `json:"title"`
		ServingSizeGrams     float64     `json:"serving_size_grams"`
		FolacinP             string      `json:"folacin_p"`
		Phosphorus           string      `json:"phosphorus"`
		NiacinP              string      `json:"niacin_p"`
		Vitb12               string      `json:"vitb12"`
		Potassium            string      `json:"potassium"`
		ServingSizeText      string      `json:"serving_size_text"`
		Fat                  string      `json:"fat"`
		FatransP             string      `json:"fatrans_p"`
		SugarsP              string      `json:"sugars_p"`
		RecipeID             int         `json:"recipe_id"`
		SfaP                 string      `json:"sfa_p"`
		Vitb12P              string      `json:"vitb12_p"`
		Success              bool        `json:"success"`
		VitaIuP              string      `json:"vita_iu_p"`
		Calcium              string      `json:"calcium"`
		Mufa                 string      `json:"mufa"`
		Iron                 string      `json:"iron"`
		Output               string      `json:"output"`
		CarbsP               string      `json:"carbs_p"`
		CholestrolP          string      `json:"cholestrol_p"`
		Sugars               string      `json:"sugars"`
		SodiumP              string      `json:"sodium_p"`
		ZincP                string      `json:"zinc_p"`
		VitcP                string      `json:"vitc_p"`
		RiboflavinP          string      `json:"riboflavin_p"`
		Protein              string      `json:"protein"`
		ProteinP             string      `json:"protein_p"`
		Vitb6                string      `json:"vitb6"`
		PufaP                string      `json:"pufa_p"`
		Fatrans              string      `json:"fatrans"`
		IronP                string      `json:"iron_p"`
		NutrientIdentifier   string      `json:"nutrient_identifier"`
		ServingSizeMls       interface{} `json:"serving_size_mls"`
		Niacin               string      `json:"niacin"`
		CalfatP              string      `json:"calfat_p"`
		FiberdtryP           string      `json:"fiberdtry_p"`
		ThiaminP             string      `json:"thiamin_p"`
		PhosphorusP          string      `json:"phosphorus_p"`
		Calfat               string      `json:"calfat"`
		Carbs                string      `json:"carbs"`
		CaloriesP            string      `json:"calories_p"`
		Cholestrol           string      `json:"cholestrol"`
		Sodium               string      `json:"sodium"`
		Calories             string      `json:"calories"`
		Riboflavin           string      `json:"riboflavin"`
		PotassiumP           string      `json:"potassium_p"`
		Folacin              string      `json:"folacin"`
		ServingsPerContainer string      `json:"servings_per_container"`
		Pufa                 string      `json:"pufa"`
	} `json:"result"`
}
