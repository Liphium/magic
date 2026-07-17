package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/moby/moby/client"
)

// CreateContainer creates Redis container with password auth.
func (rd *RedisDriver) CreateContainer(ctx context.Context, c *client.Client, a mconfig.ContainerAllocation) (string, error) {
	if rd.Image == "" {
		return "", fmt.Errorf("please specify a proper image")
	}

	return mservices.CreateContainer(ctx, redisLog, c, a, mservices.ManagedContainerOptions{
		Image: rd.Image,
		Env: []string{
			fmt.Sprintf("REDIS_PASSWORD=%s", RedisPassword),
		},
		Ports: []string{
			"6379/tcp",
		},
		Volumes: []mservices.ContainerVolume{
			{NameSuffix: "data", Target: "/data"},
		},
	})
}

// IsHealthy checks Redis readiness via redis-cli ping.
func (rd *RedisDriver) IsHealthy(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) (bool, error) {
	readyCmd := "redis-cli -a " + RedisPassword + " ping"
	cmd := strings.Split(readyCmd, " ")

	resp, err := mservices.ExecuteCommand(ctx, c, container.ID, cmd)
	if err != nil {
		return false, fmt.Errorf("couldn't execute command for readiness of container: %s", err)
	}

	if mconfig.VerboseLogging {
		redisLog.Println("Redis health check response code:", resp.ExitCode)
	}

	return resp.ExitCode == 0, nil
}

// Initialize does nothing. Redis needs no schema setup / create databases.
func (rd *RedisDriver) Initialize(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) error {
	return nil
}
