package config

type Config struct {
	Interpreters []Interpreter `json:"interpreters" jsonschema_description:"interpreter definitions"`

	Common Options `json:"common,omitempty" jsonschema_description:"default options for all interpreters"`
}

type History struct {
	Disabled *bool `json:"disabled,omitempty" jsonschema_description:"do not save history"`
}

type Options struct {
	History History `json:"history,omitempty" jsonschema_description:"history related options"`

	BracketedPaste *bool `json:"bracketed-paste,omitempty" jsonschema_description:"use bracketed paste when inserting text"`

	Selector *string `json:"selector,omitempty" jsonschema_description:"TODO:"`
}

type Interpreter struct {
	Name string `json:"name" jsonschema_description:"name of the interpreter"`

	Prompts []string `json:"prompts" jsonschema_description:"list of regexes to match prompt"`

	Options Options `json:"options,omitempty" jsonschema_description:"command options specific to this interpreter"`
}
