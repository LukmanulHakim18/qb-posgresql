package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
)

// ======================================= Execute Section =======================================

//execute query
func (qb *queryBuilder) executeQuerySelect() (*sql.Rows, error) {
	query := fmt.Sprintf("%s ", qb.selectToString())            // set table and field
	query = fmt.Sprintf("%s %s", query, qb.conditionToString()) // add condition to string
	query = fmt.Sprintf("%s %s", query, qb.orderByToString())   // add orderBy to string
	query = fmt.Sprintf("%s %s", query, qb.limitToString())     // add limit to string
	query = qb.ValidatedQueryAndMapping(query)                  // merubah symbol ? -> $n
	if qb.DebugMode {
		fmt.Println(query)
	}
	rows, err := qb.DBCon.Query(query, qb.Args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// return 1 value filter by condition
func (qb *queryBuilder) FindOne(entity interface{}) error {
	defer qb.deferFunc()
	if qb.Action != ACTION_SELECT {
		return errors.New("this function only for select")
	}
	if !qb.isPointer(entity) {
		return errors.New("parameter must pointer of entities")
	}
	qb.Limit(1)
	qb.OrderBy("id", ORDER_BY_ASC)

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
func (qb *queryBuilder) Find(entities interface{}) error {
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

func (qb *queryBuilder) execute(query string) error {
	defer qb.deferFunc()
	query = qb.ValidatedQueryAndMapping(query) // merubah symbol ? -> $n
	if qb.DebugMode {
		fmt.Println(query)
	}
	_, err := qb.DBCon.Exec(query, qb.Args...)
	if err != nil {
		return err
	}
	return nil
}

func (qb *queryBuilder) save(query string) (interface{}, error) {
	defer qb.deferFunc()
	var primeryKey interface{}
	query = qb.ValidatedQueryAndMapping(query) // merubah symbol ? -> $n
	queryCallback := fmt.Sprintf("%s RETURNING id", query)
	if qb.DebugMode {
		fmt.Println(query)
	}
	err := qb.DBCon.QueryRow(queryCallback, qb.Args...).Scan(&primeryKey)
	if err != nil {
		return primeryKey, err
	}
	return primeryKey, nil
}

func (qb *queryBuilder) executeRaw(query string) (*sql.Rows, error) {
	defer qb.deferFunc()
	query = qb.ValidatedQueryAndMapping(query) // merubah symbol ? -> $n
	if qb.DebugMode {
		fmt.Println(query)
	}
	res, err := qb.DBCon.Query(query, qb.Args...)
	if err != nil {
		return res, err
	}
	return res, nil
}
