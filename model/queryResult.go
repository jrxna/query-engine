package model

type QueryResult struct {
	Name   string                 `json:"name"`
	Result map[string]interface{} `json:"result"`
}
