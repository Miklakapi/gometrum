package version

import (
	"runtime/debug"
)

// Set from the release tag by the release build.
var releaseVersion string

func VersionString() string {
	if releaseVersion != "" {
		return releaseVersion
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	if info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "devel"
	}

	return info.Main.Version
}
