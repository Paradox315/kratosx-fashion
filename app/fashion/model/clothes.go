package model

import "time"

type Clothes struct {
	ItemId     string    `json:"ItemId"`
	IsHidden   bool      `json:"IsHidden"`
	Categories []string  `json:"Categories"`
	Timestamp  time.Time `json:"Timestamp"`
	Labels     []string  `json:"Labels"`
	Comment    string    `json:"Comment"`
}
