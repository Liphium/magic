package mservices

import (
	"context"
	"fmt"
	"log"
	"net/netip"
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
}

// CreateContainer finds and removes any existing container with the
// given name, then creates a fresh one from the provided options.
//
// Existing Docker volumes are always preserved so that data survives a
// container re-creation. Returns the ID of the newly created container.
func CreateContainer(ctx context.Context, log *log.Logger, c *client.Client, a mconfig.ContainerAllocation, opts ManagedContainerOptions) (string, error) {
	if opts.Image == "" {
		return "", fmt.Errorf("please specify a proper image")
	}

	existingMounts, err := removeExistingContainer(ctx, log, c, a, opts)
	if err != nil {
		return "", err
	}

	mounts, err := buildMounts(a, opts.Volumes, existingMounts)
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

// removeExistingContainer looks for an existing container with the allocation's
// name, recovers its mounts, and removes it. Returns a map of volume NameSuffix
// -> mount so the new container can reuse the same volumes.
func removeExistingContainer(ctx context.Context, log *log.Logger, c *client.Client, a mconfig.ContainerAllocation, opts ManagedContainerOptions) (map[string]mount.Mount, error) {
	f := make(client.Filters)
	f.Add("name", a.Name)
	summary, err := c.ContainerList(ctx, client.ContainerListOptions{
		Filters: f,
		All:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't list containers: %s", err)
	}

	existingMounts := map[string]mount.Mount{}

	for _, ct := range summary.Items {
		for _, n := range ct.Names {
			if !strings.HasSuffix(n, a.Name) {
				continue
			}

			log.Println("Found existing container, recovering mounts...")
			image, err := recoverMounts(ctx, c, ct.ID, opts.Volumes, existingMounts)
			if err != nil {
				return nil, err
			}

			// Check if the old container is using an image of a different image
			majorCurrent := GetImageMajorVersion(image)
			majorNew := GetImageMajorVersion(opts.Image)
			if majorCurrent == -1 || majorNew == -1 {
				log.Printf("Skipping major version check: unable to parse image versions (%s -> %s)", image, opts.Image)
			} else if majorCurrent != majorNew {
				return nil, fmt.Errorf("major version mismatch for %s: please upgrade or delete your old container before starting the app (%d -> %d)", opts.Image, majorCurrent, majorNew)
			}

			log.Println("Removing old container...")
			if _, err := c.ContainerRemove(ctx, ct.ID, client.ContainerRemoveOptions{
				RemoveVolumes: false,
				Force:         true,
			}); err != nil {
				return nil, fmt.Errorf("couldn't remove existing container: %s", err)
			}
		}
	}

	return existingMounts, nil
}

// recoverMounts inspects a container and indexes its mounts by the matching
// ContainerVolume.NameSuffix into the provided map.
// also returns the name of the image of the container for convenience.
func recoverMounts(ctx context.Context, c *client.Client, containerID string, volumes []ContainerVolume, out map[string]mount.Mount) (string, error) {
	resp, err := c.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return "", fmt.Errorf("couldn't inspect container: %s", err)
	}

	for _, m := range resp.Container.HostConfig.Mounts {
		for _, vol := range volumes {
			if m.Target == vol.Target {
				out[vol.NameSuffix] = m
			}
		}
	}

	return resp.Container.Config.Image, nil
}

// buildMounts constructs the mount list for the new container. Any volume whose
// target was found in existingMounts is reused as-is; otherwise a fresh named
// volume is created using "<containerName>-<nameSuffix>".
func buildMounts(a mconfig.ContainerAllocation, volumes []ContainerVolume, existingMounts map[string]mount.Mount) ([]mount.Mount, error) {
	mounts := make([]mount.Mount, 0, len(volumes))

	for _, vol := range volumes {
		if existing, ok := existingMounts[vol.NameSuffix]; ok {
			mounts = append(mounts, existing)
		} else {
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeVolume,
				Source: fmt.Sprintf("%s-%s", a.Name, vol.NameSuffix),
				Target: vol.Target,
			})
		}
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
