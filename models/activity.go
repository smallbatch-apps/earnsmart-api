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
	Type    ActivityType
	Message string
	UserID  uint `json:"user_id" gorm:"index"`
	User    User
}

func (activity Activity) MarshalJSON() ([]byte, error) {
	type Alias Activity

	return json.Marshal(&struct {
		ID      uint   `json:"id"`
		Type    string `json:"type"`
		Message string `json:"message"`
		Alias
	}{
		ID:      activity.ID,
		Type:    string(activity.Type),
		Message: activity.Message,
		Alias:   (Alias)(activity),
	})
}
