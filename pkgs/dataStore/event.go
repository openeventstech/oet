package dataStore

import (
	"fmt"
	"time"
)

type Format int

const (
	InPerson Format = iota
	Hybrid
	Virtual
)

type Event struct {
	Kind        string
	Name        string
	Description string
	URL         string
	OrganizerID string
	Start       time.Time
	End         time.Time
	Format      Format
	LocationID  string
	CfP         struct {
		URL  string
		From time.Time
		To   time.Time
	}
	Topics []string
}

func (ds *DataStore) AddEvent(id string, item Event) error {
	if ds.Events == nil {
		ds.Events = map[string]Event{}
	}
	if _, exists := ds.Locations[id]; exists {
		return fmt.Errorf("event %s already exists", id)
	}
	ds.Events[id] = item
	return nil
}
