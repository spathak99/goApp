package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PropertyMap is a json type that assists with inserting jsons into postgres
type PropertyMap map[string]interface{}

// Value unmarshals the json data
func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Scan takes the raw data from the DB and transforms it into the desired type
func (p *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}

	return nil
}
