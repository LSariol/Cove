package database

import "time"

type Secret struct {
	Id           string
	Key          string
	Value        string
	Version      int
	TimesPulled  int
	DateAdded    time.Time
	LastModified time.Time
}
