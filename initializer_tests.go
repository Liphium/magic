package magic

import (
	"testing"

	"github.com/Liphium/magic/mconfig"
)

// Call this function if you want some tests to rely on anything Magic can set up for you.
//
// In a case where you want to maybe run this test in parallel with other tests that include Magic, make sure to set
// the profile to something different for both of them so they don't collide. Otherwise they'll always be waiting for
// each other to be executed, and running them in parallel will become pointless.
//
// If you don't need tests to run in parallel, you could even set the profile to an empty string. In that case, Magic will
// use the default profile: default. However, we recommend you try to use different profiles for different tests, anyway,
// since you might want to run your tests in parallel in the future.
//
// Test profiles use the same system as a profile you can pass over the --profile (-p) flag. But Magic automatically appends
// the test- prefix when you use this method, so you don't have to worry about your profile choice colliding with any other
// profile you may use outside of tests.
//
// The handler will be called once everything is ready. No more than one handler can run at once under one profile.
func TestRunner(t *testing.T, config Config, profile string, handler func(*testing.T, *mconfig.Plan)) {
	if profile == "" {
		profile = "default"
	}

	// Start all the containers using Magic
	factory, runner := prepare(config, profile)
	if factory == nil || runner == nil {
		t.Fatal("Couldn't prepare containers with Magic")
		return
	}

	// Handle any potential crashes and stop all the containers
	defer func() {
		recover()
		factory.Unlock()
		runner.StopContainers()
	}()

	handler(t, runner.Plan())
}
