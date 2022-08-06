package querybuilder

import (
	"database/sql"
	"os"
	"strconv"
)

type QueryBuilder interface {
	// ection
	Select(tableName string) *queryBuilder
	Update(tableName string, entity interface{}) error
	Insert(tableName string, entity interface{}) (int64, error)
	Delete(tableName string, entity interface{}) error
	Raw(query string, args ...interface{}) (*sql.Rows, error)

	// condition
	Where(columnName string, symbol string, value interface{}) *queryBuilder
	OrWhere(columnName string, symbol string, value interface{}) *queryBuilder
	WhereIn(columnName string, value ...interface{}) *queryBuilder
	WhereNotIn(columnName string, value ...interface{}) *queryBuilder

	// OrWhere(columnName string, symbol string, value interface{}) *queryBuilder
	OrderBy(columnName, orderBy string) *queryBuilder

	// Execute
	FindOne(entity interface{}) error
	Find(entities interface{}) error

	//Scanner
	ScanRow(rows *sql.Rows, obj interface{}) error
	ScanRows(rows *sql.Rows, obj interface{}) error

	//Close Connection
	Close()
}

func NewQueryBuilder(db *sql.DB) QueryBuilder {
	return &queryBuilder{
		DBCon:     db,
		DebugMode: getBoolEnv("DEBUG_MODE"),
	}
}
func getStrEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		return ""
	}
	return val
}
func getBoolEnv(key string) bool {
	val := getStrEnv(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return ret
}
