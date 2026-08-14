package mrunner

import (
	"context"
	"slices"
	"strconv"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/Liphium/magic/v3/util"
)

// Get the current plan (might not be set yet, call GeneratePlan first)
func (r *Runner) Plan() *mconfig.Plan {
	return r.plan
}

// Deploy the containers for the magic context
func (r *Runner) GeneratePlan() *mconfig.Plan {
	if r.ctx == nil {
		util.Log.Fatalln("no context set")
	}

	// Set basic stuff
	r.plan.AppName = r.ctx.AppName()
	r.plan.Profile = r.ctx.Profile()

	// Collect all the ports that should be allocated (also for the service drivers obv)
	portsToAllocate := r.ctx.Ports()
	reusedPorts := []uint{}
	containerMap := map[string]mconfig.ContainerAllocation{}
	for _, driver := range r.ctx.Services() {
		if _, ok := containerMap[driver.GetUniqueId()]; ok {
			util.Log.Fatalln("ERROR: You can't create multiple drivers of the same type at the moment.")
		}

		requiredPorts := driver.GetRequiredPorts()
		alloc := mconfig.ContainerAllocation{
			Name:  mconfig.PlannedContainerName(r.plan, driver),
			Ports: []uint{},
		}

		// See if there is already a container (if yes, try find the previously allocated ports and reuse)
		inspect, err := mservices.FindContainer(context.Background(), r.client, alloc.Name)
		if err == nil && inspect != nil && inspect.Config != nil && inspect.HostConfig != nil {
			currentPorts := make([]uint, len(requiredPorts))

			// Append all of the existing port bindings to the list, they can then be re-used (order does not matter as we're talking about the ports on the host)
			found := []string{}
			for port, bindings := range inspect.HostConfig.PortBindings {
				if !slices.Contains(requiredPorts, port.String()) {
					continue
				}
				found = append(found, port.String())

				// Magic never creates more than one binding, if there is more the container has been modified
				if len(bindings) != 1 {
					if mconfig.VerboseLogging {
						util.Log.Println("Previous container broken for", alloc.Name+":", "more port bindings than one on the host for port", port.String())
					}
					goto error
				}

				// Find index for the port in the drivers port list
				i := slices.Index(requiredPorts, port.String())

				for _, binding := range bindings {
					port, err := strconv.Atoi(binding.HostPort)
					if err != nil {
						if mconfig.VerboseLogging {
							util.Log.Println("Previous container broken for", alloc.Name, "and port", port)
						}
						goto error
					}

					currentPorts[i] = uint(port)
				}
			}

			if len(found) != len(requiredPorts) {
				if mconfig.VerboseLogging {
					util.Log.Println("Previous container broken for", alloc.Name+":", "found", len(found), "required", len(requiredPorts))
				}
				goto error
			}

			// When the ports aren't the same length, just clear them again
			if len(currentPorts) != len(driver.GetRequiredPorts()) {
				if mconfig.VerboseLogging {
					util.Log.Println("Previous container broken for", alloc.Name+":", "allocated", len(currentPorts))
				}
				goto error
			}

			// Add the ports to the allocation and the list of ports to allocate
			alloc.Ports = currentPorts
			if inspect.State.Running {
				reusedPorts = append(reusedPorts, currentPorts...)
			} else {
				portsToAllocate = append(portsToAllocate, currentPorts...)
			}
			containerMap[driver.GetUniqueId()] = alloc
			continue
		}
	error:

		for range len(requiredPorts) {

			// Make sure we're not allocating a port that's already taken
			currentPort := util.RandomPort(DefaultStartPort, DefaultEndPort)
			for slices.Contains(portsToAllocate, currentPort) && !util.ScanPort(currentPort) {
				currentPort = util.RandomPort(DefaultStartPort, DefaultEndPort)
			}

			// Allocate one of the default ports for the container
			portsToAllocate = append(portsToAllocate, currentPort)
			alloc.Ports = append(alloc.Ports, currentPort)
		}

		containerMap[driver.GetUniqueId()] = alloc
	}

	// Prepare all of the ports
	allocatedPorts := map[uint]uint{}
	if len(portsToAllocate) > 0 {
		for _, port := range portsToAllocate {

			// Generate a new port in case the current one is taken
			toAllocate := port
			for !util.ScanPort(toAllocate) || slices.Contains(portsToAllocate, toAllocate) || slices.Contains(reusedPorts, toAllocate) {
				toAllocate = util.RandomPort(DefaultStartPort, DefaultEndPort)
			}

			// Add the port to the plan
			allocatedPorts[port] = toAllocate
		}
	}
	for _, port := range reusedPorts {
		allocatedPorts[port] = port
	}

	// Load into plan
	r.plan.Containers = containerMap
	r.plan.AllocatedPorts = allocatedPorts

	// Add all the services to the plan
	r.plan.Services = map[string]string{}
	for _, driver := range r.ctx.Services() {
		data, err := driver.Save()
		if err != nil {
			util.Log.Fatalln("couldn't persist service driver of type", driver.GetUniqueId()+":", err)
		}
		r.plan.Services[driver.GetUniqueId()] = data
	}

	// Generate the environment variables and add to plan
	environment := map[string]string{}
	if r.ctx.Environment() != nil {
		environment = r.ctx.Environment().Generate()
	}
	r.plan.Environment = environment
	return r.plan
}
