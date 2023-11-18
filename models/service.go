package models

import "database/sql"

type DBService interface {
	Init(*sql.DB)
}
