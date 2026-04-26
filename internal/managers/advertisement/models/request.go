package models

type CreateAdvertisementRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Quantity    int      `json:"quantity"`
	Photos      []string `json:"photos"`
}
