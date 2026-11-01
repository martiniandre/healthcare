package integration

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestMain(testingMain *testing.M) {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupts
		environment.teardown()
		os.Exit(1)
	}()

	exitCode := runTestsWithDatabaseTeardown(testingMain)
	signal.Stop(interrupts)
	os.Exit(exitCode)
}

func runTestsWithDatabaseTeardown(testingMain *testing.M) int {
	defer environment.teardown()
	return testingMain.Run()
}
