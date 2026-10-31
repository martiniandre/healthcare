package tests

import (
	"os"
	"testing"

	"github.com/healthcare/backend/internal/modules/auth"
)

func TestMain(testingMain *testing.M) {
	if initError := auth.InitJWT("auth-module-tests-secret-0123456789abcdef"); initError != nil {
		panic(initError)
	}
	os.Exit(testingMain.Run())
}
