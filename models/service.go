package models

import "database/sql"

type RpcServices interface {
	Init(*sql.DB)
}

type RestServices interface {
	Init(*sql.DB)
}
