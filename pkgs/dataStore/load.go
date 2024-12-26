package dataStore

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func idFromPath(path string, prefix string) (string, error) {
	sp := strings.Split(path, "/"+prefix+"/")
	if len(sp) <= 1 {
		return "", fmt.Errorf("invalid path")
	}
	rp := sp[len(sp)-1]
	if strings.HasSuffix(rp, ".yml") {
		return rp[:len(rp)-4], nil
	}
	if strings.HasSuffix(rp, ".yaml") {
		return rp[:len(rp)-5], nil
	}
	return rp, nil
}

func LoadFolder(path string) (DataStore, error) {
	var err error
	ds := DataStore{}
	if err = filepath.WalkDir(path+"/locations/", ds.loadLocations); err != nil {
		slog.Error(err.Error())
		return DataStore{}, err
	}
	if err = filepath.WalkDir(path+"/organizers/", ds.loadOrganizers); err != nil {
		slog.Error(err.Error())
		return DataStore{}, err
	}
	if err = filepath.WalkDir(path+"/events/", ds.loadEvents); err != nil {
		slog.Error(err.Error())
		return DataStore{}, err
	}
	return ds, nil
}

func (ds *DataStore) loadLocations(path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	fc, err := os.ReadFile(path)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	item := make(map[string]any)
	if err = yaml.Unmarshal([]byte(fc), &item); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	switch item["kind"] {
	case "location.openevents.tech/v1alpha1":
		l := Location{
			Name:       item["name"].(string),
			Country:    item["country"].(string),
			Region:     item["region"].(string),
			PostalCode: item["postalCode"].(string),
			Locality:   item["locality"].(string),
			Address:    item["address"].(string),
		}
		id, err := idFromPath(path, "locations")
		if err != nil {
			return err
		}
		if err := ds.AddLocation(id, l); err != nil {
			return err
		}
	default:
		err = fmt.Errorf("invalid kind")
	}
	return nil
}

func (ds *DataStore) loadOrganizers(path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	fc, err := os.ReadFile(path)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	item := make(map[string]any)
	if err = yaml.Unmarshal([]byte(fc), &item); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	switch item["kind"] {
	case "organizer.openevents.tech/v1alpha1":
		o := Organizer{
			Name: item["name"].(string),
		}
		if v, ok := item["url"].(string); ok {
			o.URL = v
		}
		id, err := idFromPath(path, "organizers")
		if err != nil {
			return err
		}
		if err := ds.AddOrganizer(id, o); err != nil {
			return err
		}
	default:
		err = fmt.Errorf("invalid kind")
	}
	return nil
}

func (ds *DataStore) loadEvents(path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	fc, err := os.ReadFile(path)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	item := make(map[string]any)
	if err = yaml.Unmarshal([]byte(fc), &item); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	switch item["kind"] {
	case "event.openevents.tech/v1alpha1":
		e := Event{
			Name: item["name"].(string),
			URL:  item["url"].(string),
		}
		if v, ok := item["description"].(string); ok {
			e.Description = v
		}
		if v, ok := item["organizer"].(string); ok {
			if _, exists := ds.Organizers[v]; !exists {
				return fmt.Errorf("%s is an invalid organizer", v)
			}
			e.OrganizerID = v
		}
		if d, ok := item["startDate"].(string); ok {
			var st time.Time
			var err error
			if t, ok := item["startTime"].(string); ok {
				st, err = time.Parse("2006-01-02 15:04:05 -0700", d+" "+t)
				if err != nil {
					return fmt.Errorf("%s %s is an invalid date", d, t)
				}
			} else {
				st, err = time.Parse("2006-01-02", d)
				if err != nil {
					return fmt.Errorf("%s is an invalid date", d)
				}
			}
			e.Start = st
		} else {
			return fmt.Errorf("startDate is required")
		}
		if d, ok := item["endDate"].(string); ok {
			var et time.Time
			var err error
			if t, ok := item["endTime"].(string); ok {
				et, err = time.Parse("2006-01-02 15:04:05 -0700", d+" "+t)
				if err != nil {
					return fmt.Errorf("%s %s is an invalid date", d, t)
				}
			} else {
				et, err = time.Parse("2006-01-02", d)
				if err != nil {
					return fmt.Errorf("%s is an invalid date", d)
				}
			}
			e.End = et
		} else {
			return fmt.Errorf("endDate is required")
		}
		if v, ok := item["format"].(string); ok {
			switch v {
			case "in-person":
				e.Format = InPerson
			case "hybrid":
				e.Format = Hybrid
			case "virtual":
				e.Format = Virtual
			default:
				return fmt.Errorf("%s is an invalid format", v)
			}
		}
		if v, ok := item["location"].(string); ok {
			if _, exists := ds.Locations[v]; !exists {
				return fmt.Errorf("%s is an invalid location", v)
			}
			e.LocationID = v
		}
		if cfp, ok := item["cfp"].(map[string]any); ok {
			fmt.Println(cfp)
			if v, ok := cfp["url"].(string); ok {
				e.CfP.URL = v
			}
			if v, ok := cfp["from"].(string); ok {
				t, err := time.Parse("2006-01-02", v)
				if err != nil {
					return fmt.Errorf("%s is an invalid date", v)
				}
				e.CfP.From = t
			}
			if v, ok := cfp["to"].(string); ok {
				t, err := time.Parse("2006-01-02", v)
				if err != nil {
					return fmt.Errorf("%s is an invalid date", v)
				}
				e.CfP.To = t
			}
		}
		if v, ok := item["topics"].([]string); ok {
			e.Topics = v
		}
		id, err := idFromPath(path, "events")
		if err != nil {
			return err
		}
		if err := ds.AddEvent(id, e); err != nil {
			return err
		}
	default:
		err = fmt.Errorf("invalid kind")
	}
	return nil
}
