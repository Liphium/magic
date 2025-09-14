package starter

import (
	"github.com/Liphium/magic"
	"github.com/Liphium/magic/mconfig"
)

func BuildMagicConfig() magic.Config {
	return magic.Config{
		AppName: "magic-example-real-project",
		PlanDeployment: func(ctx *mconfig.Context) {

		},
		StartFunction: func() {},
	}
}
