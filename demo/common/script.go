package common

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/updater"
)

var _ updater.IUpgradeScripts = new(upgradeScript)

type upgradeScript struct {
	v *semver.Version
}

func NewDemoScript(v *semver.Version) *upgradeScript {
	return &upgradeScript{
		v: v,
	}
}

func (p *upgradeScript) UpgradeInfo() (updater.IUpgradeScriptsInfo, error) {

	mp := versionConstraints()

	vv, ok := mp[p.v.String()]

	if !ok {
		return nil, fmt.Errorf("no this version script")
	}

	return vv, nil
}

func (p *upgradeScript) UpgradeExec(from *semver.Version) error {

	fmt.Println("exec upgrade: ", from, "to", p.v)
	return nil
}
