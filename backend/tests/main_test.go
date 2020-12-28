package tests

import (
	"os"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/startup"
)

//TestMain is here to do the setup needed before all of the tests,
//populates the UnitDataMap for tests
func TestMain(m *testing.M) {
	startup.StartUp(true)
	code := m.Run()
	os.Exit(code)
}
