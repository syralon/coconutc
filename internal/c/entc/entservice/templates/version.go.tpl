// @file: version/version.go

package version

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var (
    BuildTime = "unknown"
    Version = "unknown"
)

func Show() {
	fmt.Printf("%s %s %s\n", path.Base(strings.ReplaceAll(os.Args[0], "\\", "/")), Version, BuildTime)
}
