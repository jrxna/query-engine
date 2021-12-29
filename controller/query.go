package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

func objectFormatter(values []interface{}, columns []string) map[string]interface{} {
	row := make(map[string]interface{})
	for index, value := range values {
		encodedData := fmt.Sprint(value)
		if next, ok := strconv.ParseFloat(encodedData, 64); ok == nil {
			row[columns[index]] = next
		} else if booleanValue, ok := strconv.ParseBool(encodedData); ok == nil {
			row[columns[index]] = booleanValue
		} else if "string" == fmt.Sprintf("%T", encodedData) {
			row[columns[index]] = string(value.([]byte))
		} else {
			fmt.Printf("Failed for type %T of %v\n", encodedData, encodedData)
		}
	}
	return row
}

func columnFormatter(rows *sqlx.Rows, values []interface{}, columns []string, scanArgs []interface{}) map[string][]interface{} {
	result := make(map[string][]interface{})
	for rows.Next() {
		error := rows.Scan(scanArgs...)
		if error != nil {
			panic(error.Error())
			// TODO: Pass context here
		}
		for index, value := range values {
			encodedData := fmt.Sprint(value)
			if next, ok := strconv.ParseFloat(encodedData, 64); ok == nil {
				result[columns[index]] = append(result[columns[index]], next)
			} else if booleanValue, ok := strconv.ParseBool(encodedData); ok == nil {
				result[columns[index]] = append(result[columns[index]], booleanValue)
			} else if "string" == fmt.Sprintf("%T", encodedData) {
				result[columns[index]] = append(result[columns[index]], string(value.([]byte)))
			} else {
				fmt.Printf("Failed for type %T of %v\n", encodedData, encodedData)
			}
		}
	}

	return result
}

func rowFormatter(values []interface{}, columns []string) []interface{} {
	row := make([]interface{}, 0)
	for _, value := range values {
		encodedData := fmt.Sprint(value)
		if next, ok := strconv.ParseFloat(encodedData, 64); ok == nil {
			row = append(row, next)
		} else if booleanValue, ok := strconv.ParseBool(encodedData); ok == nil {
			row = append(row, booleanValue)
		} else if "string" == fmt.Sprintf("%T", encodedData) {
			row = append(row, string(value.([]byte)))
		} else {
			fmt.Printf("Failed for type %T of %v\n", encodedData, encodedData)
		}
	}
	return row
}

func (ht *QueryController) GetQueryResult(ctx *gin.Context) {
	var request Request
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

	if resource["type"] == "mysql" {
		mysqlConfig := resource["mysql"].(map[string]interface{})

		cfg := mysql.Config{
			User:   fmt.Sprint(mysqlConfig["databaseUserName"]),
			Passwd: fmt.Sprint(mysqlConfig["databasePassword"]),
			Net:    "tcp",
			Addr:   fmt.Sprint(mysqlConfig["host"]) + ":" + fmt.Sprint(mysqlConfig["port"]),
			DBName: fmt.Sprint(mysqlConfig["databaseName"]),
		}

		db, err := sqlx.Connect("mysql", cfg.FormatDSN())
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			panic(err.Error())
		}
		defer db.Close()

		rows, err := db.Queryx(fmt.Sprint(query["content"]), request.Variables...)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			panic(err.Error())
		}

		columns, err := rows.Columns()
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			panic(err.Error())
		}
		count := len(columns)
		if count == 0 {
			ctx.IndentedJSON(http.StatusAccepted, gin.H{
				"success": true,
			})
			return
		}

		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		/* If the requested format is column, result will be of type map[string][]interface{}. In all other cases,
		 * result will be an array of interfaces. These interfaces could be maps of rows or just arrays of rows.
		 * Since result of the type map[string][]interface{} needs to be populated in a slightly different manner,
		 * it has been separated by the following if-else block from the other two types of formats.
		 */
		if request.Format == "column" {
			var result map[string][]interface{}
			result = columnFormatter(rows, values, columns, scanArgs)
			ctx.IndentedJSON(http.StatusAccepted, gin.H{
				"success": true,
				"result":  result,
			})
		} else {
			var result []interface{}
			for rows.Next() {
				error := rows.Scan(scanArgs...)
				if error != nil {
					panic(error.Error())
				}
				if request.Format == "object" {
					result = append(result, objectFormatter(values, columns))
				} else if request.Format == "row" {
					result = append(result, rowFormatter(values, columns))
				}
			}
			ctx.IndentedJSON(http.StatusAccepted, gin.H{
				"success": true,
				"result":  result,
			})
		}
	}
}
