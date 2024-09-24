package output

import (
	"fmt"
	"testing"

	csvcounters "github.com/sayden/counters/input"
)

func TestGetBuildFileDataForCounters(t *testing.T) {
	template, err := csvcounters.ReadCSVCounters("../../counters.csv")
	if err != nil {
		t.Fatal(err)
	}

	buildData, err := GetVassalDataForCounters(template, "../../template.xml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(buildData)
}
