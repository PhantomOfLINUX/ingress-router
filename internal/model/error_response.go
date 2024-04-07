package model

type ErrorResponse struct {
	Response   string `json:"response"`
	Details    string `json:"details"`
	Error      string `json:"error"`
	StatusCode int    `json:"-"`
}