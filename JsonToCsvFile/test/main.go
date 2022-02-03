package test

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type EmployeeDetails struct {
	EmpId     string      `json:"emp_id" bson:"emp_id"`
	FirstName string      `json:"first_name" bson:"first_name"`
	LastName  string      `json:"last_name" bson:"last_name"`
	EmailId   string      `json:"email_id" bson:"email_id"`
	Password  string      `json:"password" bson:"password"`
	StartDate interface{} `json:"start_date" bson:"start_date"`
	EndDate   interface{} `json:"end_date" bson:"end_date"`
}

//func main() {
//	jsonData := getDataFromDB()
//
//	//// Unmarshal JSON data
//	//var jsonData []EmployeeDetails
//	//err = json.Unmarshal([]byte(jsonDataFromFile), &jsonData)
//	//
//	//if err != nil {
//	//	fmt.Println(err)
//	//}
//	//
//	csvFile, err := os.Create("./data.csv")
//
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer csvFile.Close()
//
//	writer := csv.NewWriter(csvFile)
//
//	for key, usance := range jsonData {
//		fmt.Println(key)
//		var row []string
//		row = append(row, usance.EmpId)
//		row = append(row, usance.FirstName)
//		row = append(row, usance.LastName)
//		writer.Write(row)
//	}
//
//	// remember to flush!
//	writer.Flush()
//}

var empDB = DbConnect()

func GetDataFromDB() []interface{} {
	cursor, err := empDB.Find(context.TODO(), bson.D{})

	var employees []interface{}

	var emp map[string]string

	if err != nil {
		log.Println("Failed To fetch Data From DB")
		defer cursor.Close(context.TODO())
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&emp)
		if err != nil {
			log.Println("Failed To Decode Data .............", err)
		}
		employees = append(employees, emp)
	}
	return employees
}

var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")

//DbConnect
//@note This method is to establish connection with mongoDb
//@return DBCollection pointer/
func DbConnect() *mongo.Collection {
	log.Println("dbConnection.dbConnect() Connecting to DB.....")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Failed To Connect To DB")
	}
	db := client.Database("goapi")
	employeeCollection := db.Collection("books")
	return employeeCollection
}
