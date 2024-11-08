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
	OwnableModel `gorm:"index;uniqueIndex:idx_user_name"`
	Name         string      `json:"name" gorm:"uniqueIndex:idx_user_name"`
	Type         SettingType `json:"type"`
	Value        string      `json:"value"`
}

func (setting Setting) MarshalJSON() ([]byte, error) {
	type Alias Setting

	return json.Marshal(&struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(setting.Type),
		Alias: (Alias)(setting),
	})
}
