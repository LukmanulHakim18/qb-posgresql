package qb_implement

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	ACTION_DELETE = "DELETE"
	ACTION_INSERT = "INSERT"
	ACTION_SELECT = "SELECT"
	ACTION_UPDATE = "UPDATE"

	ORDER_BY_ASC  = "ASC"
	ORDER_BY_DESC = "DESC"
)

type QueryBuilder struct {
	DBConnection      *sql.DB       // accommodate connection DB
	DBTransaction     *sql.Tx       // accommodate if query in transaction checking
	DebugMode         bool          // accommodate debug config, if true then print query and param
	Action            string        // accommodate Action
	Args              []interface{} // accommodate all arguments for query
	Columns           []Column      // accommodate data set for update, select or insert
	Conditions        []Condition   // accommodate condition
	TableName         string        // accommodate table name
	LimitVal          int           // accommodate limit
	OffsetValue       int64         // accommodate offset
	OrderByConditions []OrderBy     // accommodate order by
	GroupByConditions []string      // accommodate order by
	PrimeryKey        interface{}   // accommodate Primary Key value
}

// ======================================= DB Section =======================================
func (qb *QueryBuilder) Close() {
	qb.DBConnection.Close()
}

func (qb *QueryBuilder) TrxBegin() {
	qb.DBTransaction, _ = qb.DBConnection.Begin()
}
func (qb *QueryBuilder) TrxRollback() {
	qb.DBTransaction.Rollback()
	qb.DBTransaction = nil
}
func (qb *QueryBuilder) TrxCommit() {
	qb.DBTransaction.Commit()
	qb.DBTransaction = nil
}

// ======================================= Action Section =======================================

// Make query select * from table_name
func (qb *QueryBuilder) Select(tableName string) *QueryBuilder {
	qb.TableName = tableName
	qb.Action = ACTION_SELECT
	return qb
}

// Make query UPDATE by entity
func (qb *QueryBuilder) Update(tableName string, entity interface{}) error {
	qb.TableName = tableName
	qb.Action = ACTION_UPDATE

	qb.ScanEntity(entity)
	if qb.PrimeryKey == nil {
		return errors.New("tags id not found in this entity")
	}
	qb.Where("id", "=", qb.PrimeryKey)

	queryUpdate := qb.updateToString()
	querySet := qb.mappingDataSetUpdate()
	queryWhere := qb.conditionToString()
	queryFull := fmt.Sprintf("%s %s %s", queryUpdate, querySet, queryWhere)
	return qb.execute(queryFull)
}

// Make query INSERT by entity
func (qb *QueryBuilder) Insert(tableName string, entity interface{}) (int64, error) {
	qb.TableName = tableName
	qb.Action = ACTION_INSERT

	qb.ScanEntity(entity)

	queryUpdate := qb.insertToString()
	querySet := qb.mappingDataSetInsert()
	queryFull := fmt.Sprintf("%s %s", queryUpdate, querySet)

	primaryKey, err := qb.save(queryFull)
	return primaryKey.(int64), err
}

// Make query DELETE by entity
func (qb *QueryBuilder) Delete(tableName string, entity interface{}) error {
	qb.TableName = tableName
	qb.Action = ACTION_DELETE

	qb.ScanEntity(entity)
	if qb.PrimeryKey == nil {
		return errors.New("tags id not found in this entity")
	}
	qb.Where("id", "=", qb.PrimeryKey)

	queryDelete := qb.deleteToString()
	queryWhere := qb.conditionToString()
	queryFull := fmt.Sprintf("%s %s", queryDelete, queryWhere)
	return qb.execute(queryFull)
}

// Make query RAW
// you can creat enything query with param ?, and pass argument after query
// returning *sql.Rows
func (qb *QueryBuilder) Raw(query string, args ...interface{}) (*sql.Rows, error) {
	qb.Args = append(qb.Args, args...)
	res, err := qb.executeRawQuery(query)
	return res, err
}

// ======================================= Condition Section =======================================

// Where make condition AND in query
func (qb *QueryBuilder) Where(columnName string, symbol string, value interface{}) *QueryBuilder {
	where := Condition{
		ColumnName: columnName,
		Symbol:     symbol,
		Value:      value,
		Connector:  CONDITION_WHERE_AND,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

// Where make condition OR in query
func (qb *QueryBuilder) OrWhere(columnName string, symbol string, value interface{}) *QueryBuilder {
	where := Condition{
		ColumnName: columnName,
		Symbol:     symbol,
		Value:      value,
		Connector:  CONDITION_WHERE_OR,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

// Where make condition AND IN (...) in query
func (qb *QueryBuilder) WhereIn(columnName string, value ...interface{}) *QueryBuilder {
	where := Condition{
		ColumnName: columnName,
		Value:      value,
		Connector:  CONDITION_WHERE_IN,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

// Where make condition AND NOT IN (...) in query
func (qb *QueryBuilder) WhereNotIn(columnName string, value ...interface{}) *QueryBuilder {
	where := Condition{
		ColumnName: columnName,
		Value:      value,
		Connector:  CONDITION_WHERE_NOT_IN,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

func (qb *QueryBuilder) WhereBetween(columnName string, val1, val2 interface{}) *QueryBuilder {
	value := []interface{}{val1, val2}
	where := Condition{
		ColumnName: columnName,
		Value:      value,
		Connector:  CONDITION_WHERE_BETWEEN,
	}
	qb.Conditions = append(qb.Conditions, where)
	return qb
}

// Where make condition ORDER BY ASC/DESC in query
func (qb *QueryBuilder) OrderBy(columnName, orderByKey string) *QueryBuilder {
	orderByUpper := strings.ToUpper(orderByKey)
	if orderByUpper == ORDER_BY_ASC {
		orderBy := OrderBy{
			ColumnName: columnName,
			Value:      ORDER_BY_ASC,
		}
		qb.OrderByConditions = append(qb.OrderByConditions, orderBy)
	}
	if orderByUpper == ORDER_BY_DESC {
		orderBy := OrderBy{
			ColumnName: columnName,
			Value:      ORDER_BY_DESC,
		}
		qb.OrderByConditions = append(qb.OrderByConditions, orderBy)
	}
	return qb
}

// Make condition LIMIT in query
func (qb *QueryBuilder) Limit(limitVal int) *QueryBuilder {
	qb.LimitVal = limitVal
	return qb
}

// Make condition LIMIT in query
func (qb *QueryBuilder) Offset(offsetVal int64) *QueryBuilder {
	qb.OffsetValue = offsetVal
	return qb
}

// Deprecated: GroupBy is deprecated.
func (qb *QueryBuilder) GroupBy(columnName ...string) *QueryBuilder {
	qb.GroupByConditions = append(qb.GroupByConditions, columnName...)
	return qb
}
