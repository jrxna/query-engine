package model

type Query struct {
	Id       string `json:"_id"`
	Name     string `json:"name"`
	Resource string `json:"resource"`
	Content  string `json:"content"`
}
