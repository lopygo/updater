package updater

import "github.com/Masterminds/semver/v3"

type IRemoteScripts interface {
	// RemoteLatestVersion Get the latest version from the remote server.
	RemoteLatestVersion() (*semver.Version, error)

	// RemoteGetUpgradeScripts Fetch the upgrade script from the remote server.
	RemoteGetUpgradeScripts(*semver.Version) (IUpgradeScripts, error)

	// RemoteVersions Retrieve all versions that meet the criteria from the remote server.
	RemoteVersions(constraints ...*semver.Constraints) ([]*semver.Version, error)
}
