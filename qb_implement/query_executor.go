package qb_implement

import (
	"database/sql"
	"errors"
	"fmt"

	util "github.com/LukmanulHakim18/qb-posgresql/utility"
)

// ======================================= Execute Section =======================================

//execute query

// return 1 value filter by condition
func (qb *QueryBuilder) FindOne(entity interface{}) error {
	if qb.Action != ACTION_SELECT {
		return errors.New("this function only for select")
	}
	qb.Limit(1)
	if qb.PrimeryKey != nil {
		qb.OrderBy(qb.PrimeryKey.Name, ORDER_BY_DESC)
	}

	rows, err := qb.Get()

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

	rows, err := qb.Get()
	if err != nil {
		return err
	}
	err = qb.ScanRows(rows, entities)

	if err != nil {
		return err
	}
	return nil
}

// this function for save data and returning id
func (qb *QueryBuilder) save() (interface{}, error) {
	defer qb.deferFunc()
	var primeryKey interface{}
	qb.Build()
	query := qb.ValidatedQueryAndMapping(qb.Query) // merubah symbol ? -> $n
	if primeryKey != nil {
		query = fmt.Sprintf("%s RETURNING %s", query, qb.PrimeryKey.Name)
		util.DebugQueryAndParams(qb.DebugMode, query, qb.Args) // for debuging
		if qb.DBTransaction != nil {
			err := qb.DBTransaction.QueryRow(query, qb.Args...).Scan(&primeryKey)
			if err != nil {
				return primeryKey, err
			}
			return primeryKey, nil
		}
		err := qb.DBConnection.QueryRow(query, qb.Args...).Scan(&primeryKey)
		if err != nil {
			return nil, err
		}
		return primeryKey, err
	} else {
		if qb.DBTransaction != nil {
			res, err := qb.DBTransaction.Exec(query, qb.Args...)
			if err != nil {
				return primeryKey, err
			}

			return res.RowsAffected()
		}
		res, err := qb.DBConnection.Exec(query, qb.Args...)
		if err != nil {
			return nil, err
		}

		return res.RowsAffected()
	}

}

// returning rows affected and error
func (qb *QueryBuilder) Exec() (int64, error) {
	defer qb.deferFunc()
	qb.Build()                                             // Build Query
	query := qb.ValidatedQueryAndMapping(qb.Query)         // merubah symbol ? -> $n
	util.DebugQueryAndParams(qb.DebugMode, query, qb.Args) // for debuging
	if qb.DBTransaction != nil {
		res, err := qb.DBTransaction.Exec(query, qb.Args...)
		if err != nil {
			return 0, err
		}
		ra, _ := res.RowsAffected()
		return ra, nil
	}
	res, err := qb.DBConnection.Exec(query, qb.Args...)
	if err != nil {
		return 0, err
	}
	ra, _ := res.RowsAffected()
	return ra, nil
}

// returning *sql.Rows, and error
func (qb *QueryBuilder) Get() (*sql.Rows, error) {
	defer qb.deferFunc()
	qb.Build()                                             // build query
	query := qb.ValidatedQueryAndMapping(qb.Query)         // merubah symbol ? -> $n
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
