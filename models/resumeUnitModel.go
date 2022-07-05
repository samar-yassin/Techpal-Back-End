package models

type resumeUnit struct {
	Name string   `json:"name"`
	Hide bool     `json:"hide"`
	Data []string `json:"data"`
}
