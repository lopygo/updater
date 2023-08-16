package demo

import (
	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/updater"
)

var _ updater.IUpgradeScriptsInfo = new(upgradeScriptInfoDemo)

type upgradeScriptInfoDemo struct {
	Version    string `json:"version"`
	Constraint string `json:"constraint"`
}

// UpgradeInfoCurrent
func (p *upgradeScriptInfoDemo) UpgradeInfoCurrent() (*semver.Version, error) {
	return semver.NewVersion(p.Version)
}

// UpgradeInfoConstraint
func (p *upgradeScriptInfoDemo) UpgradeInfoConstraint() (*semver.Constraints, error) {
	return semver.NewConstraint(p.Constraint)
}
