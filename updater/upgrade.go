package updater

import "github.com/Masterminds/semver/v3"

type IUpgradeScripts interface {
	UpgradeInfo() (IUpgradeScriptsInfo, error)

	UpgradeExec() error
}

type IUpgradeScriptsInfo interface {

	// UpgradeInfoCurrent
	UpgradeInfoCurrent() (*semver.Version, error)

	// UpgradeInfoConstraint
	UpgradeInfoConstraint() (*semver.Constraints, error)
}
