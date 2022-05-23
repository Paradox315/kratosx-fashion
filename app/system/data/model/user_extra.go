package model

type UserExtra struct {
	Address     string   `json:"address,omitempty"`
	Country     string   `json:"country,omitempty"`
	City        string   `json:"city,omitempty"`
	Birthday    string   `json:"birthday,omitempty"`
	Description string   `json:"description,omitempty"`
	Figures     []string `json:"figures,omitempty"`
}
