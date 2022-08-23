package qb_implement

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// scaning entity to get tag db and value
func (qb *QueryBuilder) ScanEntity(entity interface{}) *QueryBuilder {
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	for i := 0; i < t.NumField(); i++ {
		tagDB := t.Field(i).Tag.Get("db")
		tagDB = strings.TrimSpace(tagDB)
		tags := strings.Split(strings.TrimSpace(tagDB), ", ")
		if len(tags) > 1 && tags[1] == "primaryKey" {
			name := t.Field(i).Name
			fv := v.FieldByName(name)
			PkColumn := Column{Name: tags[0], Value: fv.Interface()}
			qb.PrimeryKey = &PkColumn
		} else {
			name := t.Field(i).Name
			fv := v.FieldByName(name)
			if tagDB != "" && tagDB != "id" {
				colUpdate := Column{Name: tagDB, Value: fv.Interface()}
				qb.Columns = append(qb.Columns, colUpdate)
			}
			if tagDB == "id" {
				PkColumn := Column{Name: "id", Value: fv.Interface()}
				qb.PrimeryKey = &PkColumn
			}
		}

	}
	return qb
}

// use this for scan one row data
func (qb *QueryBuilder) ScanRow(rows *sql.Rows, entity interface{}) error {
	defer rows.Close()
	if !rows.Next() {
		return errors.New("data not found")
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
	// spew.Dump(resultMap)
	return qb.populateData(entity, resultMap)
}

// use this for scan many rows datas
func (qb *QueryBuilder) ScanRows(rows *sql.Rows, entities interface{}) error {
	defer rows.Close()
	columns, _ := rows.Columns()
	resultSlice := []map[string]interface{}{}
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
	return qb.populateDatas(entities, resultSlice)
}

// for support populating datas pointer of struct when ScanRow
func (qb *QueryBuilder) populateData(ptrEntity interface{}, data map[string]interface{}) error {

	// mencari tau apakah param adalah pointer
	if ptrEntity == nil {
		return errors.New("must be pointer of struct")
	}

	val := reflect.ValueOf(ptrEntity)
	rt := reflect.TypeOf(ptrEntity)

	// If it's an interface or a pointer, unwrap it.
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
		rt = rt.Elem()
	} else {
		return errors.New("must be pointer of struct")
	}
	valNumFields := rt.NumField()
	for i := 0; i < valNumFields; i++ {
		field := rt.Field(i)
		tagDB := field.Tag.Get("db")
		tagDB = strings.TrimSpace(tagDB)
		tags := strings.Split(strings.TrimSpace(tagDB), ", ")
		tag := tags[0]
		if tag == "" || tag == "-" {
			continue
		}
		name := field.Name
		refValByname := val.FieldByName(name)
		refdataVal := reflect.ValueOf(data[tag])
		switch field.Type.String() {
		case "*time.Time":
			if data[tag] != nil {
				t, _ := data[tag].(time.Time)
				refdataVal := reflect.ValueOf(&t)
				refValByname.Set(refdataVal.Convert(field.Type))
			}
		case "sql.NullString":
			newData := sql.NullString{}
			newData.Scan(data[tag])
			refdataVal := reflect.ValueOf(newData)
			refValByname.Set(refdataVal.Convert(field.Type))
		default:
			if data[tag] != nil {
				refValByname.Set(refdataVal.Convert(field.Type))
			}
		}
	}
	return nil
}

// for support populating datas pointer of struct slice when ScanRows
func (qb *QueryBuilder) populateDatas(ptrEntities interface{}, datas []map[string]interface{}) error {

	// mencari tau apakah param adalah pointer
	if ptrEntities == nil {
		return errors.New("must be pointer of slice struct")
	}
	sliceVal := reflect.ValueOf(ptrEntities).Elem()

	if sliceVal.Kind() != reflect.Slice {
		return errors.New("must be pointer of slice struct")
	}

	elemType := sliceVal.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return errors.New("must be pointer of slice struct")
	}
	for _, data := range datas {
		newVal := reflect.New(elemType).Elem()

		valNumFields := elemType.NumField()
		for i := 0; i < valNumFields; i++ {
			field := elemType.Field(i)
			tagDB := field.Tag.Get("db")
			tagDB = strings.TrimSpace(tagDB)
			tags := strings.Split(strings.TrimSpace(tagDB), ", ")
			tag := tags[0]
			if tag == "" || tag == "-" {
				continue
			}
			name := field.Name
			refValByname := newVal.FieldByName(name)
			refDataVal := reflect.ValueOf(data[tag])
			switch field.Type.String() {
			case "*time.Time":
				if data[tag] != nil {
					t, _ := data[tag].(time.Time)
					refdataVal := reflect.ValueOf(&t)
					refValByname.Set(refdataVal.Convert(field.Type))
				}
			case "sql.NullString":
				newData := sql.NullString{}
				newData.Scan(data[tag])
				refdataVal := reflect.ValueOf(newData)
				refValByname.Set(refdataVal.Convert(field.Type))
			default:
				if data[tag] != nil {
					refValByname.Set(refDataVal.Convert(field.Type))
				}
			}

		}
		sliceVal.Set(reflect.Append(sliceVal, newVal))
	}

	return nil
}
