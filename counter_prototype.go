package counters

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"github.com/qdm12/reprint"
)

type CounterPrototype struct {
	Counter

	ImagesPrototypes []ImagePrototype `json:"image_prototypes,omitempty"`
	TextsPrototypes  []TextPrototype  `json:"text_prototypes,omitempty"`
}

type ImagePrototype struct {
	Image

	PathList []string `json:"path_list"`
}

type TextPrototype struct {
	Text

	StringList []string `json:"string_list"`
}

func ParsePrototypedTemplate(counterTemplate *CounterTemplate) (*CounterTemplate, error) {
	// JSON counters to Counters
	newTemplate, err := JsonPrototypeToTemplate(counterTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to convert a counter template into another counter template")
	}
	newTemplate.Scaling = 1

	byt, err := json.Marshal(newTemplate)
	if err != nil {
		return nil, err
	}

	newTemplate, err = ParseTemplate(byt)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse JSON file")
	}

	return newTemplate, nil
}

func JsonPrototypeToTemplate(ct *CounterTemplate) (t *CounterTemplate, err error) {
	// JSON counters to Counters, check Prototype in CounterTemplate
	if ct.Prototypes != nil {
		if ct.Counters == nil {
			ct.Counters = make([]Counter, 0)
		}

		// sort prototypes by name, to ensure consistent output filenames this is a small
		// inconvenience, because iterating over maps in Go returns keys in random order
		names := make([]string, 0, len(ct.Prototypes))
		for name := range ct.Prototypes {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, prototypeName := range names {
			counter := ct.Prototypes[prototypeName]

			// You can prototype texts and images, so one of the two must be present, get their length
			length := 0
			if len(counter.TextsPrototypes) > 0 && len(counter.TextsPrototypes[0].StringList) > 0 {
				length = len(counter.TextsPrototypes[0].StringList)
				if len(counter.ImagesPrototypes) > 0 && len(counter.ImagesPrototypes[0].PathList) != length {
					return nil, errors.New("the number of images and texts prototypes must be the same")
				}
			} else if len(counter.ImagesPrototypes) > 0 && len(counter.ImagesPrototypes[0].PathList) > 0 {
				length = len(counter.ImagesPrototypes[0].PathList)
				if len(counter.TextsPrototypes) > 0 && len(counter.TextsPrototypes) != length {
					return nil, errors.New("the number of images and texts prototypes must be the same")
				}
			} else {
				return nil, errors.New("no prototypes found in the counter template")
			}

			for i := 0; i < length; i++ {
				var newCounter Counter
				if err = reprint.FromTo(counter.Counter, &newCounter); err != nil {
					return nil, err
				}

				if counter.TextsPrototypes != nil {
					for _, textPrototype := range counter.TextsPrototypes {
						originalText := Text{}
						if err = reprint.FromTo(textPrototype.Text, &originalText); err != nil {
							return nil, err
						}
						originalText.String = textPrototype.StringList[i]
						newCounter.Texts = append(newCounter.Texts, originalText)
					}
				}

				if counter.ImagesPrototypes != nil {
					for _, imagePrototype := range counter.ImagesPrototypes {
						originalImage := Image{}
						if err = reprint.FromTo(imagePrototype.Image, &originalImage); err != nil {
							return nil, err
						}
						originalImage.Path = imagePrototype.PathList[i]
						newCounter.Images = append(newCounter.Images, originalImage)
					}
				}

				ct.Counters = append(ct.Counters, newCounter)
			}
		}

		ct.Prototypes = nil

		return ct, nil
	}

	return ct, nil
}
