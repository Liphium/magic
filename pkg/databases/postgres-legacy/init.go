package postgres_legacy

import "github.com/Liphium/magic/v2/mconfig"

func init() {
	mconfig.RegisterDriver(&PostgresDriver{})
}
