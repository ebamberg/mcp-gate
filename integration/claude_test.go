package integration

import (
	"os"
	"testing"
	"time"

	"github.com/ebamberg/mcp-gate/tests"
)

func TestAddMCPGateToClaudeDesktopConfig(t *testing.T) {

}

func TestBackupClaudeDesktopConfig(t *testing.T) {
	var file *os.File
	var err error
	timestamp := time.Now().GoString()
	{
		file, err = os.CreateTemp("", "mcpgate_unittest_")
		tests.FailOnError(t, err)
		defer file.Close()

		file.WriteString(timestamp)
	}
	defer os.Remove(file.Name())

	err = BackupClaudeDesktopConfig(file.Name())
	tests.FailOnError(t, err)
	defer os.Remove(file.Name() + ".bak")

	tests.AssertFileExist(t, file.Name()+".bak")

	bytes, err := os.ReadFile(file.Name() + ".bak")
	tests.FailOnError(t, err)
	if string(bytes) != timestamp {
		t.Errorf("content of backupfile doesn't match original content")
	}
}
