package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"query-engine/model"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var queryResult = model.QueryResult{
	Name: "query_1",
	Result: map[string]interface{}{
		"name": "Joel",
	},
}

var queryText = "SELECT * FROM user;"

func GetQueryResult(c *gin.Context) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query(queryText)
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var result []interface{}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
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

	c.IndentedJSON(http.StatusOK, result)
}
