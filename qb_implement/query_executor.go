package qb_implement

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	util "github.com/LukmanulHakim18/qb-posgresql/utility"
)

// ======================================= Execute Section =======================================

//execute query

// this function for executing select with attachment query condition
func (qb *QueryBuilder) executeQuerySelect() (*sql.Rows, error) {
	queries := []string{qb.selectToString()}                             // set table and field
	queries = append(queries, qb.conditionToString())                    // add condition to string
	queries = append(queries, qb.groupByToString())                      // add orderBy to string
	queries = append(queries, qb.orderByToString())                      // add orderBy to string
	queries = append(queries, qb.limitToString())                        // add limit to string
	queries = append(queries, qb.offsetToString())                       // add offset to string
	queryFull := qb.ValidatedQueryAndMapping(strings.Join(queries, " ")) // merubah symbol ? -> $n

	util.DebugQueryAndParams(qb.DebugMode, queryFull, qb.Args) // for debuging
	if qb.DBTransaction != nil {
		rows, err := qb.DBTransaction.Query(queryFull, qb.Args...)
		if err != nil {
			return nil, err
		}
		return rows, nil
	}
	rows, err := qb.DBConnection.Query(queryFull, qb.Args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// return 1 value filter by condition
func (qb *QueryBuilder) FindOne(entity interface{}) error {
	defer qb.deferFunc()
	if qb.Action != ACTION_SELECT {
		return errors.New("this function only for select")
	}
	qb.Limit(1)
	qb.OrderBy("id", ORDER_BY_DESC)

	rows, err := qb.executeQuerySelect()
	if err != nil {
		return err
	}
	err = qb.ScanRow(rows, entity)

	if err != nil {
		return err
	}
	return nil
}

// return many value filter by condition
func (qb *QueryBuilder) Find(entities interface{}) error {
	defer qb.deferFunc()
	if qb.Action != ACTION_SELECT {
		return errors.New("this function only for select")
	}
	rows, err := qb.executeQuerySelect()
	if err != nil {
		return err
	}
	err = qb.ScanRows(rows, entities)

	if err != nil {
		return err
	}
	return nil
}

// this function for executing query condition
func (qb *QueryBuilder) execute(query string) error {
	defer qb.deferFunc()
	query = qb.ValidatedQueryAndMapping(query)             // merubah symbol ? -> $n
	util.DebugQueryAndParams(qb.DebugMode, query, qb.Args) // for debuging

	if qb.DBTransaction != nil {
		_, err := qb.DBTransaction.Exec(query, qb.Args...)
		if err != nil {
			return err
		}
		return nil
	}
	_, err := qb.DBConnection.Exec(query, qb.Args...)
	if err != nil {
		return err
	}
	return nil
}

// this function for save data and returning id
func (qb *QueryBuilder) save(query string) (interface{}, error) {
	defer qb.deferFunc()
	var primeryKey interface{}
	query = qb.ValidatedQueryAndMapping(query) // merubah symbol ? -> $n
	queryCallback := fmt.Sprintf("%s RETURNING id", query)
	util.DebugQueryAndParams(qb.DebugMode, query, qb.Args) // for debuging
	if qb.DBTransaction != nil {
		err := qb.DBTransaction.QueryRow(queryCallback, qb.Args...).Scan(&primeryKey)
		if err != nil {
			return primeryKey, err
		}
		return primeryKey, nil
	}
	err := qb.DBConnection.QueryRow(queryCallback, qb.Args...).Scan(&primeryKey)
	if err != nil {
		return primeryKey, err
	}
	return primeryKey, nil
}

// this function for query raw and can custom query
func (qb *QueryBuilder) executeRawQuery(query string) (*sql.Rows, error) {
	defer qb.deferFunc()
	query = qb.ValidatedQueryAndMapping(query)             // merubah symbol ? -> $n
	util.DebugQueryAndParams(qb.DebugMode, query, qb.Args) // for debuging
	if qb.DBTransaction != nil {
		res, err := qb.DBTransaction.Query(query, qb.Args...)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	res, err := qb.DBConnection.Query(query, qb.Args...)
	if err != nil {
		return res, err
	}
	return res, nil
}
