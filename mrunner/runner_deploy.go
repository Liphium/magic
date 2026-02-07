package mrunner

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/util"
	_ "github.com/lib/pq"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
)

// Deploy all the containers nessecary for the application
func (r *Runner) Deploy(deleteContainers bool) error {
	ctx := context.Background()

	// Make sure the Docker connection is working
	_, err := r.client.Info(ctx, client.InfoOptions{})
	if client.IsErrConnectionFailed(err) {
		return fmt.Errorf("please make sure Docker is running, and that Magic (or the Go toolchain) has access to it. (%s)", err)
	}

	// Clear all state in case wanted
	if deleteContainers {
		util.Log.Println("Deleting all containers and volumes...")
		if err := r.DeleteEverything(); err != nil {
			return fmt.Errorf("couldn't clear state: %v", err)
		}
	}

	// Pull all of the images in case they are not downloaded yet
	if err := r.pullServiceImages(ctx); err != nil {
		return err
	}

	// Start all of the service containers
	if err := r.startServiceContainers(ctx); err != nil {
		return err
	}

	// Load environment variables into current application
	for key, value := range r.plan.Environment {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("couldn't set environment variable %s: %s", key, err)
		}
	}

	util.Log.Println("Deployment finished.")
	return nil
}

// Pull all of the images for the services the runner has registered
func (r *Runner) pullServiceImages(ctx context.Context) error {
	for _, driver := range r.services {
		image := driver.GetImage()

		// Check if the image exists locally
		_, err := r.client.ImageInspect(ctx, image)
		if err != nil {

			// Image not found, need to pull it
			util.Log.Println("Pulling image", image+"...")

			reader, err := r.client.ImagePull(ctx, image, client.ImagePullOptions{})
			if err != nil {
				return fmt.Errorf("couldn't pull image %s: %s", image, err)
			}
			defer reader.Close()

			// Track progress with updates every second
			lastUpdate := time.Now()
			buf := make([]byte, 1024)
			for {
				n, err := reader.Read(buf)
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					return fmt.Errorf("error while pulling image %s: %s", image, err)
				}

				// Print progress update every second
				if time.Since(lastUpdate) >= time.Second {
					util.Log.Println("Downloading", image+"...")
					lastUpdate = time.Now()
				}

				if n == 0 {
					break
				}
			}

			util.Log.Println("Successfully pulled image", image)
		}
	}

	return nil
}

// Create all the service containers and start them + wait until healthy and initialize
func (r *Runner) startServiceContainers(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(r.services))

	// Deploy all service containers in parallel
	for _, driver := range r.services {
		wg.Add(1)
		go func(driver mconfig.ServiceDriver) {
			defer wg.Done()

			name := r.plan.Containers[driver.GetUniqueId()].Name

			// Generate a proper port list from the allocated ports of the plan
			containerPorts := []uint{}
			for _, port := range r.plan.Containers[driver.GetUniqueId()].Ports {
				containerPorts = append(containerPorts, r.plan.AllocatedPorts[port])
			}

			// Create the container using the driver
			containerID, err := driver.CreateContainer(ctx, r.client, mconfig.ContainerAllocation{
				Name:  name,
				Ports: containerPorts,
			})
			if err != nil {
				errChan <- fmt.Errorf("couldn't create container for service %s: %s", driver.GetUniqueId(), err)
				return
			}

			// Start the container
			if _, err := r.client.ContainerStart(ctx, containerID, client.ContainerStartOptions{}); err != nil {
				errChan <- fmt.Errorf("couldn't start container for service %s: %s", driver.GetUniqueId(), err)
				return
			}

			// Monitor health until the container is ready
			util.Log.Println("Waiting for", name, "to be healthy...")
			containerInfo := mconfig.ContainerInformation{
				ID:    containerID,
				Name:  name,
				Ports: r.plan.Containers[driver.GetUniqueId()].Ports,
			}

			for {
				healthy, err := driver.IsHealthy(ctx, r.client, containerInfo)
				if err != nil {
					errChan <- fmt.Errorf("couldn't check health for service %s: %s", driver.GetUniqueId(), err)
					return
				}
				if healthy {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
			time.Sleep(200 * time.Millisecond) // Some extra time, some services are a little weird with healthy state

			// Initialize the container
			if err := driver.Initialize(ctx, r.client, containerInfo); err != nil {
				errChan <- fmt.Errorf("couldn't initialize service %s: %s", driver.GetUniqueId(), err)
				return
			}

			util.Log.Println("Service", name, "is ready")
		}(driver)
	}

	// Wait for all services to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to iterate over containers and execute a callback for each found container
func (r *Runner) forEachContainer(ctx context.Context, callback func(service string, containerID string, container mconfig.ContainerAllocation) error) error {
	for service, container := range r.plan.Containers {

		// Try to find the container for the type
		f := make(client.Filters)
		f.Add("name", container.Name)
		summary, err := r.client.ContainerList(ctx, client.ContainerListOptions{
			Filters: f,
		})
		if err != nil {
			util.Log.Fatalln("Couldn't list containers:", err)
		}
		containerId := ""
		for _, c := range summary.Items {
			for _, n := range c.Names {
				if strings.Contains(n, container.Name) {
					containerId = c.ID
				}
			}
		}

		// If there is no container, nothing to do
		if containerId == "" {
			continue
		}

		// Execute the callback with the container ID, container info, and key
		if err := callback(service, containerId, container); err != nil {
			return err
		}
	}

	return nil
}

// Stop all containers
func (r *Runner) StopContainers() error {
	ctx := context.Background()
	return r.forEachContainer(ctx, func(_, containerID string, container mconfig.ContainerAllocation) error {

		// Stop the container
		util.Log.Println("Stopping container", container.Name+"...")
		if _, err := r.client.ContainerStop(ctx, containerID, client.ContainerStopOptions{}); err != nil {
			return fmt.Errorf("Couldn't stop database container: %v", err)
		}

		return nil
	})
}

// Delete all containers + their attached volumes and reset all state
func (r *Runner) DeleteEverything() error {
	ctx := context.Background()
	return r.forEachContainer(ctx, func(_ string, containerID string, container mconfig.ContainerAllocation) error {

		// Get all the attached volumes to delete them manually
		containerInfo, err := r.client.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
		if err != nil {
			util.Log.Println("Warning: Couldn't inspect container:", err)
		}
		var volumeNames []string
		if containerInfo.Container.Mounts != nil {
			for _, mnt := range containerInfo.Container.Mounts {
				if mnt.Type == mount.TypeVolume && mnt.Name != "" {
					volumeNames = append(volumeNames, mnt.Name)
				}
			}
		}

		// Delete the container
		util.Log.Println("Deleting container", container.Name+"...")
		if _, err := r.client.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{
			RemoveVolumes: false,
			Force:         true,
		}); err != nil {
			return fmt.Errorf("Couldn't delete database container: %v", err)
		}

		// Delete all the attached volumes
		for _, volumeName := range volumeNames {
			util.Log.Println("Deleting volume", volumeName+"...")
			if _, err := r.client.VolumeRemove(ctx, volumeName, client.VolumeRemoveOptions{
				Force: true,
			}); err != nil {
				return fmt.Errorf("Couldn't delete volume %s: %v", volumeName, err)
			}
		}

		return nil
	})
}

// Clear the content of all tables from databases, at runtime
func (r *Runner) DropTables() error {
	ctx := context.Background()
	return r.forEachContainer(ctx, func(service, containerID string, container mconfig.ContainerAllocation) error {
		driver, ok := mconfig.GetDriver(service)
		if !ok {
			return fmt.Errorf("couldn't find service driver for service type: %s", service)
		}

		// Convert the ports
		containerPorts := []uint{}
		for _, port := range container.Ports {
			containerPorts = append(containerPorts, r.plan.AllocatedPorts[port])
		}

		if err := driver.HandleInstruction(ctx, r.client, mconfig.ContainerInformation{
			ID:    containerID,
			Name:  container.Name,
			Ports: containerPorts,
		}, mconfig.InstructionDropTables); err != nil {
			return fmt.Errorf("couldn't drop tables: %v", err)
		}

		return nil
	})
}

// Delete all database tables from databases, at runtime
func (r *Runner) ClearTables() error {
	ctx := context.Background()
	return r.forEachContainer(ctx, func(service, containerID string, container mconfig.ContainerAllocation) error {
		driver, ok := mconfig.GetDriver(service)
		if !ok {
			return fmt.Errorf("couldn't find service driver for service type: %s", service)
		}

		// Convert the ports
		containerPorts := []uint{}
		for _, port := range container.Ports {
			containerPorts = append(containerPorts, r.plan.AllocatedPorts[port])
		}

		if err := driver.HandleInstruction(ctx, r.client, mconfig.ContainerInformation{
			ID:    containerID,
			Name:  container.Name,
			Ports: containerPorts,
		}, mconfig.InstructionClearTables); err != nil {
			return fmt.Errorf("couldn't clear tables: %v", err)
		}

		return nil
	})
}
