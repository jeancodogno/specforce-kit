package agent
import (
	"os"
	"testing"
)
func TestCWD(t *testing.T) {
	cwd, _ := os.Getwd()
	t.Logf("CWD: %s", cwd)
}
