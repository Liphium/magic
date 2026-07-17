package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/moby/moby/client"
)

// For all clear instructions, flush everything away. Similar to other databases, but since we don't have schemas we flush everything.
func (rd *RedisDriver) HandleInstruction(ctx context.Context, c *client.Client, container mconfig.ContainerInformation, instruction mconfig.Instruction) error {
	switch instruction {
	case mconfig.InstructionDropTables, mconfig.InstructionClearTables:
		return rd.FlushAll(ctx, c, container)
	}
	return nil
}

// FlushAll executes FLUSHALL with async to clear all Redis data.
func (rd *RedisDriver) FlushAll(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) error {
	cmd := strings.Split("redis-cli -a "+RedisPassword+" flushall async", " ")

	resp, err := mservices.ExecuteCommand(ctx, c, container.ID, cmd)
	if err != nil {
		return fmt.Errorf("couldn't flush redis: %s", err)
	}
	if resp.ExitCode != 0 {
		return fmt.Errorf("redis flush failed with exit code %d", resp.ExitCode)
	}

	return nil
}
