package models

type Photo struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	SortOrder int    `json:"sort_order"`
	URL       string `json:"url"`
}

type User struct {
	BirthDate string `json:"birth_date,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	ID        string `json:"id,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateAdvertisementResponse struct {
	CreatedAt   string  `json:"created_at"`
	Description string  `json:"description"`
	ID          string  `json:"id"`
	Photos      []Photo `json:"photos"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Title       string  `json:"title"`
	UpdatedAt   string  `json:"updated_at"`
	UserID      string  `json:"user_id"`
}

type GetAdvertisementResponse struct {
	CreatedAt   string  `json:"created_at"`
	Description string  `json:"description"`
	ID          string  `json:"id"`
	Photos      []Photo `json:"photos"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Title       string  `json:"title"`
	UpdatedAt   string  `json:"updated_at"`
	User        User    `json:"user"`
	UserID      string  `json:"user_id"`
}
