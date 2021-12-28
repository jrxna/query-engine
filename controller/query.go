package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

type QueryController struct {
	Database *mongo.Client
}

type Request struct {
	Name      string        `json:"name" validate:"min=1,max=128,regexp=^[a-zA-Z_][a-zA-Z_0-9-]+[a-zA-Z_0-9]$"`
	Variables []interface{} `json:"variables"`
	Format    string        `json:"format"`
}

var queryText = "SELECT * FROM user;"

func (ht *QueryController) GetQueryResult(ctx *gin.Context) {
	var request Request
	//var query model.Query
	if error := ctx.ShouldBindJSON(&request); error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		panic(error.Error())
	}
	if error := validator.Validate(request); error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		panic(error.Error())
	}
	if request.Format != "row" && request.Format != "column" && request.Format != "object" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "The format can only be one of the enumeration [row, column, object]",
		})
		panic("The format can only be one of enumeration [row, column, object]")
	}

	queryCollection := ht.Database.Database("hypertool").Collection("querytemplates")
	var query map[string]interface{}
	err := queryCollection.FindOne(ctx, bson.M{"name": request.Name}).Decode(&query)
	if err != nil {
		log.Fatal(err)
	}

	resourceCollection := ht.Database.Database("hypertool").Collection("resources")
	var resource map[string]interface{}
	err = resourceCollection.FindOne(ctx, bson.M{"_id": query["resource"]}).Decode(&resource)
	if err != nil {
		log.Fatal(err)
	}

	var result []interface{}

	if resource["type"] == "mysql" {
		mysqlConfig := resource["mysql"].(map[string]interface{})

		cfg := mysql.Config{
			User:   fmt.Sprint(mysqlConfig["databaseUserName"]),
			Passwd: fmt.Sprint(mysqlConfig["databasePassword"]),
			Net:    "tcp",
			Addr:   fmt.Sprint(mysqlConfig["host"]) + ":" + fmt.Sprint(mysqlConfig["port"]),
			DBName: fmt.Sprint(mysqlConfig["databaseName"]),
		}

		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		log.Println("Connected to the database successfully...")

		log.Println("Executing query...")
		log.Println(query["content"])
		rows, error := db.Query(queryText)
		if error != nil {
			panic(error.Error())
		}
		columns, error := rows.Columns()
		if error != nil {
			panic(error.Error())
		}
		count := len(columns)
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			error := rows.Scan(scanArgs...)
			if error != nil {
				panic(error.Error())
			}
			row := make(map[string]interface{})
			for index, value := range values {

				encodedData := value.([]byte)

				/**
				 * From the Go Blog: JSON and GO - 25 Jan 2011:
				 * The json package uses map[string]interface{} and []interface{} values to store arbitrary JSON objects and arrays;
				 * it will happily unmarshal any valid JSON blob into a plain interface{} value. The default concrete Go types are:
				 *
				 * bool for JSON booleans,
				 * float64 for JSON numbers,
				 * string for JSON strings, and
				 * nil for JSON null.
				 **/
				if next, ok := strconv.ParseFloat(string(encodedData), 64); ok == nil {
					row[columns[index]] = next
				} else if booleanValue, ok := strconv.ParseBool(string(encodedData)); ok == nil {
					row[columns[index]] = booleanValue
				} else if "string" == fmt.Sprintf("%T", string(encodedData)) {
					row[columns[index]] = string(encodedData)
				} else {
					fmt.Printf("Failed for type %T of %v\n", encodedData, encodedData)
				}
			}
			result = append(result, row)

		}
	}

	ctx.IndentedJSON(http.StatusOK, result)
}
