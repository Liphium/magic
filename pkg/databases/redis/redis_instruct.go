package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/moby/moby/client"
)

func (rd *RedisDriver) HandleInstruction(ctx context.Context, c *client.Client, container mconfig.ContainerInformation, instruction mconfig.Instruction) error {
	switch instruction {
	case mconfig.InstructionClearTables:
		return rd.FlushDB(ctx, c, container)
	case mconfig.InstructionDropTables:
		return rd.FlushDB(ctx, c, container)
	}
	return nil
}

func (rd *RedisDriver) FlushDB(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) error {
	redisLog.Println("Flushing Redis database...")
	cmd := strings.Split("redis-cli flushall", " ")

	respInspect, err := mservices.ExecuteCommand(ctx, c, container.ID, cmd)
	if err != nil {
		return fmt.Errorf("couldn't execute flushall: %s", err)
	}

	if respInspect.ExitCode != 0 {
		return fmt.Errorf("flushall failed with exit code %d", respInspect.ExitCode)
	}

	return nil
}
