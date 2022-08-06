package querybuilder

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// ======================================= Mapping Section =======================================

func (qb *queryBuilder) ValidatedQueryAndMapping(query string) string {
	for k, _ := range qb.Args {
		key := "$" + strconv.Itoa(k+1)
		query = strings.Replace(query, "?", key, 1)
	}
	return query
}
func (qb *queryBuilder) ScanEntity(entity interface{}) {
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	for i := 0; i < t.NumField(); i++ {
		tagDB := t.Field(i).Tag.Get("db")
		name := t.Field(i).Name
		fv := v.FieldByName(name)
		if tagDB != "" && tagDB != "id" {
			colUpdate := Column{Name: tagDB, Value: fv.Interface()}
			qb.Columns = append(qb.Columns, colUpdate)
		}
		if tagDB == "id" {
			qb.PrimeryKey = fv.Interface()
		}
	}
}

func (qb *queryBuilder) ScanRow(rows *sql.Rows, entity interface{}) error {
	defer rows.Close()
	if !qb.isPointer(entity) {
		return errors.New("parameter must pointer of entity")
	}
	if !rows.Next() {
		return nil
	}
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i, _ := range values {
		pointers[i] = &values[i]
	}
	err := rows.Scan(pointers...)
	if err != nil {
		return err
	}

	resultMap := make(map[string]interface{})
	for i, val := range values {
		typeOf := reflect.ValueOf(val).Kind()
		if typeOf == reflect.Slice {
			value := string(val.([]byte))
			strToFloat, _ := strconv.ParseFloat(value, 64)
			val = strToFloat
		}
		resultMap[columns[i]] = val
	}
	byt, err := json.Marshal(resultMap)
	if err != nil {
		return err
	}
	json.Unmarshal(byt, entity)

	return nil
}
func (qb *queryBuilder) ScanRows(rows *sql.Rows, entities interface{}) error {
	defer rows.Close()
	if !qb.isPointer(entities) {
		return errors.New("parameter must pointer of entities")
	}
	columns, _ := rows.Columns()
	resultSlice := []interface{}{}
	for rows.Next() {

		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i, _ := range values {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			return err
		}

		resultMap := make(map[string]interface{})
		for i, val := range values {
			typeOf := reflect.ValueOf(val).Kind()
			if typeOf == reflect.Slice {
				value := string(val.([]byte))
				strToFloat, _ := strconv.ParseFloat(value, 64)
				val = strToFloat
			}
			resultMap[columns[i]] = val
		}

		resultSlice = append(resultSlice, resultMap)
	}
	byt, err := json.Marshal(resultSlice)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byt, entities)

	return nil
}

func (qb *queryBuilder) clearData() {
	qb.TableName = ""
	qb.Action = ""
	qb.Conditions = nil
	qb.OrderByConditions = nil
	qb.Columns = nil
	qb.LimitVal = 0
	qb.OffsetValue = 0
	qb.PrimeryKey = nil
	qb.Args = nil
}

func (qb *queryBuilder) isStruct(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Struct
}
func (qb *queryBuilder) isPointer(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Pointer
}

func (qb *queryBuilder) deferFunc() {
	qb.clearData()
}
