package main

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/demo/common"
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
	v := common.VersionRand()
	return semver.NewVersion(v)
}

var _ updater.IRemoteScripts = new(remoteDemo)

type remoteDemo struct {
}

// RemoteLatestVersion Get the latest version from the remote server.
func (p *remoteDemo) RemoteLatestVersion() (*semver.Version, error) {
	v := common.VersionRand()
	return semver.NewVersion(v)
}

// RemoteGetUpgradeScripts Fetch the upgrade script from the remote server.
func (p *remoteDemo) RemoteGetUpgradeScripts(v *semver.Version) (updater.IUpgradeScripts, error) {

	return common.NewDemoScript(v), nil

}

// RemoteVersions Retrieve all versions that meet the criteria from the remote server.
func (p *remoteDemo) RemoteVersions(constraints ...*semver.Constraints) ([]*semver.Version, error) {
	l := common.VersionAll()

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
