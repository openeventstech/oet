package dataStore

import (
	"fmt"
)

type Location struct {
	Name       string
	Country    string
	Region     string
	PostalCode string
	Locality   string
	Address    string
}

func (ds *DataStore) AddLocation(id string, item Location) error {
	if ds.Locations == nil {
		ds.Locations = map[string]Location{}
	}
	if _, exists := ds.Locations[id]; exists {
		return fmt.Errorf("location %s already exists", id)
	}
	ds.Locations[id] = item
	return nil
}
