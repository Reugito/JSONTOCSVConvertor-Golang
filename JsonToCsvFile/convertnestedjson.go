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

	_, err = JSONTOCSVConverter2("complexjsondata.json", "compjsonData.csv")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("SuccessFully Converted json data to csv ")
	}
}

// JSONTOCSVConverter2
//@note :	This method Will write the json data to given csv file
//@param :		jsonFileName, newCSVFileName
//				jsonFileName : The path of json File which contains Data
//				newCSVFileName : The csv file to store the converted json Data
//@return :		boolean data and error depending on status of operation
func JSONTOCSVConverter2(jsonFileName, newCSVFileName string) (bool, error) {
	log.Println("JSONTOCSVConverter()..........")
	if jsonFileName == "" || newCSVFileName == "" {
		return false, errors.New("fileName is Empty")
	}

	_, err := os.Stat(jsonFileName)
	if err != nil {
		return false, err
	}

	jsonDataFromFile, err := ioutil.ReadFile(jsonFileName)

	if err != nil {
		log.Println("JSONTOCSVConverter() Failed To read jsonFile from the Path...........")
		return false, errors.New("failed To Read jsonFile from the Path ")
	}

	var data []interface{}

	err = json.Unmarshal([]byte(jsonDataFromFile), &data)
	if err != nil {
		log.Println("JSONTOCSVConverter() failed To Unmarshal the json Data...........")
		return false, errors.New("failed To Unmarshal the json Data ")
	}

	if data == nil {
		log.Println("JSONTOCSVConverter() data is not available to write.............")
		return false, errors.New("data is not available to write")
	} else {

		csvFile, err := os.Create(newCSVFileName)

		if err != nil {
			log.Println("JSONTOCSVConverter() failed to create new CSV file with given name  .............")
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
		var rowValues []string
		for index, val := range data {
			if index == 0 {
				values := decodeJson(val.(map[string]interface{}), "")
				headers, err = getHeadersForCSV(values)
				if err != nil {
					return false, err
				}
				err = writer.Write(headers)
				if err != nil {
					return false, err
				}
				rowValues, err = getValuesForCSVRow(headers, values)
				if err != nil {
					return false, err
				}
				err = writer.Write(rowValues)
				if err != nil {
					return false, err
				}
			} else {
				values := decodeJson(val.(map[string]interface{}), "")
				rowValues, err = getValuesForCSVRow(headers, values)
				if err != nil {
					return false, err
				}
				err = writer.Write(rowValues)
				if err != nil {
					return false, err
				}
			}
		}
		writer.Flush()
		return true, nil
	}
}

//decodeJson	 nestedKeyName - pass "" while calling the function
//
//@note :	This method is used to create New Dictionary which contains newHeaders/newKey and data with respect to  json key
//				here newHeaders/newKey created by concatenating nested json keys with _ and appending respective values to created header
//@param :		jsonDataForCSVRow, nestedKeyName
//				jsonDataForCSVRow :	It contains jsonData which will be decoded and appended to new Row in CSV File
//				nestedKeyName :		It is newHeader/key which will be created automatically by reading nested json keys
//@return :		jsonData
func decodeJson(jsonDataForCSVRow map[string]interface{}, nestedKeyName string) map[string]interface{} {
	log.Println("decodeJson()..............")

	if jsonDataForCSVRow == nil {
		log.Println("decodeJson() ", errors.New("given data is nil"))
		return nil
	} else {
		var DataMap = make(map[string]interface{})

		for key, val := range jsonDataForCSVRow {
			//fmt.Println(reflect.TypeOf(val))
			switch vv := val.(type) {
			case map[string]interface{}:
				for k, v := range decodeJson(vv, nestedKeyName+key+"_") {
					DataMap[k] = v
				}
			case []interface{}:
				DataMap[nestedKeyName+key] = "null"
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

//getHeadersForCSV
//@note : 	This method is for fetching the headers for csv file from given jsonData
//@param : 	mapData - contains json data
//@return :	array of string headers
func getHeadersForCSV(mapData map[string]interface{}) ([]string, error) {
	log.Println("getHeadersForCSV()................")
	if mapData == nil {
		log.Println("getHeadersForCSV() jsonData is nil ................")
		return nil, errors.New("jsonData is nil")
	} else {
		var headers []string
		for key := range mapData {
			headers = append(headers, key)
		}
		return headers, nil
	}
}

//getValuesForCSVRow
//@note : 	This method is for fetching the row values for csv file from given jsonData
//@param : 	headers - csv headers, mapData - contains json data
//@return :	array of string values
func getValuesForCSVRow(headers []string, mapData map[string]interface{}) ([]string, error) {
	log.Println("getValuesForCSVRow().............")
	if headers == nil || mapData == nil {
		log.Println("getValuesForCSVRow() either of headers or jason data is nil .........")
		return nil, errors.New("either of headers or jason data is nil")
	} else {
		var row []string
		for _, header := range headers {
			for key, val := range mapData {
				if header == key {
					row = append(row, fmt.Sprintf("%s", val))
				}
			}
		}
		return row, nil
	}
}
