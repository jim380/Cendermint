package services

import "database/sql"

type RpcServiceExecutor interface {
	Init(*sql.DB)
}

type RestServiceExecutor interface {
	Init(*sql.DB)
}
