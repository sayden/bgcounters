package counters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCounterFilename(t *testing.T) {
	tests := []struct {
		name           string
		counter        Counter
		position       int
		suffix         string
		filenumber     int
		filenamesInUse map[string]bool
		expected       string
	}{
		{
			name: "Basic case",
			counter: Counter{
				Texts: []Text{
					{
						Settings: Settings{Position: 0},
						String:   "Test",
					},
				},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     -1,
			filenamesInUse: map[string]bool{},
			expected:       "Test suffix.png",
		},
		{
			name: "With Extra Title",
			counter: Counter{
				Texts: []Text{
					{
						Settings: Settings{Position: 0},
						String:   "Test",
					},
				},
				Extra: &Extra{
					Title: "ExtraTitle",
				},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     -1,
			filenamesInUse: map[string]bool{},
			expected:       "Test ExtraTitle suffix.png",
		},
		{
			name: "With Title Position",
			counter: Counter{
				Texts: []Text{
					{
						Settings: Settings{Position: 0},
						String:   "Test",
					},
					{
						Settings: Settings{Position: 1},
						String:   "TitlePosition",
					},
				},
				Extra: &Extra{
					TitlePosition: newInt(1),
				},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     -1,
			filenamesInUse: map[string]bool{},
			expected:       "TitlePosition suffix.png",
		},
		{
			name: "Filename in use",
			counter: Counter{
				Texts: []Text{
					{
						Settings: Settings{Position: 0},
						String:   "Test",
					},
				},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     1,
			filenamesInUse: map[string]bool{"Test suffix": true},
			expected:       "Test suffix 001.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.counter.GetCounterFilename(tt.position, tt.suffix, tt.filenumber, tt.filenamesInUse)
			if got != tt.expected {
				t.Errorf("GetCounterFilename() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestApplyCounterScaling(t *testing.T) {
	counter := Counter{
		Settings: Settings{
			Width: 300,
		},
		Images: []Image{{Scale: 0.5}},
		Texts:  []Text{{Settings: Settings{FontHeight: 10}}},
	}
	ApplyCounterScaling(&counter, 2)

	assert.Equal(t, 600, counter.Settings.Width)
	assert.Equal(t, float64(1), counter.Images[0].Scale)
	assert.Equal(t, 20.0, counter.Texts[0].Settings.FontHeight)
}

func newInt(i int) *int {
	return &i
}
