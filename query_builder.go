package qb_posgresql

import (
	"database/sql"

	qbImplement "github.com/LukmanulHakim18/qb-posgresql/qb_implement"
	qbInterface "github.com/LukmanulHakim18/qb-posgresql/qb_interface"
	utility "github.com/LukmanulHakim18/qb-posgresql/utility"
)

// for create query builder and representation interface
func NewQueryBuilder(db *sql.DB) qbInterface.QueryBuilder {
	return &qbImplement.QueryBuilder{
		DBConnection: db,
		DebugMode:    utility.GetBoolEnv("DEBUG_MODE"),
	}
}
