package qb_implement

import (
	"fmt"
	"strings"
)

const (
	CONDITION_WHERE_IN      = "IN"
	CONDITION_WHERE_OR      = "OR"
	CONDITION_WHERE_AND     = "AND"
	CONDITION_WHERE_NOT_IN  = "NOT IN"
	CONDITION_WHERE_BETWEEN = "BETWEEN"
)

type Condition struct {
	ColumnName string      // accommodate column name
	Symbol     string      // accommodate symbol for compairing condition
	Value      interface{} // accommodate value
	Connector  string      // accommodate conector like "AND, OR"
}

type OrderBy struct {
	ColumnName string // accommodate column name
	Value      string // accommodate value
}

// set condition to string and make argument for query params
func (qb *QueryBuilder) conditionToString() string {
	if len(qb.Conditions) == 0 {
		return ""
	}
	conditions := []string{}

	// looping condition field
	for _, v := range qb.Conditions {
		switch v.Connector {
		case CONDITION_WHERE_IN, CONDITION_WHERE_NOT_IN: // jika kondisi where in
			params := []string{}
			args := v.Value.([]interface{})
			for _, v := range args {
				qb.Args = append(qb.Args, v)
				params = append(params, "?")

			}
			paramsString := strings.Join(params, ", ")
			queryWhereIn := fmt.Sprintf("AND %s %s (%s)", v.ColumnName, v.Connector, paramsString)
			conditions = append(conditions, queryWhereIn)
		case CONDITION_WHERE_BETWEEN:
			args := v.Value.([]interface{})
			queryBetwen := fmt.Sprintf("AND (%s %s (?) AND (?))", v.ColumnName, v.Connector)
			conditions = append(conditions, queryBetwen)
			qb.Args = append(qb.Args, args...)
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
func (qb *QueryBuilder) limitToString() string {
	if qb.LimitVal <= 0 {
		return ""
	}
	queryLimit := fmt.Sprintf("LIMIT %d", qb.LimitVal)
	return queryLimit

}
func (qb *QueryBuilder) offsetToString() string {
	if qb.OffsetValue <= 0 {
		return ""
	}
	queryLimit := fmt.Sprintf("OFFSET %d", qb.OffsetValue)
	return queryLimit

}

// set order by to query string
func (qb *QueryBuilder) orderByToString() string {
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
func (qb *QueryBuilder) groupByToString() string {
	if len(qb.GroupByConditions) == 0 {
		return ""
	}

	querGroupBy := fmt.Sprintf("GROUP BY %s", strings.Join(qb.GroupByConditions, ", "))
	return querGroupBy

}
