package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/moby/moby/client"
)

func (rd *RedisDriver) CreateContainer(ctx context.Context, c *client.Client, a mconfig.ContainerAllocation) (string, error) {
	if rd.Image == "" {
		return "", fmt.Errorf("please specify a proper image")
	}

	return mservices.CreateContainer(ctx, redisLog, c, a, mservices.ManagedContainerOptions{
		Image: rd.Image,
		Ports: []string{
			"6379/tcp",
		},
		Volumes: []mservices.ContainerVolume{
			{NameSuffix: "data", Target: "/data"},
		},
	})
}

func (rd *RedisDriver) IsHealthy(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) (bool, error) {
	readyCommand := "redis-cli ping"
	cmd := strings.Split(readyCommand, " ")

	respInspect, err := mservices.ExecuteCommand(ctx, c, container.ID, cmd)
	if err != nil {
		return false, fmt.Errorf("couldn't execute command for readiness of container: %s", err)
	}

	return respInspect.ExitCode == 0, nil
}

func (rd *RedisDriver) Initialize(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) error {
	return nil
}
