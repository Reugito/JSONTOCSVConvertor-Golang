package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {

	_, err := JSONTOCSVConverter2("JsonData.json", "jsonData.csv")

	if err != nil {
		log.Println(err)
	} else {
		log.Println("SuccessFully Converted json data to csv ")
	}

	JSONTOCSVConverter2("complexjsondata.json", "compjsonData.csv")
}

func decodeJson(m map[string]interface{}, nestedKeyName string) map[string]interface{} {
	if m == nil {
		log.Println(errors.New("given data is nil"))
		return nil
	} else {
		var DataMap = make(map[string]interface{})

		for key, val := range m {
			switch vv := val.(type) {
			case map[string]interface{}:
				for k, v := range decodeJson(vv, nestedKeyName+key+"_") {
					DataMap[k] = v
				}
			case string:
				DataMap[nestedKeyName+key] = vv
			case float64:
				DataMap[nestedKeyName+key] = strconv.FormatFloat(vv, 'f', -1, 64)
			case bool:
				DataMap[nestedKeyName+key] = strconv.FormatBool(vv)
			}
		}
		return DataMap
	}
}

func writeToCSVFile2(data []interface{}, newCSVFileName string) (bool, error) {
	if data == nil {
		return false, errors.New("data is not available to write")
	} else {

		csvFile, err := os.Create(newCSVFileName)

		if err != nil {
			return false, errors.New("failed to create new CSV file with given name ")
		}
		defer func(csvFile *os.File) {
			err := csvFile.Close()
			if err != nil {
				log.Println(err)
			}
		}(csvFile)

		writer := csv.NewWriter(csvFile)
		var headers []string
		for index, val := range data {
			if index == 0 {
				values := decodeJson(val.(map[string]interface{}), "")
				headers = getHeaders(values)
				err = writer.Write(headers)
				if err != nil {
					return false, err
				}
				err = writer.Write(getKeys(headers, values))
				if err != nil {
					return false, err
				}
			} else {
				values := decodeJson(val.(map[string]interface{}), "")
				err = writer.Write(getKeys(headers, values))
				if err != nil {
					return false, err
				}
			}
		}
		writer.Flush()
		return true, nil
	}
}

func getKeys(headers []string, mapData map[string]interface{}) []string {
	if headers == nil || mapData == nil {
		return nil
	} else {
		var row []string
		for _, header := range headers {
			for key, val := range mapData {
				if header == key {
					row = append(row, fmt.Sprintf("%s", val))
				}
			}
		}
		return row
	}
}

func getHeaders(mapData map[string]interface{}) []string {
	if mapData == nil {
		return nil
	} else {
		var headers []string
		for key := range mapData {
			headers = append(headers, key)
		}
		return headers
	}
}

func JSONTOCSVConverter2(jsonFileName, newCSVFileName string) (bool, error) {
	if jsonFileName == "" || newCSVFileName == "" {
		return false, errors.New("fileName is Empty")
	}

	_, err := os.Stat(jsonFileName)
	if err != nil {
		return false, err
	}

	jsonDataFromFile, err := ioutil.ReadFile(jsonFileName)

	if err != nil {
		return false, errors.New("failed To Read jsonFile from the Path ")
	}

	var jsonData []interface{}

	err = json.Unmarshal([]byte(jsonDataFromFile), &jsonData)
	if err != nil {
		return false, errors.New("failed To Unmarshal the json Data ")
	}
	return writeToCSVFile2(jsonData, newCSVFileName)
}
