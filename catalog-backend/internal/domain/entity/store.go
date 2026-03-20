package entity

type StoreInfo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url,omitempty"`
	Features    []Feature `json:"features,omitempty"`
}

type Feature struct {
	ID          string `json:"id"` // для удобства редактирования на фронтенде
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
