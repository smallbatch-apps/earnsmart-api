package models

import (
	"encoding/json"
)

type ActivityType string

const (
	ActivityTypeUser  ActivityType = "user"
	ActivityTypeAdmin ActivityType = "admin"
)

type Activity struct {
	CustomModel
	OwnableModel
	Type    ActivityType `json:"type"`
	Message string       `json:"message"`
}

func (activity Activity) MarshalJSON() ([]byte, error) {
	type Alias Activity
	return json.Marshal(&struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(activity.Type),
		Alias: (Alias)(activity),
	})
}
