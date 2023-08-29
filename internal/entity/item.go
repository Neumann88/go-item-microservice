package entity

type Item struct {
	ID    int `json:"id" db:"id"`
	Value int `json:"value" db:"value"`
}
