package qb_implement

import (
	"strconv"
	"strings"
)

// for mapping query, change ? to $n based on the number of arguments
func (qb *QueryBuilder) ValidatedQueryAndMapping(query string) string {
	for k, _ := range qb.Args {
		key := "$" + strconv.Itoa(k+1)
		query = strings.Replace(query, "?", key, 1)
	}
	return query
}

// cleare data for next query when use singel repo
func (qb *QueryBuilder) clearData() {
	qb.TableName = ""
	qb.Action = ""
	qb.Conditions = nil
	qb.OrderByConditions = nil
	qb.Columns = nil
	qb.LimitVal = 0
	qb.OffsetValue = 0
	qb.PrimeryKey = nil
	qb.Args = nil
	qb.Query = ""
}

// function defer
func (qb *QueryBuilder) deferFunc() {
	qb.clearData()
}
