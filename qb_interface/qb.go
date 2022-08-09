package qbinterface

import (
	"database/sql"

	qb "github.com/LukmanulHakim18/qb-posgresql/qb_implement"
)

type QueryBuilder interface {
	// ection
	Select(tableName string) *qb.QueryBuilder
	Update(tableName string, entity interface{}) error
	Insert(tableName string, entity interface{}) (int64, error)
	Delete(tableName string, entity interface{}) error
	Raw(query string, args ...interface{}) (*sql.Rows, error)

	// condition
	Where(columnName string, symbol string, value interface{}) *qb.QueryBuilder
	OrWhere(columnName string, symbol string, value interface{}) *qb.QueryBuilder
	WhereIn(columnName string, value ...interface{}) *qb.QueryBuilder
	WhereNotIn(columnName string, value ...interface{}) *qb.QueryBuilder

	OrderBy(columnName, orderBy string) *qb.QueryBuilder
	// Deprecated: GroupBy is deprecated.
	GroupBy(columnName ...string) *qb.QueryBuilder
	Limit(int) *qb.QueryBuilder
	Offset(int64) *qb.QueryBuilder

	// Execute
	FindOne(entity interface{}) error
	Find(entities interface{}) error

	//Scanner
	ScanRow(rows *sql.Rows, obj interface{}) error
	ScanRows(rows *sql.Rows, obj interface{}) error

	//Close Connection
	Close()
}
