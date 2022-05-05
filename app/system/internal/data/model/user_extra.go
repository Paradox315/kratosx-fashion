package model

type UserExtra struct {
	Address     string `json:"address"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Birthday    string `json:"birthday"`
	Description string `json:"description"`
}
