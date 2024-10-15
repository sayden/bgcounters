package counters

import (
	"encoding/json"
	"testing"
)

func TestParseTemplate(t *testing.T) {
	ct := CounterTemplate{
		Rows:                      7,
		Columns:                   7,
		Mode:                      "template",
		OutputFolder:              "generated",
		PositionNumberForFilename: 0,

		Counters: []Counter{{
			Texts:  []Text{{Settings: Settings{Position: 3}, String: "String"}},
			Images: []Image{{Settings: Settings{Position: 0}, Path: "path/to/image"}}}},

		Prototypes: map[string]CounterPrototype{
			"prototype1": {
				TextsPrototypes: []TextPrototype{
					{StringList: []string{"String1", "String2"}}},
				ImagesPrototypes: []ImagePrototype{{
					Image:    Image{Settings: Settings{Position: 0}},
					PathList: []string{"path/to/image1", "path/to/image2"}}},
				Counter: Counter{
					Texts:  []Text{{Settings: Settings{Position: 3}, String: "String"}},
					Images: []Image{{Settings: Settings{Position: 0}, Path: "path/to/image"}}}}},

		Settings: Settings{
			FontPath:        "BebasNeue-Regular.ttf",
			FontHeight:      10,
			Width:           100,
			Height:          100,
			Margins:         3,
			FontColorS:      "black",
			BackgroundColor: "white",
			ImageScaling:    "fitWidth"},
	}

	byt, err := json.MarshalIndent(ct, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling counter template: %v", err)
	}

	parsed, err := ParseTemplate(byt)
	if err != nil {
		t.Fatalf("Error parsing counter template: %v", err)
	}

	byt, err = json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling counter template: %v", err)
	}

	// fmt.Println(string(byt))
}
