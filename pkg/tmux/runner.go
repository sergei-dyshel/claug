package tmux

import (
	"io"
	"os/exec"

	"github.com/sergei-dyshel/claug/internal/utils"

	"github.com/alessio/shellescape"
)

type cmdRunner interface {
	run(args []string, input *string) ([]byte, error)
}

type cmdLineRunner struct{}

var CmdLine = &cmdLineRunner{}

func (*cmdLineRunner) run(args []string, input *string) ([]byte, error) {
	execCmd := exec.Command("tmux", args...)
	if input != nil {
		stdin, err := execCmd.StdinPipe()
		if err != nil {
			return nil, utils.Wrapf(err, "could not create stdin pipe")
		}
		go func() {
			defer stdin.Close()
			// TODO: handle err
			_, err := io.WriteString(stdin, *input)
			utils.AssertErr(err)
		}()
	}
	out, err := execCmd.Output()
	if err != nil {
		// TODO: add full command to error msg
		return nil, utils.Wrapf(err, "tmux command failed")
	}
	return out, nil
}

func wrapErr(err error, args []string) error {
	return utils.Wrapf(err, "failed to run command: %s", shellescape.QuoteCommand(args))
}

func Run(runner cmdRunner, cmd any, args ...string) (out []byte, err error) {
	out, err = runner.run(append(Serialize(cmd), args...), nil /* input */)
	if err != nil {
		err = wrapErr(err, args)
	}
	return
}

func RunInput(runner cmdRunner, input string, cmd any, args ...string) (out []byte, err error) {
	out, err = runner.run(append(Serialize(cmd), args...), &input)
	if err != nil {
		err = wrapErr(err, args)
	}
	return
}

func InsertText(runner cmdRunner, text string, bracketedPaste bool) error {
	_, err := RunInput(runner, text, &LoadBuffer{})
	if err != nil {
		return utils.Wrapf(err, "failed to load buffer")
	}
	_, err = Run(runner, &PasteBuffer{Delete: true, NoReplace: true, Bracketed: bracketedPaste})
	if err != nil {
		return utils.Wrapf(err, "failed to paste buffer")
	}
	return nil
}
