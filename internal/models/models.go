package models

import (
	"gorm.io/datatypes"
	_ "gorm.io/datatypes"
)

type ScheduleModel struct {
	Name        string                       ` json:"name" gorm:"size:191;unique"`
	DisplayName string                       `json:"displayname"`
	TimeZone    string                       `json:"timezone"`
	Schedule    datatypes.JSONType[Schedule] `json:"schedule"`
}

type Schedule struct {
	Dtype   string  `json:"dtype"`
	Corder  bool    `json:"Corder"`
	Shape   []int   `json:"Shape"`
	NdArray [][]int `json:"__ndarray__"`
}

type Policy struct {
	Name         string                                 `gorm:"size:191;unique" json:"name"`
	DisplayName  string                                 `json:"displayname"`
	Projects     datatypes.JSONSlice[Project]           `json:"projects"`
	Tags         datatypes.JSONSlice[map[string]string] `json:"tags"`
	ScheduleName string                                 `json:"schedulename"`
	Schedule     ScheduleModel                          `gorm:"foreignKey:ScheduleName;references:Name"`
	Provider     string                                 `json:"provider,omitempty"`
}

type Project struct {
	Name          string `json:"name"`
	CredentialRef string `json:"credentialRef"`
}
