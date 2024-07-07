package models

import (
	"encoding/json"
)

type SettingType string

const (
	SettingTypeUser SettingType = "user"
	SettingTypeApp  SettingType = "app"
	SettingTypePerm SettingType = "perm"
)

type Setting struct {
	CustomModel
	Name   string
	Type   SettingType
	Value  string
	UserID uint `json:"user_id" gorm:"index"`
	User   User
}

func (setting Setting) MarshalJSON() ([]byte, error) {
	type Alias Setting

	return json.Marshal(&struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
		Alias
	}{
		ID:    setting.ID,
		Name:  setting.Name,
		Type:  string(setting.Type),
		Value: setting.Value,
		Alias: (Alias)(setting),
	})
}
