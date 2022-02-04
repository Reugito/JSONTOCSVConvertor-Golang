package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func main() {

	//const data = `{"0": null, "1": 0, "2": "s", "3": [null, 0, "string", [], {}], "4": {"0": null, "1": 0, "2": "s", "3": [], "4": {}}}`

	jsonDataFromFile, err := ioutil.ReadFile("JsonData.json")
	fmt.Println(err)

	var data []map[string]interface{}

	err = json.Unmarshal([]byte(jsonDataFromFile), &data)
	fmt.Println(err)
	fmt.Println("data :", data[0])
	//keys, err := objectKeys([]byte(jsonDataFromFile[0]))
	//fmt.Println(keys, err)
}

func objectKeys(b []byte) ([]string, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	if t != json.Delim('{') {
		return nil, errors.New("expected start of object")
	}
	var keys []string
	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		if t == json.Delim('}') {
			return keys, nil
		}
		keys = append(keys, t.(string))
		if err := skipValue(d); err != nil {
			return nil, err
		}
	}
}
func skipValue(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	switch t {
	case json.Delim('['), json.Delim('{'):
		for {
			if err := skipValue(d); err != nil {
				if err == end {
					break
				}
				return err
			}
		}
	case json.Delim(']'), json.Delim('}'):
		return end
	}
	return nil
}

var end = errors.New("invalid end of array or object")
