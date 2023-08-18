package main

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/updater"
)

func main() {
	fmt.Println("updater running")

	up := updater.NewUpdater(new(localDemo), new(remoteDemo))

	list, err := up.GenerateUpgrades()
	if err != nil {
		panic(err)
	}

	err = up.Exec(list)
	if err != nil {
		panic(err)
	}

	fmt.Println("done")

}

var _ updater.ILocalScripts = new(localDemo)

type localDemo struct {
}

func (p *localDemo) LocalCurrentVersion() (*semver.Version, error) {
	v := "v0.1.2"
	return semver.NewVersion(v)
}

var _ updater.IRemoteScripts = new(remoteDemo)

type remoteDemo struct {
}

// RemoteLatestVersion Get the latest version from the remote server.
func (p *remoteDemo) RemoteLatestVersion() (*semver.Version, error) {
	v := "v1.0.0"
	return semver.NewVersion(v)
}

// RemoteGetUpgradeScripts Fetch the upgrade script from the remote server.
func (p *remoteDemo) RemoteGetUpgradeScripts(v *semver.Version) (updater.IUpgradeScripts, error) {

	return &upgradeScript{
		v: v,
	}, nil

}

// RemoteVersions Retrieve all versions that meet the criteria from the remote server.
func (p *remoteDemo) RemoteVersions(constraints ...*semver.Constraints) ([]*semver.Version, error) {
	l := []string{
		"v1.0.0",
		"v0.3.2",
		"v0.3.1",
		"v0.3.0",
		"v0.2.2",
		"v0.2.1",
		"v0.2.0",
		"v0.1.3",
		"v0.1.2",
	}

	res := make([]*semver.Version, 0)
	for _, v := range l {
		vv, err := semver.NewVersion(v)
		if err != nil {
			continue
		}
		res = append(res, vv)
	}

	return res, nil
}

func versionConstraint() map[string]*upgradeScriptInfo {

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

var _ updater.IUpgradeScripts = new(upgradeScript)

type upgradeScript struct {
	v *semver.Version
}

func (p *upgradeScript) UpgradeInfo() (updater.IUpgradeScriptsInfo, error) {

	mp := versionConstraint()

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
