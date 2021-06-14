package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PropertyMap is a json type that assists with inserting jsons into postgres
type PropertyMap map[string]interface{}

// Value unmarshals the json data
func (pMap PropertyMap) Value() (driver.Value, error) {
	unmarshalledJson, err := json.Marshal(pMap)
	return unmarshalledJson,err

}

// Scan takes the raw data from the DB and transforms it into the desired type
func (pMap *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var intFace interface{}
	err := json.Unmarshal(source, &intFace)
	if err != nil {
		return err
	}

	*pMap, ok = intFace.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}

	return nil
}
