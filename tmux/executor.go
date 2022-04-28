package tmux

import (
	"fmt"
	"io"
	"os/exec"
)

type Executor interface {
	Run(command Command) ([]byte, error)
}

type CmdLineExecutor struct{}

func (*CmdLineExecutor) runImpl(params CommandParams) ([]byte, error) {
	err := params.err
	if err != nil {
		return nil, fmt.Errorf("error generating arguments for command: %w", err)
	}
	args := []string{params.name}
	args = append(args, params.args.Values...)
	execCmd := exec.Command("tmux", args...)
	if params.input != "" {
		stdin, err := execCmd.StdinPipe()
		if err != nil {
			return nil, fmt.Errorf("could not create stdin pipe: %w", err)
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, params.input)
		}()
	}
	if err != nil {
		// TODO: add full command to error msg
		return nil, fmt.Errorf("failed to run tmux command: %w", err)
	}
	return execCmd.Output()
}

func (executor *CmdLineExecutor) Run(command Command) ([]byte, error) {
	params := command.Params()
	out, err := executor.runImpl(params)
	if err != nil {
		return nil, fmt.Errorf("failed to run command '%s': %w", params.name, err)
	}
	return out, err
}

func InsertText(executor Executor, text string, bracketedPaste bool) error {
	_, err := executor.Run(&LoadBuffer{Text: text})
	if err != nil {
		return fmt.Errorf("failed to load buffer: %w", err)
	}
	_, err = executor.Run(&PasteBuffer{Delete: true, NoReplace: true, Bracketed: bracketedPaste})
	if err != nil {
		return fmt.Errorf("failed to paste buffer: %w", err)
	}
	return nil
}
