package models

import (
	"time"
)

type Status int64

const (
	Inactive Status = 0
	Active   Status = 1
)

type Hat struct {
	Id       int        `json:"id" db:"id"`
	UsedBy   *string    `json:"usedby" db:"usedby"`
	Priority int        `json:"priority" db:"priority"`
	Cleaning *time.Time `json:"cleaning" db:"cleaning"`
}

type Party struct {
	Id            string `json:"id" db:"id"`
	Status        int    `json:"status" db:"status"`
	HatsRequested int    `json:"hatsrequired" db:"hatsrequired"`
	Hats          []Hat  `json:"Hats" db:"Hats"`
}
