package mservices

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

// Simply execute a command inside of a container.
func ExecuteCommand(ctx context.Context, c *client.Client, id string, cmd []string) (client.ExecInspectResult, error) {
	execIDResp, err := c.ExecCreate(ctx, id, client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return client.ExecInspectResult{}, fmt.Errorf("couldn't create command for readiness of container: %s", err)
	}
	execStartCheck := client.ExecStartOptions{Detach: false, TTY: false}
	if _, err := c.ExecStart(ctx, execIDResp.ID, execStartCheck); err != nil {
		return client.ExecInspectResult{}, fmt.Errorf("couldn't start command for readiness of container: %s", err)
	}
	return c.ExecInspect(ctx, execIDResp.ID, client.ExecInspectOptions{})
}
