package mrunner

import (
	"slices"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/util"
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
	currentPort := DefaultStartPort
	containerMap := map[string]mconfig.ContainerAllocation{}
	for _, driver := range r.ctx.Services() {
		if _, ok := containerMap[driver.GetUniqueId()]; ok {
			util.Log.Fatalln("ERROR: You can't create multiple drivers of the same type at the moment.")
		}

		alloc := mconfig.ContainerAllocation{
			Name:  mconfig.PlannedContainerName(r.plan, driver),
			Ports: []uint{},
		}

		for range driver.GetRequiredPortAmount() {

			// Make sure we're not allocating a port that's already taken
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
	if len(portsToAllocate) >= 0 {
		for _, port := range portsToAllocate {

			// Generate a new port in case the current one is taken
			toAllocate := port
			for slices.Contains(portsToAllocate, port) && !util.ScanPort(toAllocate) {
				toAllocate = util.RandomPort(DefaultStartPort, DefaultEndPort)
			}

			// Add the port to the plan
			allocatedPorts[port] = toAllocate
		}
	}

	// Load into plan
	r.plan.Containers = containerMap
	r.plan.AllocatedPorts = allocatedPorts

	// Generate the environment variables and add to plan
	environment := map[string]string{}
	if r.ctx.Environment() != nil {
		environment = r.ctx.Environment().Generate()
	}
	r.plan.Environment = environment
	return r.plan
}
