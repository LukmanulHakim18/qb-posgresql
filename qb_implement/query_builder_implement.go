package qb_implement

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	ACTION_DELETE = "DELETE"
	ACTION_INSERT = "INSERT"
	ACTION_SELECT = "SELECT"
	ACTION_UPDATE = "UPDATE"
	ACTION_RAW    = "RAW"

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
	PrimeryKey        *Column       // accommodate Primary Key value
	Query             string        // accommodate Query Raw in transaction or not
}

// ======================================= DB Section =======================================
func (qb *QueryBuilder) Close() {
	qb.DBConnection.Close()
}

// start transaction
func (qb *QueryBuilder) TrxBegin() (err error) {
	qb.DBTransaction, err = qb.DBConnection.Begin()
	if err != nil {
		return err
	}
	if qb.DebugMode {
		fmt.Println("transaction begin")
	}
	return nil
}

// Rolback Transaction and close transaction
func (qb *QueryBuilder) TrxRollback() error {
	if qb.DBTransaction == nil {
		fmt.Println("not in transaction")
		return nil
	}
	if err := qb.DBTransaction.Rollback(); err != nil {
		if qb.DebugMode {
			fmt.Println("rollback failed")
		}
		return err
	}
	qb.DBTransaction = nil
	if qb.DebugMode {
		fmt.Println("rollback transaction")
	}
	return nil
}

// Commit transaction and close transaction
func (qb *QueryBuilder) TrxCommit() error {
	if qb.DBTransaction == nil {
		fmt.Println("not in transaction")
		return nil
	}
	if err := qb.DBTransaction.Commit(); err != nil {
		if qb.DebugMode {
			fmt.Println("commit failed")
		}
		return err
	}
	qb.DBTransaction = nil
	if qb.DebugMode {
		fmt.Println("commit transaction")
	}
	return nil
}

// ======================================= Action Section =======================================

// Make query select * from table_name
func (qb *QueryBuilder) Select(tableName string) *QueryBuilder {
	qb.TableName = tableName
	qb.Action = ACTION_SELECT
	return qb
}

// Make query INSERT by entity
func (qb *QueryBuilder) Insert(tableName string, entity interface{}) (int64, error) {
	qb.TableName = tableName
	qb.Action = ACTION_INSERT

	qb.ScanEntity(entity)

	primaryKey, err := qb.save()
	npk := int64(0)
	if primaryKey != nil {
		npk = primaryKey.(int64)
	}
	return npk, err
}

// Make query UPDATE by entity
func (qb *QueryBuilder) Update(tableName string) (int64, error) {
	qb.TableName = tableName
	qb.Action = ACTION_UPDATE

	if qb.PrimeryKey != nil {
		qb.Where(qb.PrimeryKey.Name, "=", qb.PrimeryKey.Value)
	}
	return qb.Exec()
}

// Make query DELETE by entity
func (qb *QueryBuilder) Delete(tableName string) (int64, error) {
	qb.TableName = tableName
	qb.Action = ACTION_DELETE

	if qb.PrimeryKey != nil {
		qb.Where(qb.PrimeryKey.Name, "=", qb.PrimeryKey.Value)
	}

	return qb.Exec()
}

// Make query RAW
// you can creat enything query with param ?, and pass argument after query
// returning *sql.Rows
func (qb *QueryBuilder) Raw(query string, args ...interface{}) *QueryBuilder {
	qb.Args = append(qb.Args, args...)
	qb.Query = query
	return qb
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

func (qb *QueryBuilder) Build() {
	queries := []string{} // set table and field
	switch qb.Action {
	case ACTION_SELECT:
		queries = append(queries, qb.selectToString())    // set table and field
		queries = append(queries, qb.conditionToString()) // add condition to string
		queries = append(queries, qb.groupByToString())   // add orderBy to string
		queries = append(queries, qb.orderByToString())   // add orderBy to string
		queries = append(queries, qb.limitToString())     // add limit to string
		queries = append(queries, qb.offsetToString())    // add offset to string
		qb.Query = strings.Join(queries, " ")             // merge all query
	case ACTION_UPDATE:
		queries = append(queries, qb.updateToString())
		queries = append(queries, qb.mappingDataSetUpdate())
		queries = append(queries, qb.conditionToString())
		qb.Query = strings.Join(queries, " ") // merge all query
	case ACTION_DELETE:
		queries = append(queries, qb.deleteToString())
		queries = append(queries, qb.conditionToString())
		qb.Query = strings.Join(queries, " ") // merge all query
	case ACTION_INSERT:
		queries = append(queries, qb.insertToString())
		queries = append(queries, qb.mappingDataSetInsert())
		qb.Query = strings.Join(queries, " ") // merge all query
	}
}
