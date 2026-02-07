package mrunner

import (
	"fmt"
	"slices"

	"github.com/Liphium/magic/v2/integration"
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

	// Collect all the ports that should be allocated (also for the service drivers obv)
	portsToAllocate := r.ctx.Ports()
	startPort := DefaultStartPort
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
			for slices.Contains(portsToAllocate, startPort) {
				startPort++
			}

			// Allocate one of the default ports for the container
			portsToAllocate = append(portsToAllocate, startPort)
			startPort++
		}

		containerMap[driver.GetUniqueId()] = alloc
	}

	// Prepare all of the ports
	allocatedPorts := map[uint]uint{}
	if len(portsToAllocate) >= 0 {
		for _, port := range portsToAllocate {

			// Generate a new port in case the current one is taken
			toAllocate := port
			if integration.ScanPort(port) {
				var err error
				toAllocate, err = scanForOpenPort()
				if err != nil {
					util.Log.Fatalln("Couldn't find open port for", port, ":", err)
				}
			}

			// Add the port to the plan
			allocatedPorts[port] = toAllocate
		}
	}

	// Load into plan
	r.plan = &mconfig.Plan{
		AppName:        r.ctx.AppName(),
		Profile:        r.ctx.Profile(),
		Containers:     containerMap,
		AllocatedPorts: allocatedPorts,
	}

	// Generate the environment variables and add to plan
	environment := map[string]string{}
	if r.ctx.Environment() != nil {
		environment = r.ctx.Environment().Generate()
	}
	r.plan.Environment = environment
	return r.plan
}

// Scan for an open port in the default range
func scanForOpenPort() (uint, error) {
	openPort, err := integration.ScanForOpenPort(DefaultStartPort, DefaultEndPort)
	if err != nil {
		return 0, fmt.Errorf("couldn't find open port: %e", err)
	}
	return openPort, err
}
