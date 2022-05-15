package tmux

type SendKeys struct {
	Literal bool `opt:"l"` // processes the keys as literal UTF-8 characters
}

type PasteBuffer struct {
	Name      string `opt:"b"` // buffer name
	Delete    bool   `opt:"d"` // delete buffer after pasting
	NoReplace bool   `opt:"r"` // do not replace separator
	Bracketed bool   `opt:"p"` // bracketed paste
}

type LoadBuffer struct {
	Name string `opt:"b"` // buffer name
	Path string // path to file
}

type CapturePane struct {
	Pane   string `opt:"t"` // target pane
	StdOut bool   `opt:"b"` // dump to stdout
	Join   bool   `opt:"J"` // TODO:
	Escape bool   `opt:"e"`
	Start  string `opt:"S"`
	End    string `opt:"E"`
}
