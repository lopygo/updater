package updater

import "github.com/Masterminds/semver/v3"

type ILocalScripts interface {
	// LocalCurrentVersion Get the current version locally
	LocalCurrentVersion() (*semver.Version, error)
}
