package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jesusrmoreno/parse"
)

// Notification Object
type Notification struct {
	RecipeID int
	Name     string
	Day      int
	Month    int
	Year     int
	Seen     bool
	OnDate   time.Time `json:"onDate"`
	MenuName string    `json:"menuName"`
	MealName string    `json:"mealName"`
	Venue    string    `json:"venueKey"`
}

// DateObject ...
type DateObject struct {
	Type string    `json:"__type"`
	ISO  time.Time `json:"iso"`
}

// ParseNotification object
type ParseNotification struct {
	UUID     string     `json:"uuid"`
	ID       string     `json:"objectId"`
	Class    string     `json:"-"`
	RecipeID int        `json:"recipeID"`
	Name     string     `json:"recipeName"`
	Day      int        `json:"day"`
	Month    int        `json:"month"`
	Year     int        `json:"year"`
	Seen     bool       `json:"seen"`
	For      CreatedBy  `json:"for"`
	OnDate   DateObject `json:"onDate"`
	MenuName string     `json:"menuName"`
	MealName string     `json:"mealName"`
	Venue    string     `json:"venueKey"`
	Created  time.Time  `json:"createdAt"`
}

// GenerateUUID gives UUID
func (o ParseNotification) GenerateUUID() string {
	uuidStr := fmt.Sprintf("%d%d%d%s%s%s%s", o.Month, o.Day, o.Year, o.MenuName, o.MealName, o.Venue, o.For.ObjectID)
	return GetMD5Hash(uuidStr)
}

// SetID ...
func (o ParseNotification) SetID(id string) parse.Object {
	o.ID = id
	return o
}

// SetClass ...
func (o ParseNotification) SetClass(class string) parse.Object {
	o.Class = class
	return o
}

// JSON ...
func (o ParseNotification) JSON() (string, error) {
	j, err := json.Marshal(o)
	return string(j), err
}

// ClassName ...
func (o ParseNotification) ClassName() string {
	return o.Class
}

// ObjectID ...
func (o ParseNotification) ObjectID() string {
	return o.ID
}

// CreatedAt ...
func (o ParseNotification) CreatedAt() time.Time {
	return o.Created
}
