package tmux

import (
	"errors"

	"github.com/sergei-dyshel/claug/utils"
)

type CommandParams struct {
	name  string
	args  utils.Slice[string]
	input string
	err   error
}

type Command interface {
	Params() CommandParams
}

type SendKeys struct {
	Keys    []string
	Literal bool // processes the keys as literal UTF-8 characters
}

func (cmd *SendKeys) Params() (params CommandParams) {
	params.name = "send-keys"
	params.args.AppendIf(cmd.Literal, "-l").Append(cmd.Keys...)
	return
}

type BufferOpts struct {
	BufferName string
}

func (opts BufferOpts) Args() []string {
	if opts.BufferName != "" {
		return []string{"-b", opts.BufferName}
	}
	return []string{}
}

type PasteBuffer struct {
	BufferOpts

	Delete    bool // delete buffer after pasting
	NoReplace bool // do not replace separator
	Bracketed bool // bracketed paste
}

func (cmd *PasteBuffer) Params() (params CommandParams) {
	params.name = "paste-buffer"
	params.args.Values = cmd.BufferOpts.Args()
	params.args.AppendIf(cmd.Delete, "-d").
		AppendIf(cmd.NoReplace, "-r").
		AppendIf(cmd.Bracketed, "-p")
	return
}

type LoadBuffer struct {
	BufferOpts

	Path string
	Text string
}

func (cmd *LoadBuffer) Params() (params CommandParams) {
	params.name = "load-buffer"
	params.args.Values = cmd.BufferOpts.Args()

	if cmd.Path != "" && cmd.Text == "" {
		params.args.Append(cmd.Path)
	} else if cmd.Path == "" && cmd.Text != "" {
		params.args.Append("-")
		params.input = cmd.Text
	} else {
		params.err = errors.New("LoadBuffer: exactly one of Path or Text must be specified")
	}
	return
}

type TargetPaneOpts struct {
	TargetPane string
}

func (opts TargetPaneOpts) Args() []string {
	if opts.TargetPane != "" {
		return []string{"-t", opts.TargetPane}
	}
	return []string{}
}

type CapturePane struct {
	TargetPaneOpts

	StdOut bool
	Join   bool
	Escape bool
	Start  string
	End    string
}

func (cmd *CapturePane) Params() (params CommandParams) {
	params.name = "capture-pane"
	params.args.Values = cmd.TargetPaneOpts.Args()
	params.args.AppendIf(cmd.StdOut, "-p").
		AppendIf(cmd.Escape, "-e").
		AppendIf(cmd.Join, "-J").
		AppendIf(cmd.Start != "", "-S", cmd.Start).
		AppendIf(cmd.End != "", "-E", cmd.End)
	return
}
