package counters

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
