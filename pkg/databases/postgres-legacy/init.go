package postgres_legacy

import "github.com/Liphium/magic/v3/mconfig"

func init() {
	mconfig.RegisterDriver(&PostgresDriver{})
}
