package buildinfo

import (
	"runtime/debug"
	"sync"

	"go.opentelemetry.io/otel/attribute"
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
	if len(revision) > 8 {
		revision = revision[:8]
	}
	if modified {
		revision += "+dirty"
	}
	return revision
})

func VCSInfo() string {
	return vcsInfo()
}

func VCSAttribute() attribute.KeyValue {
	return attribute.String("vcs.revision", VCSInfo())
}
