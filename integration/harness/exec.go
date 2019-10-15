package integration

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	// compileMtx guards access to the executable path so that the project is
	// only compiled once.
	compileMtx sync.Mutex

	// executablePath is the path to the compiled executable. This is an empty
	// string until the initial compilation. It should not be accessed directly;
	// use the poetExecutablePath() function instead.
	executablePath string
)

// nodeExecutablePath returns a path to the mock node server executable.
// To ensure the code tests against the most up-to-date version, this method
// compiles a mock node server the first time it is called. After that, the
// generated binary is used for subsequent requests.
func nodeExecutablePath(baseDir string) (string, error) {
	compileMtx.Lock()
	defer compileMtx.Unlock()

	// If mock node has already been compiled, just use that.
	if len(executablePath) != 0 {
		return executablePath, nil
	}

	// Build mock node and output an executable in basedir path.
	outputPath := filepath.Join(baseDir, "mocknode")
	if runtime.GOOS == "windows" {
		outputPath += ".exe"
	}

	// TODO  add node repo
	// TODO "github.com/spacemeshos/poet"
	cmd := exec.Command("go", "build", "-i", "-o", outputPath)

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to build mock node: %v", err)
	}

	// Save executable path so future calls do not recompile.
	executablePath = outputPath
	return executablePath, nil
}
