package tmux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	require.Equal(t, []string{"capture-pane", "-J", "-S", "50", "-E", "60"}, Serialize(
		&CapturePane{Join: true, Start: "50", End: "60"},
	))

	require.Equal(t, []string{"load-buffer", "-b", "buffer", "path"}, Serialize(
		&LoadBuffer{Name: "buffer", Path: "path"},
	))
}
