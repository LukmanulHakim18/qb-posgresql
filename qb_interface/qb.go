package qbinterface

import (
	"database/sql"

	qb "github.com/LukmanulHakim18/qb-posgresql/qb_implement"
)

type QueryBuilder interface {
	// Action
	Select(tableName string) *qb.QueryBuilder
	Insert(tableName string, entity interface{}) (Id int64, err error)
	Update(tableName string) (RowsAffected int64, err error)
	Delete(tableName string) (RowsAffected int64, err error)

	Raw(query string, args ...interface{}) *qb.QueryBuilder

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
	Get() (*sql.Rows, error)

	// returning RowsAffected and error
	Exec() (int64, error)

	//Scanner
	ScanRow(rows *sql.Rows, obj interface{}) error
	ScanRows(rows *sql.Rows, obj interface{}) error

	// Manipulation Connection
	TrxBegin() error
	TrxRollback() error
	TrxCommit() error
	Close()

	// Scanner
	ScanEntity(entity interface{}) *qb.QueryBuilder
}
