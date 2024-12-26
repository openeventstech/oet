package dataStore

import (
	"fmt"
)

type Organizer struct {
	Name string
	URL  string
}

func (ds *DataStore) AddOrganizer(id string, item Organizer) error {
	if ds.Organizers == nil {
		ds.Organizers = map[string]Organizer{}
	}
	if _, exists := ds.Organizers[id]; exists {
		return fmt.Errorf("organizer %s already exists", id)
	}
	ds.Organizers[id] = item
	return nil
}
