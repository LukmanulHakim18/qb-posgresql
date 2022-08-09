package qb_posgresql

import (
	"database/sql"

	utility "github.com/LukmanulHakim18/qb-posgresql/utility"
	qbImplement "github.com/LukmanulHakim18/qb-posgresql/qb_implement"
	qbInterface "github.com/LukmanulHakim18/qb-posgresql/qb_interface"
)

// for create query builder and representation interface
func NewQueryBuilder(db *sql.DB) qbInterface.QueryBuilder {
	return &qbImplement.QueryBuilder{
		DBCon:     db,
		DebugMode: utility.GetBoolEnv("DEBUG_MODE"),
	}
}
