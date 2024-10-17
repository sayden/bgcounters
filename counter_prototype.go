package counters

import (
	"github.com/pkg/errors"
	deepcopy "github.com/qdm12/reprint"
)

type CounterPrototype struct {
	Counter
	ImagePrototypes []ImagePrototype `json:"image_prototypes,omitempty"`
	TextPrototypes  []TextPrototype  `json:"text_prototypes,omitempty"`
}

type ImagePrototype struct {
	Image
	PathList []string `json:"path_list"`
}

type TextPrototype struct {
	Text
	StringList []string `json:"string_list"`
}

func (p *CounterPrototype) ToCounters() ([]Counter, error) {
	cts := make([]Counter, 0)

	// You can prototype texts and images, so one of the two must be present, get their length
	length := 0
	if len(p.TextPrototypes) > 0 && len(p.TextPrototypes[0].StringList) > 0 {
		length = len(p.TextPrototypes[0].StringList)
		if len(p.ImagePrototypes) > 0 && len(p.ImagePrototypes[0].PathList) != length {
			return nil, errors.New("the number of images and texts prototypes must be the same")
		}
	} else if len(p.ImagePrototypes) > 0 && len(p.ImagePrototypes[0].PathList) > 0 {
		length = len(p.ImagePrototypes[0].PathList)
		if len(p.TextPrototypes) > 0 && len(p.TextPrototypes) != length {
			return nil, errors.New("the number of images and texts prototypes must be the same")
		}
	} else {
		return nil, errors.New("no prototypes found in the counter template")
	}

	for i := 0; i < length; i++ {
		var newCounter Counter
		if err := deepcopy.FromTo(p.Counter, &newCounter); err != nil {
			return nil, err
		}

		if p.TextPrototypes != nil {
			for _, textPrototype := range p.TextPrototypes {
				originalText := Text{}
				if err := deepcopy.FromTo(textPrototype.Text, &originalText); err != nil {
					return nil, err
				}
				originalText.String = textPrototype.StringList[i]
				newCounter.Texts = append(newCounter.Texts, originalText)
			}
		}

		if p.ImagePrototypes != nil {
			for _, imagePrototype := range p.ImagePrototypes {
				originalImage := Image{}
				if err := deepcopy.FromTo(imagePrototype.Image, &originalImage); err != nil {
					return nil, err
				}
				originalImage.Path = imagePrototype.PathList[i]
				newCounter.Images = append(newCounter.Images, originalImage)
			}
		}
		cts = append(cts, newCounter)

	}

	return cts, nil
}
