package dataStore

type DataStore struct {
	Events     map[string]Event
	Locations  map[string]Location
	Organizers map[string]Organizer
}
