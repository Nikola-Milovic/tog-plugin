package tests

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/startup"
)

//TestMain is here to do the setup needed before all of the tests,
//populates the UnitDataMap for tests
func TestMain(m *testing.M) {
	//os.Stdout, _ = os.Open(os.DevNull)
	log.SetOutput(ioutil.Discard)
	startup.StartUp(true)
	createTempFile()
	code := m.Run()
	os.Exit(code)
}
