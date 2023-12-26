package model

import "time"

type Urls []Url
type Url struct {
	ID     int       `json:"id"`
	Url    string    `json:"url"`
	Count  int       `json:"count"`
	Short  string    `json:"short"`
	Expire time.Time `json:"expire"`
}
