package sqldb

import (
	"fmt"
	"strconv"

	"github.com/gilperopiola/obreros/core"
)

type Operation string

const (
	Where Operation = "where"
	And   Operation = "and"
	Or    Operation = "or"
	Like  Operation = "like"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - High Level SQL Options -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func WithID(id int32) core.SqlDBOpt {
	return WithCondition(Where, "id", strconv.Itoa(int(id)))
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Low Level SQL Options -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func WithCondition(operation Operation, field, value string) core.SqlDBOpt {
	if field == "" {
		return func(db core.SqlDB) {} // No-op
	}

	return func(db core.SqlDB) {
		if operation == Where || operation == And { // Where / And
			db.Where(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Or { // Or
			db.Or(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Like { // Like
			db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
			return
		}
	}
}
