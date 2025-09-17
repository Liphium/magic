package magic

import (
	"os"
	"testing"
	"time"

	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/util"
)

var startSignalChan = make(chan struct{})
var magicTestRunner *mrunner.Runner = nil

// Call this function in any TestMain function if you want some tests to rely on your actual app or database.
//
// When calling this method, Magic will automatically start your app and any required containers, just like when you
// normally run your app.
//
// Please make sure to call the magic.AppStarted() function for the test runner to work properly. You can read its
// comments if you should still have questions about what it does. You can adjust the timeout for exiting in the Magic
// config.
//
// You can provide a profile to make sure you can run multiple tests in multiple packages in parallel. Magic will automatically
// append test- to it, so don't worry about it colliding with profiles you set over the --profile (-p) flag.
//
// The handler will be called once everything is ready.
func PrepareTesting(t *testing.M, profile string, config Config) {
	if profile == "" {
		profile = "default"
	}

	// Start all the containers using Magic
	factory, runner := prepare(config, profile)
	if factory == nil || runner == nil {
		util.Log.Fatal("Couldn't prepare containers with Magic")
		return
	}

	// Load environment
	util.Log.Println("Loading environment...")
	for key, value := range runner.Plan().Environment {
		if err := os.Setenv(key, value); err != nil {
			util.Log.Fatalf("couldn't set environment variable %s: %s", key, err)
		}
	}
	util.Log.Println("Setup finished.")

	// Stop all containers and unlock once testing is done
	defer func() {
		recover()
		factory.Unlock()
		runner.StopContainers()
	}()

	// Start the app
	go func() {
		config.StartFunction()
	}()

	// Wait for the app's start signal
	util.Log.Println("Waiting for start signal...")
	util.Log.Println("If you don't call magic.AppStarted() when your app starts, this will fail.")

	timeoutChan := make(chan struct{})
	go func() {
		if config.TestAppTimeout == nil {
			config.TestAppTimeout = Ptr(time.Second * 10)
		}
		time.Sleep(*config.TestAppTimeout)
		timeoutChan <- struct{}{}
	}()

	select {
	case <-timeoutChan:
		util.Log.Fatalln("Couldn't get start signal in time.")
	case <-startSignalChan:
		util.Log.Println("Signal received. Everything successfully prepared!")
	}

	magicTestRunner = runner
	t.Run()
}

// This function has to be called when your app successfully started.
//
// It's used to start the test once your app is up and running in testing. The test runner does have a timeout of
// 10 seconds though, so if you're app takes longer than that to startup, you can modify the timeout in your Magic
// config.
func AppStarted() {
	startSignalChan <- struct{}{}
}

// Get the current runner active while testing.
//
// For this to work, please make sure you call PrepareTesting in your TestMain function.
func GetTestRunner() *mrunner.Runner {
	return magicTestRunner
}
