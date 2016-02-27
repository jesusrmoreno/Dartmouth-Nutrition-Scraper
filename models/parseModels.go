package models

import (
	"encoding/json"
	"time"

	"github.com/jesusrmoreno/parse"
)

// ParseOffering ...
type ParseOffering struct {
	Venue    string `json:"venueKey"`
	ID       string `json:"objectId"`
	Day      int    `json:"day"`
	Month    int    `json:"month"`
	Year     int    `json:"year"`
	MenuName string `json:"menuName"`
	MealName string `json:"mealName"`
	Recipes  struct {
		Op      string   `json:"__op"`
		Objects []Object `json:"objects"`
	} `json:"recipes"`
	Class   string    `json:"-"`
	Created time.Time `json:"-"`
	UUID    string    `json:"uuid"`
}

// Object ...
type Object struct {
	Type      string `json:"__type"`
	Classname string `json:"className"`
	ObjectID  string `json:"objectId"`
}

// AddRecipe ...
func (o *ParseOffering) AddRecipe(objectID string) {
	o.Recipes.Op = "AddRelation"
	o.Recipes.Objects = append(o.Recipes.Objects, Object{
		Type:      "Pointer",
		Classname: "Recipe",
		ObjectID:  objectID,
	})
}

// ClassName ...
func (o ParseOffering) ClassName() string {
	return o.Class
}

// ObjectID ...
func (o ParseOffering) ObjectID() string {
	return o.ID
}

// CreatedAt ...
func (o ParseOffering) CreatedAt() time.Time {
	return o.Created
}

// SetID ...
func (o ParseOffering) SetID(id string) parse.Object {
	o.ID = id
	return o
}

// SetClass ...
func (o ParseOffering) SetClass(class string) parse.Object {
	o.Class = class
	return o
}

// JSON ...
func (o ParseOffering) JSON() (string, error) {
	j, err := json.Marshal(o)
	return string(j), err
}

// CreatedBy ....
type CreatedBy struct {
	Kind      string `json:"__type"`
	ClassName string `json:"className"`
	ObjectID  string `json:"objectId"`
}

// ParseRecipe ...
type ParseRecipe struct {
	ID          string               `json:"objectId"`
	Class       string               `json:"-"`
	Created     time.Time            `json:"createdAt"`
	Updated     time.Time            `json:"updatedAt"`
	Name        string               `json:"name"`
	Category    string               `json:"category"`
	DartmouthID int                  `json:"dartmouthId"`
	Rank        int                  `json:"rank"`
	UUID        string               `json:"uuid"`
	Nutrients   NutrientInfoResponse `json:"nutrients"`
	CreatedBy   CreatedBy            `json:"createdBy"`
}

// JSON ...
func (o ParseRecipe) JSON() (string, error) {
	j, err := json.Marshal(o)
	return string(j), err
}

// ClassName ...
func (o ParseRecipe) ClassName() string {
	return o.Class
}

// ObjectID ...
func (o ParseRecipe) ObjectID() string {
	return o.ID
}

// CreatedAt ...
func (o ParseRecipe) CreatedAt() time.Time {
	return o.Created
}

// SetID ...
func (o ParseRecipe) SetID(id string) parse.Object {
	o.ID = id
	return o
}

// SetClass ...
func (o ParseRecipe) SetClass(class string) parse.Object {
	o.Class = class
	return o
}

// ParseNutrients ...
type ParseNutrients struct {
	ID                   string    `json:"objectId"`
	UUID                 string    `json:"uuid"`
	Class                string    `json:"-"`
	Created              time.Time `json:"-"`
	DartmouthID          int       `json:"dartmouthId"`
	VitaIu               string    `json:"vita_iu"`
	Vitb6P               string    `json:"vitb6_p"`
	Sfa                  string    `json:"sfa"`
	CalciumP             string    `json:"calcium_p"`
	Thiamin              string    `json:"thiamin"`
	MufaP                string    `json:"mufa_p"`
	Zinc                 string    `json:"zinc"`
	Vitc                 string    `json:"vitc"`
	Message              string    `json:"message"`
	FatP                 string    `json:"fat_p"`
	Fiberdtry            string    `json:"fiberdtry"`
	Title                string    `json:"title"`
	ServingSizeGrams     float64   `json:"serving_size_grams"`
	FolacinP             string    `json:"folacin_p"`
	Phosphorus           string    `json:"phosphorus"`
	NiacinP              string    `json:"niacin_p"`
	Vitb12               string    `json:"vitb12"`
	Potassium            string    `json:"potassium"`
	ServingSizeText      string    `json:"serving_size_text"`
	Fat                  string    `json:"fat"`
	FatransP             string    `json:"fatrans_p"`
	SugarsP              string    `json:"sugars_p"`
	RecipeID             int       `json:"recipe_id"`
	SfaP                 string    `json:"sfa_p"`
	Vitb12P              string    `json:"vitb12_p"`
	Success              bool      `json:"success"`
	VitaIuP              string    `json:"vita_iu_p"`
	Calcium              string    `json:"calcium"`
	Mufa                 string    `json:"mufa"`
	Iron                 string    `json:"iron"`
	Output               string    `json:"output"`
	CarbsP               string    `json:"carbs_p"`
	CholestrolP          string    `json:"cholestrol_p"`
	Sugars               string    `json:"sugars"`
	SodiumP              string    `json:"sodium_p"`
	ZincP                string    `json:"zinc_p"`
	VitcP                string    `json:"vitc_p"`
	RiboflavinP          string    `json:"riboflavin_p"`
	Protein              string    `json:"protein"`
	ProteinP             string    `json:"protein_p"`
	Vitb6                string    `json:"vitb6"`
	PufaP                string    `json:"pufa_p"`
	Fatrans              string    `json:"fatrans"`
	IronP                string    `json:"iron_p"`
	NutrientIdentifier   string    `json:"nutrient_identifier"`
	Niacin               string    `json:"niacin"`
	CalfatP              string    `json:"calfat_p"`
	FiberdtryP           string    `json:"fiberdtry_p"`
	ThiaminP             string    `json:"thiamin_p"`
	PhosphorusP          string    `json:"phosphorus_p"`
	Calfat               string    `json:"calfat"`
	Carbs                string    `json:"carbs"`
	CaloriesP            string    `json:"calories_p"`
	Cholestrol           string    `json:"cholestrol"`
	Sodium               string    `json:"sodium"`
	Calories             string    `json:"calories"`
	Riboflavin           string    `json:"riboflavin"`
	PotassiumP           string    `json:"potassium_p"`
	Folacin              string    `json:"folacin"`
	ServingsPerContainer string    `json:"servings_per_container"`
	Pufa                 string    `json:"pufa"`
	Vegetarian           bool      `json:"vegetarian"`
	Gluten               bool      `json:"glutenFree"`
	Local                bool      `json:"local"`
	Kosher               bool      `json:"kosher"`
	Halal                bool      `json:"halal"`
	Vegan                bool      `json:"vegan"`
	Eggs                 bool      `json:"eggs"`
	Fish                 bool      `json:"fish"`
	Dairy                bool      `json:"dairy"`
	TreeNuts             bool      `json:"treeNuts"`
	Peanuts              bool      `json:"peanutes"`
	Pork                 bool      `json:"pork"`
	Soy                  bool      `json:"soy"`
	ShellFish            bool      `json:"shellFish"`
	Wheat                bool      `json:"wheat"`
}

// ClassName ...
func (o ParseNutrients) ClassName() string {
	return o.Class
}

// ObjectID ...
func (o ParseNutrients) ObjectID() string {
	return o.ID
}

// CreatedAt ...
func (o ParseNutrients) CreatedAt() time.Time {
	return o.Created
}

// SetID ...
func (o ParseNutrients) SetID(id string) parse.Object {
	o.ID = id
	return o
}

// SetClass ...
func (o ParseNutrients) SetClass(class string) parse.Object {
	o.Class = class
	return o
}

// JSON ...
func (o ParseNutrients) JSON() (string, error) {
	j, err := json.Marshal(o)
	return string(j), err
}
