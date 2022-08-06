package querybuilder

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	CONDITION_WHERE_IN     = "IN"
	CONDITION_WHERE_OR     = "OR"
	CONDITION_WHERE_AND    = "AND"
	CONDITION_WHERE_NOT_IN = "NOT IN"
)

type condition struct {
	ColumnName string
	Symbol     string
	Value      interface{}
	Connector  string
}

type orderBy struct {
	ColumnName string
	Value      string
}

// set condition to string and make argument for query params
func (qb *queryBuilder) conditionToString() string {
	if len(qb.Conditions) == 0 {
		return ""
	}
	conditions := []string{}

	// looping condition field
	for _, v := range qb.Conditions {
		switch v.Connector {
		case CONDITION_WHERE_IN, CONDITION_WHERE_NOT_IN: // jika kondisi where in
			params := []string{}
			switch reflect.TypeOf(v.Value).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(v.Value)
				for i := 0; i < s.Len(); i++ {
					qb.Args = append(qb.Args, s.Index(i).Interface())
					params = append(params, " ?")
				}
			}
			paramsString := strings.Join(params, ", ")
			queryWhereIn := fmt.Sprintf("AND %s %s (%s)", v.ColumnName, v.Connector, paramsString)
			conditions = append(conditions, queryWhereIn)
		default:
			conditions = append(conditions, fmt.Sprintf("%s %s %s ?", v.Connector, v.ColumnName, v.Symbol))
			qb.Args = append(qb.Args, v.Value)

		}
	}

	queryCondition := strings.Join(conditions, " ")
	replaceAtBeginning := strings.Replace(queryCondition, "AND", "WHERE", 1)
	return replaceAtBeginning
}

// set limit to query string
func (qb *queryBuilder) limitToString() string {
	if qb.LimitVal <= 0 {
		return ""
	}
	queryLimit := fmt.Sprintf("LIMIT %d", qb.LimitVal)
	return queryLimit

}

// set order by to query string
func (qb *queryBuilder) orderByToString() string {
	if len(qb.OrderByConditions) == 0 {
		return ""
	}
	orderBy := []string{}
	for _, v := range qb.OrderByConditions {
		orderBy = append(orderBy, fmt.Sprintf("%s %s", v.ColumnName, v.Value))
	}

	querOrderby := fmt.Sprintf("ORDER BY %s", strings.Join(orderBy, ", "))
	return querOrderby

}
