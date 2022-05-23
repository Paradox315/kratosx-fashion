package do

type Clothes struct {
	Id          string   `json:"id,omitempty"`
	Type        string   `json:"type,omitempty"`
	Description string   `json:"description,omitempty"`
	Image       string   `json:"image,omitempty"`
	Brand       string   `json:"brand,omitempty"`
	Style       string   `json:"style,omitempty"`
	Region      string   `json:"region,omitempty"`
	Time        string   `json:"time,omitempty"`
	Price       float32  `json:"price,omitempty"`
	Colors      []string `json:"colors,omitempty"`
}

type ClothesComment struct {
	Description string `json:"description"`
	Price       string `json:"price"`
	Region      string `json:"region"`
	Image       string `json:"image"`
	Colors      string `json:"colors"`
}
