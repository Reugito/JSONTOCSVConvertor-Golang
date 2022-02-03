package main

import (
	"Digital_Credit_API/JsonToCsvFile/test"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

func main() {

	fmt.Println("Json to CSV App Is running ....................")
	dataFromDB := test.GetDataFromDB()

	_, err := JSONTOCSVConverter("JsonData.json", "jsonData.csv")

	if err != nil {
		log.Println(err)
	} else {
		log.Println("SuccessFully Converted json data to csv ")
	}

	_, err = dbDataToCSV(dataFromDB, "DBData.csv")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("SuccessFully Converted json data to csv ")
	}

}

func dbDataToCSV(dbData []interface{}, newCSVFileName string) (bool, error) {
	return writeToCSVFile(dbData, newCSVFileName)
}

func JSONTOCSVConverter(jsonFileName, newCSVFileName string) (bool, error) {
	jsonDataFromFile, err := ioutil.ReadFile(jsonFileName)

	if err != nil {
		return false, errors.New("failed To Read jsonFile from the Path ")
	}

	var jsonData []interface{}

	err = json.Unmarshal([]byte(jsonDataFromFile), &jsonData)
	if err != nil {
		return false, errors.New("failed To Unmarshal the json Data ")
	}
	return writeToCSVFile(jsonData, newCSVFileName)
}

func writeToCSVFile(data []interface{}, newCSVFileName string) (bool, error) {

	csvFile, err := os.Create(newCSVFileName)

	if err != nil {
		return false, errors.New("failed to create new CSV file with given name ")
	}
	defer csvFile.Close()

	rowInMapFormat := data[0]

	jStr, _ := json.Marshal(rowInMapFormat)

	unorderedKeys, err := getUnorderedKeys(rowInMapFormat)
	if err != nil {
		return false, err
	}

	orderedKeys := getOrderedKeys(fmt.Sprintf("%s", jStr), unorderedKeys)

	orderedReflectKeys := getOrderedValues(reflect.ValueOf(rowInMapFormat).MapKeys(), orderedKeys)

	writer := csv.NewWriter(csvFile)
	writer.Write(orderedKeys)

	for _, rowData := range data {
		row, _ := getDataToWrite(rowData, orderedReflectKeys)
		writer.Write(row)
	}
	writer.Flush()

	return true, nil
}

func getOrderedKeys(data string, keys []string) []string {
	var orderedKeys []string

	split1 := strings.Split(data, ":")
	stringD := strings.Join(split1, ",")
	split1 = strings.Split(stringD, "{")
	stringD = strings.Join(split1, "")
	split1 = strings.Split(stringD, ",")

	for _, val := range split1 {
		for _, key := range keys {
			if strings.Trim(val, "\"") == fmt.Sprintf("%s", key) {
				orderedKeys = append(orderedKeys, key)
			}
		}
	}
	return orderedKeys
}

func getUnorderedKeys(m interface{}) ([]string, error) {
	v := reflect.ValueOf(m)

	if v.Kind() != reflect.Map {
		return nil, errors.New(" given data is not a type of Map ")
	} else {
		var keys []string
		for _, key := range v.MapKeys() {
			keys = append(keys, fmt.Sprintf("%s", key))
			fmt.Println("Unordered ", key)
		}

		return keys, nil
	}
}

func getDataToWrite(m interface{}, orderedKeys []reflect.Value) ([]string, error) {
	v := reflect.ValueOf(m)

	if v.Kind() != reflect.Map {
		return nil, errors.New(" given data is not a type of Map ")
	} else {
		var row []string

		for _, val := range orderedKeys {
			if fmt.Sprintf("%s", v.MapIndex(val)) == "<invalid reflect.Value>" {
				row = append(row, "")
			} else {
				row = append(row, fmt.Sprintf("%s", v.MapIndex(val)))
			}
		}
		return row, nil
	}
}

func getOrderedValues(values []reflect.Value, orderedKeys []string) []reflect.Value {
	var keys []reflect.Value
	for _, value := range orderedKeys {
		for _, key := range values {
			if fmt.Sprintf("%s", key) == value {
				keys = append(keys, key)
			}
		}
	}
	return keys
}
