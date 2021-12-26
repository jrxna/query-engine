package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
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

	//queryCollection := ht.Database.Database("hypertool").Collection("queryTemplate")
	fmt.Println(request.Format)

	db, error := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test")
	if error != nil {
		panic(error.Error())
	}
	defer db.Close()

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

	var result []interface{}

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

	ctx.IndentedJSON(http.StatusOK, result)
}
