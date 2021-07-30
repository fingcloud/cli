package update

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/config"
)

const (
	npmPackage = "@fingcloud/cli"
)

func isUnderHomebrew() bool {
	fingBinary, err := os.Executable()
	if err != nil {
		return false
	}

	brewExe, err := exec.LookPath("brew")
	if err != nil {
		return false
	}

	brewPrefixBytes, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}

	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBytes)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(fingBinary, brewBinPrefix)
}

func isNpmAvailable() bool {
	_, err := exec.LookPath("npm")
	if err != nil {
		return false
	}
	return true
}

func isUnderNpm() bool {
	npmExe, err := exec.LookPath("npm")
	if err != nil {
		return false
	}

	npmOutput, err := exec.Command(npmExe, "list", "-g", npmPackage).Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(npmOutput), npmPackage)
}

func CheckForUpdate(ctx context.Context, version string) (*Release, error) {
	cfg, _ := config.ReadFingConfig()

	if time.Since(cfg.LastCheckedAt).Hours() > 12 {
		release, err := getLatestVersion(ctx)
		if err != nil {
			return nil, err
		}

		cfg.LastCheckedAt = time.Now()
		if err := config.WriteFingConfig(cfg); err != nil {
			return nil, err
		}

		if isNewerVersion(cli.Version, release.Version) {
			return release, nil
		}
	}

	return nil, nil
}

func ShouldCheckUpdate() bool {
	if cli.Version == "dev2" {
		return false
	}
	if os.Getenv("CI") != "" {
		return false
	}
	if os.Getenv("BUILD_NUMBER") != "" {
		return false
	}
	if os.Getenv("RUN_ID") != "" {
		return false
	}
	if os.Getenv("CODESPACES") != "" {
		return false
	}

	return true
}

func isNewerVersion(a, b string) bool {
	av, err := semver.ParseTolerant(a)
	if err != nil {
		fmt.Printf("error parsing version number '%s': %s\n", a, err)
		return false
	}

	bv, err := semver.ParseTolerant(b)
	if err != nil {
		fmt.Printf("error parsing version number '%s': %s\n", a, err)
		return false
	}

	return bv.GT(av)
}
