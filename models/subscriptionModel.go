package models

import "time"

// Subscription covers a food subscription for a given user...
type Subscription struct {
	User struct {
		Type      string `json:"__type"`
		ClassName string `json:"className"`
		ObjectID  string `json:"objectId"`
	} `json:"User"`
	CreatedAt time.Time `json:"createdAt"`
	ObjectID  string    `json:"objectId"`
	Recipes   []int     `json:"recipes"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SubscriptionSlice is a slice of subscriptions ...
type SubscriptionSlice []Subscription
