package i3ipc

import (
	"testing"
)

func TestCommand(t *testing.T) {
	ipc, _ := GetIPCSocket()
	defer ipc.Close()

	// `exec /bin/true` is a good NOP operation for testing
	success, err := Command("exec /bin/true", ipc)
	if !success {
		t.Error("Unsuccessful command.")
	}
	if err != nil {
		t.Errorf("An error occurred during command: %v", err)
	}
}