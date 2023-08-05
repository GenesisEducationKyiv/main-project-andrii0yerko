package architecture_test

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestArchitecture(t *testing.T) {
	t.Run("Controller does not depend directrly on infra", func(t *testing.T) {
		archtest.Package(t, "bitcoinrateapp/pkg/app").
			ShouldNotDependDirectlyOn("bitcoinrateapp/pkg/rateclient", "bitcoinrateapp/pkg/storage", "bitcoinrateapp/pkg/email")
	})

	t.Run("Email package is isolated", func(t *testing.T) {
		archtest.Package(t, "bitcoinrateapp/pkg/email").
			ShouldNotDependOn("bitcoinrateapp/pkg/rateclient", "bitcoinrateapp/pkg/storage")
	})
	t.Run("RateClient package is isolated", func(t *testing.T) {
		archtest.Package(t, "bitcoinrateapp/pkg/rateclient").
			ShouldNotDependOn("bitcoinrateapp/pkg/email", "bitcoinrateapp/pkg/storage")
	})

	t.Run("Storage package is isolated", func(t *testing.T) {
		archtest.Package(t, "bitcoinrateapp/pkg/storage").
			ShouldNotDependOn("bitcoinrateapp/pkg/email", "bitcoinrateapp/pkg/rateclient")
	})

	t.Run("Models are independent", func(t *testing.T) {
		archtest.Package(t, "bitcoinrateapp/pkg/storage").
			ShouldNotDependOn(
				"bitcoinrateapp/pkg/rateclient",
				"bitcoinrateapp/pkg/storage",
				"bitcoinrateapp/pkg/email",
				"bitcoinrateapp/pkg/service",
				"bitcoinrateapp/pkg/app",
			)
	})
}
