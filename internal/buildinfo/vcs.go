package buildinfo

import (
	"runtime/debug"
	"sync"
)

var vcsInfo = sync.OnceValue(func() string {

	binfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	var (
		revision string
		modified bool
	)
	for _, setting := range binfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		default:
		}
	}
	if revision == "" {
		return "unknown"
	}
	if len(revision) > 6 {
		revision = revision[:6]
	}
	if modified {
		revision += "+dirty"
	}
	return revision
})

func VCSInfo() string {
	return vcsInfo()
}
