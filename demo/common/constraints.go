package common

import (
	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/updater"
)

func versionConstraints() map[string]*upgradeScriptInfo {

	mp := map[string]*upgradeScriptInfo{
		"1.0.0": {
			constraint: ">= v0.3.2",
		},
		"0.3.3": {
			constraint: ">= v0.3.0",
		},
		"0.3.2": {
			constraint: ">= v0.3.0",
		},
		"0.3.1": {
			constraint: ">= v0.3.0",
		},
		"0.3.0": {
			constraint: "= v0.2.2",
		},
		"0.2.2": {
			constraint: ">= v0.2.0",
		},
		"0.2.1": {
			constraint: ">= v0.2.0",
		},
		"0.2.0": {
			constraint: "v0.1.3",
		},
		"0.1.3": {
			constraint: "< v0.1.3",
		},
		"0.1.2": {
			constraint: "< v0.1.2",
		},
	}

	for k, _ := range mp {
		mp[k].v = k
	}
	return mp
}

var _ updater.IUpgradeScriptsInfo = new(upgradeScriptInfo)

type upgradeScriptInfo struct {
	v          string
	constraint string
}

// UpgradeInfoCurrent
func (p *upgradeScriptInfo) UpgradeInfoCurrent() (*semver.Version, error) {
	return semver.NewVersion(p.v)
}

// UpgradeInfoConstraint
func (p *upgradeScriptInfo) UpgradeInfoConstraint() (*semver.Constraints, error) {
	return semver.NewConstraint(p.constraint)
}
