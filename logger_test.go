package artemis_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/MovieStoreGuy/artemis"
)

func TestCreateLogger(t *testing.T) {
	i := artemis.GetInstance()
	if i == nil {
		t.Fatal("Created a nil object")
	}
	if i != artemis.GetInstance() {
		t.Fatal("Should create a singleton")
	}
}

func TestLogger(t *testing.T) {
	i := artemis.GetInstance()
	if i == nil {
		t.Fatal("Created a nil object")
	}
	b := &bytes.Buffer{}
	i.Set(artemis.Debug, b)
	i.Log(artemis.Entry{
		Level: artemis.Info,
		Data:  "What is it good for",
	})
	i.Start()
	// testing starting the artemis twice
	i.Start()
	i.Log(artemis.Entry{
		Level: artemis.Info,
		Data:  "What is it good for",
	})
	i.Stop()
	expected := fmt.Sprintf("^[Trace] .* %s$", "What is it good for")
	if regexp.MustCompile(expected).MatchString(string(b.Bytes())) {
		t.Log("Expected:", expected)
		t.Log("Given:", string(b.Bytes()))
		t.Fatal("Incorrect details logged")
	}
	i.Stop()
}

func TestLogLevel(t *testing.T) {
	var level artemis.Level = -1
	if level.String() != "Unknown" {
		t.Fatal("Should report an unknown level")
	}
}
