package model

type User struct {
	Id            string   `json:"_id"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Description   string   `json:"description"`
	Organization  string   `json:"organization"`
	Gender        string   `json:"gender"`
	CountryCode   string   `json:"countryCode"`
	PictureURL    string   `json:"pictureURL"`
	EmailAddress  string   `json:"emailAddress"`
	EmailVerified bool     `json:"emailVerified"`
	Birthday      string   `json:"birthday"`
	Status        string   `json:"status"`
	Role          string   `json:"role"`
	Groups        []string `json:"groups"`
}
