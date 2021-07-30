package update

import (
	"context"
	"fmt"
	"runtime"
)

func UpdateCommand() string {
	if isUnderHomebrew() {
		return "brew upgrade fing"
	}

	if isUnderNpm() {
		return fmt.Sprintf("npm i -g %s", npmPackage)
	}

	if isNpmAvailable() {
		return fmt.Sprintf("npm i -g %s", npmPackage)
	}

	if runtime.GOOS == "windows" {
		return "iwr https://fing.ir/install.ps1 -useb | iex"
	}

	return "curl -fsSL https://fing.ir/install.sh | sh"
}

// PerformSelfUpdate downloads fing binary
func PerformSelfUpdate(ctx context.Context, version string) error {
	// TODO: implement this
	return nil
}
