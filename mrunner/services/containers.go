package mservices

import (
	"context"
	"fmt"
	"log"
	"net/netip"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// ContainerVolume describes a single named volume that should be mounted into
// the container. The volume name is derived from the container name so it
// survives container re-creation.
type ContainerVolume struct {
	// Suffix appended to the container name to form the Docker volume name,
	// e.g. "data" -> "<containerName>-data".
	NameSuffix string
	// Absolute path inside the container where the volume is mounted.
	Target string
}

// ManagedContainerOptions holds everything needed to create (or re-create) a
// managed Docker container in a reproducible way.
type ManagedContainerOptions struct {
	// Docker image to use, e.g. "postgres:17".
	Image string
	// Environment variables passed into the container.
	Env []string
	// Ports to expose. Each entry maps one container port (inside of the container) to one host port (chosen by Magic).
	Ports []string
	// Named volumes to attach. Existing mounts are reused across re-creations.
	Volumes []ContainerVolume
	// Commands to run in the container (overrides the image's default CMD).
	Cmd []string
}

// CreateContainer looks for an existing container with the given name and reuses it when its configuration still matches the requested options.
//
// When the configuration differs (image, env, command, ports, volumes), it removes the old container and creates a fresh one from the provided options.
//
// Existing Docker volumes are always preserved so that data survives a container re-creation. Returns the ID of the currently running/reusable or freshly created container.
func CreateContainer(ctx context.Context, log *log.Logger, c *client.Client, a mconfig.ContainerAllocation, opts ManagedContainerOptions) (string, error) {
	if opts.Image == "" {
		return "", fmt.Errorf("please specify a proper image")
	}

	// Look for an existing container first. If it exists and its configuration still matches, we can reuse it without recreating.
	existing, err := FindContainer(ctx, c, a.Name)
	if err != nil {
		return "", err
	}

	if existing != nil {
		existingID := existing.ID

		difference := configMatches(existing, a, opts)
		if difference == "" {
			log.Printf("Reusing existing container %q, configuration unchanged", a.Name)
			return existingID, nil
		}

		log.Printf("Existing container %q configuration changed, recreating...", a.Name)
		if mconfig.VerboseLogging {
			log.Println("Difference:", difference)
		}

		// Version guard before destroying anything
		majorCurrent := GetImageMajorVersion(existing.Config.Image)
		majorNew := GetImageMajorVersion(opts.Image)
		if majorCurrent == -1 || majorNew == -1 {
			log.Printf("Skipping major version check: unable to parse image versions (%s -> %s)", existing.Config.Image, opts.Image)
		} else if majorCurrent != majorNew {
			return "", fmt.Errorf("major version mismatch for %s: please upgrade or delete your old container before starting the app (%d -> %d)", opts.Image, majorCurrent, majorNew)
		}

		log.Println("Removing old container...")
		if _, err := c.ContainerRemove(ctx, existingID, client.ContainerRemoveOptions{
			RemoveVolumes: false,
			Force:         true,
		}); err != nil {
			return "", fmt.Errorf("couldn't remove existing container: %s", err)
		}
	}

	mounts, err := buildMounts(a, opts.Volumes)
	if err != nil {
		return "", err
	}

	exposedPorts, portBindings, err := buildPortBindings(a, opts.Ports)
	if err != nil {
		return "", err
	}

	resp, err := c.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image:        opts.Image,
			Env:          opts.Env,
			Cmd:          opts.Cmd,
			ExposedPorts: exposedPorts,
		},
		HostConfig: &container.HostConfig{
			PortBindings: portBindings,
			Mounts:       mounts,
		},
		Name: a.Name,
	})
	if err != nil {
		return "", fmt.Errorf("couldn't create container %q: %s", a.Name, err)
	}

	return resp.ID, nil
}

// FindContainer looks for a container with the given name and returns its inspect data. Returns nil when no container matches.
func FindContainer(ctx context.Context, c *client.Client, name string) (*container.InspectResponse, error) {
	f := make(client.Filters)
	f.Add("name", name)
	summary, err := c.ContainerList(ctx, client.ContainerListOptions{
		Filters: f,
		All:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't list containers: %s", err)
	}

	for _, ct := range summary.Items {
		for _, n := range ct.Names {
			if !strings.HasSuffix(n, name) {
				continue
			}

			resp, err := c.ContainerInspect(ctx, ct.ID, client.ContainerInspectOptions{})
			if err != nil {
				return nil, fmt.Errorf("couldn't inspect container %q: %s", name, err)
			}
			return &resp.Container, nil
		}
	}

	return nil, nil
}

// configMatches reports whether a found container can be reused as-is for the given options. It compares the fields Magic controls: image, env, command, exposed container ports, host port bindings and the managed volume mounts.
//
// Returns "" if there is no conflict.
func configMatches(existing *container.InspectResponse, a mconfig.ContainerAllocation, opts ManagedContainerOptions) string {
	if existing.Config == nil || existing.HostConfig == nil {
		return "existing container config is nil"
	}

	// Image
	if existing.Config.Image != opts.Image {
		return "image is different"
	}

	// Environment (order independent)
	if !envEqual(existing.Config.Env, opts.Env) {
		return "environment variables are different"
	}

	// Command
	if !reflect.DeepEqual(existing.Config.Cmd, opts.Cmd) {
		return "start command is different"
	}

	// Ports: both the exposed container port set and the host bindings. A re-rolled host port relative to a reused container means the allocation no longer matches the actual binding, so recreate.
	_, expectedBindings, err := buildPortBindings(a, opts.Ports)
	if err != nil {
		return "couldn't build port bindings"
	}
	if len(existing.HostConfig.PortBindings) != len(expectedBindings) {
		return "expected host port bindings have different length compared to current ones"
	}
	for port, bindings := range expectedBindings {
		existingBindings, ok := existing.HostConfig.PortBindings[port]
		if !ok {
			return "couldn't find host binding of existing port"
		}
		if len(existingBindings) != len(bindings) {
			return "length of existing bindings isn't equal to expected one"
		}
		for i, b := range bindings {
			if existingBindings[i] != b {
				return "existing binding is different from existing one"
			}
		}
	}

	// Managed volume mounts
	if !volumesMatch(existing, a, opts.Volumes) {
		return "volumes don't match"
	}

	return ""
}

// envEqual compares two environment variable slices regardless of order.
func envEqual(current, expected []string) bool {

	// We can ignore other environment variables (Docker for example adds PATH, which we don't care about)
	for _, ex := range expected {
		if !slices.Contains(current, ex) {
			return false
		}
	}
	return true
}

// volumesMatch verifies that every managed volume is mounted at the requested target with the expected named-volume source.
func volumesMatch(existing *container.InspectResponse, a mconfig.ContainerAllocation, volumes []ContainerVolume) bool {
	// Index the existing mounts by their target for lookup
	byTarget := map[string]mount.Mount{}
	for _, m := range existing.HostConfig.Mounts {
		byTarget[m.Target] = m
	}

	for _, vol := range volumes {
		m, ok := byTarget[vol.Target]
		if !ok {
			return false
		}
		// Only compare the pieces we control; bind mounts created by Docker
		// internals (e.g. resolv.conf) are irrelevant here.
		if m.Type != mount.TypeVolume {
			return false
		}
		if m.Source != fmt.Sprintf("%s-%s", a.Name, vol.NameSuffix) {
			return false
		}
	}

	return true
}

// buildMounts constructs the mount list for the new container.
func buildMounts(a mconfig.ContainerAllocation, volumes []ContainerVolume) ([]mount.Mount, error) {
	mounts := make([]mount.Mount, 0, len(volumes))

	for _, vol := range volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: fmt.Sprintf("%s-%s", a.Name, vol.NameSuffix),
			Target: vol.Target,
		})
	}

	return mounts, nil
}

// buildPortBindings converts the ports to what Docker actually needs.
func buildPortBindings(a mconfig.ContainerAllocation, ports []string) (network.PortSet, network.PortMap, error) {
	exposedPorts := network.PortSet{}
	portBindings := network.PortMap{}

	// Make sure the amount of ports is correct
	if len(a.Ports) != len(ports) {
		return nil, nil, fmt.Errorf("expected %d ports, received only %d", len(ports), len(a.Ports))
	}

	for i, port := range ports {
		p, err := network.ParsePort(port)
		if err != nil {
			return nil, nil, fmt.Errorf("couldn't parse container port %q: %s", port, err)
		}

		exposedPorts[p] = struct{}{}
		portBindings[p] = []network.PortBinding{
			{
				HostIP:   netip.MustParseAddr("127.0.0.1"),
				HostPort: fmt.Sprintf("%d", a.Ports[i]),
			},
		}
	}

	return exposedPorts, portBindings, nil
}

// GetImageMajorVersion extracts the major version number from a Docker image
// tag. It searches for the first sequence of digits in the tag portion of the
// image name (after the colon) and returns it as an integer. For example:
//   - "postgres:17"      -> 17
//   - "node:v20.1.0"     -> 20
//   - "nginx:1.25-alpine" -> 1
//
// Returns -1 if no version number can be found.
func GetImageMajorVersion(image string) int {
	// Use only the tag portion if a colon is present.
	tag := image
	if idx := strings.LastIndex(image, ":"); idx != -1 {
		tag = image[idx+1:]
	}

	// Find the first run of digits in the tag.
	start := -1
	for i, ch := range tag {
		if ch >= '0' && ch <= '9' {
			if start == -1 {
				start = i
			}
		} else {
			if start != -1 {
				// We have reached the end of the first digit sequence.
				v, _ := strconv.Atoi(tag[start:i])
				return v
			}
		}
	}

	// Handle the case where the tag ends with digits.
	if start != -1 {
		v, _ := strconv.Atoi(tag[start:])
		return v
	}

	return -1
}
