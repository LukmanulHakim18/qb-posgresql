package querybuilder

import (
	"fmt"
	"strings"
)

type Column struct {
	Name  string
	Value interface{}
}

func (qb *queryBuilder) selectToString() string {
	if len(qb.Columns) == 0 {
		return fmt.Sprintf("SELECT * FROM %s", qb.TableName)
	}
	columns := []string{}
	for _, v := range qb.Columns {
		columns = append(columns, v.Name)
	}
	columnsStr := strings.Join(columns, ", ")
	querySelect := fmt.Sprintf("SELECT %sFROM %s", columnsStr, qb.TableName)
	return querySelect
}

func (qb *queryBuilder) updateToString() string {
	querySelect := fmt.Sprintf("UPDATE %s ", qb.TableName)
	return querySelect
}
func (qb *queryBuilder) insertToString() string {
	querySelect := fmt.Sprintf("INSERT INTO %s ", qb.TableName)
	return querySelect
}
func (qb *queryBuilder) deleteToString() string {
	querySelect := fmt.Sprintf("DELETE FROM %s ", qb.TableName)
	return querySelect
}

func (qb *queryBuilder) mappingDataSetUpdste() string {
	fieldTemp := []string{}
	for _, column := range qb.Columns {
		fieldTemp = append(fieldTemp, fmt.Sprintf("%s = ?", column.Name))
		qb.Args = append(qb.Args, column.Value)
	}

	fieldUpdate := fmt.Sprintf("SET %s ", strings.Join(fieldTemp, ", "))
	return fieldUpdate
}
func (qb *queryBuilder) mappingDataSetInsert() string {
	fieldTemp := []string{}
	for _, column := range qb.Columns {
		fieldTemp = append(fieldTemp, fmt.Sprintf("%s", column.Name))
		qb.Args = append(qb.Args, column.Value)
	}

	columnInsert := fmt.Sprintf("(%s)", strings.Join(fieldTemp, ", "))
	valuesInsert := strings.Replace(strings.Repeat(", ? ", len(fieldTemp)), ", ", "", 1)
	queryInsert := fmt.Sprintf("%s VALUES (%s)", columnInsert, valuesInsert)
	return queryInsert
}
