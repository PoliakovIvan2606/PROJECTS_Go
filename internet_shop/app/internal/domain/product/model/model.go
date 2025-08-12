package model 

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"escription"`
	ImageID string `json:"image_id"`
	Price int `json:"price"`
	CurrencyID string `json:"currency_id"`
	Rating string `json:"rating"`
	CategoryID string `json:"category_id"`
	Specification string `json:"specification"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}


