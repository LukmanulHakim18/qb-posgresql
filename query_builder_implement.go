package querybuilder

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

type queryBuilder struct {
	DBCon             *sql.DB       // menampung connection DB
	DebugMode         bool          // jika true akan menampilkan query ketika eksekusi, tambahkan pada .env DEBUG_MODE = true
	Action            string        // menampung Action
	Args              []interface{} //menampung seluruh argument
	Columns           []Column      // menampung data set untuk update, select atau insert
	Conditions        []condition   // menampung data kondisi
	TableName         string        // menampung table name
	LimitVal          int           // menampung Limit
	OffsetValue       int           // menampung Order By value
	OrderByConditions []orderBy     // menampung Order By
	PrimeryKey        interface{}   // menampung data Primary Key
}

// ======================================= DB Section =======================================
func (qb *queryBuilder) Close() {
	qb.DBCon.Close()
}

// ======================================= Action Section =======================================
func (qb *queryBuilder) Select(tableName string) *queryBuilder {
	qb.TableName = tableName
	qb.Action = ACTION_SELECT
	return qb
}

func (qb *queryBuilder) Update(tableName string, entity interface{}) error {
	qb.TableName = tableName
	qb.Action = ACTION_UPDATE

	qb.ScanEntity(entity)
	if qb.PrimeryKey == nil {
		return errors.New("tags id not found in this entity")
	}
	qb.Where("id", "=", qb.PrimeryKey)

	queryUpdate := qb.updateToString()
	querySet := qb.mappingDataSetUpdste()
	queryWhere := qb.conditionToString()
	queryFull := fmt.Sprintf("%s %s %s", queryUpdate, querySet, queryWhere)
	return qb.execute(queryFull)
}
func (qb *queryBuilder) Insert(tableName string, entity interface{}) (int64, error) {
	qb.TableName = tableName
	qb.Action = ACTION_INSERT

	qb.ScanEntity(entity)

	queryUpdate := qb.insertToString()
	querySet := qb.mappingDataSetInsert()
	queryFull := fmt.Sprintf("%s %s", queryUpdate, querySet)

	primaryKey, err := qb.save(queryFull)
	return primaryKey.(int64), err
}

func (qb *queryBuilder) Delete(tableName string, entity interface{}) error {
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

func (qb *queryBuilder) Raw(query string, args ...interface{}) (*sql.Rows, error) {
	for _, v := range args {
		qb.Args = append(qb.Args, v)
	}
	res, err := qb.executeRaw(query)
	return res, err
}

// ======================================= Condition Section =======================================

// Where make condition AND in query
func (qb *queryBuilder) Where(columnName string, symbol string, value interface{}) *queryBuilder {
	where := condition{
		ColumnName: columnName,
		Symbol:     symbol,
		Value:      value,
		Connector:  CONDITION_WHERE_AND,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

// Where make condition OR in query
func (qb *queryBuilder) OrWhere(columnName string, symbol string, value interface{}) *queryBuilder {
	where := condition{
		ColumnName: columnName,
		Symbol:     symbol,
		Value:      value,
		Connector:  CONDITION_WHERE_OR,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

// Where make condition AND IN (...) in query
func (qb *queryBuilder) WhereIn(columnName string, value ...interface{}) *queryBuilder {
	where := condition{
		ColumnName: columnName,
		Value:      value,
		Connector:  CONDITION_WHERE_IN,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

// Where make condition AND NOT IN (...) in query
func (qb *queryBuilder) WhereNotIn(columnName string, value ...interface{}) *queryBuilder {
	where := condition{
		ColumnName: columnName,
		Value:      value,
		Connector:  CONDITION_WHERE_NOT_IN,
	}
	qb.Conditions = append(qb.Conditions, where)

	return qb
}

// Where make condition ORDER BY ASC/DESC in query
func (qb *queryBuilder) OrderBy(columnName, orderByKey string) *queryBuilder {
	orderByUpper := strings.ToUpper(orderByKey)
	if orderByUpper == ORDER_BY_ASC {
		orderBy := orderBy{
			ColumnName: columnName,
			Value:      ORDER_BY_ASC,
		}
		qb.OrderByConditions = append(qb.OrderByConditions, orderBy)
	}
	if orderByUpper == ORDER_BY_DESC {
		orderBy := orderBy{
			ColumnName: columnName,
			Value:      ORDER_BY_DESC,
		}
		qb.OrderByConditions = append(qb.OrderByConditions, orderBy)
	}
	return qb
}

// Where make condition LIMIT in query
func (qb *queryBuilder) Limit(limitVal int) *queryBuilder {
	qb.LimitVal = limitVal
	return qb
}
