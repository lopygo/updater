package updater

import (
	"fmt"
	"sort"
	"time"

	"github.com/Masterminds/semver/v3"
)

type Updater struct {
	localScript  ILocalScripts
	remoteScript IRemoteScripts
}

func NewUpdater(local ILocalScripts, remote IRemoteScripts) *Updater {

	return &Updater{
		localScript:  local,
		remoteScript: remote,
	}
}

func (p *Updater) ResolveList() (targets []IUpgradeScripts, err error) {

	fmt.Println("updater running")

	localVer, err := p.localScript.LocalCurrentVersion()
	if err != nil {

		return
	}
	fmt.Println("get version of local", localVer)

	latestVer, err := p.remoteScript.RemoteLatestVersion()
	if err != nil {
		return
	}
	fmt.Println("get latest version from remote", latestVer)

	fmt.Println("set target version, from current")
	targets = []IUpgradeScripts{}

	fmt.Println("loop")
	err = p.compareVersionFromAndTo(&targets, localVer, latestVer)

	return
}

func (p *Updater) Exec(targets []IUpgradeScripts) error {
	for _, v := range targets {
		err := v.UpgradeExec()
		if err != nil {
			return err
		}

	}

	return nil
}

func (p *Updater) compareVersionFromAndTo(targets *[]IUpgradeScripts, fromVersion, toVersion *semver.Version) (err error) {
	time.Sleep(time.Second)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("compare from and to", fromVersion, toVersion)

	if !fromVersion.LessThan(toVersion) {
		err = fmt.Errorf("local version[%s] less than target[%s] ", fromVersion, toVersion)
		return
	}

	fmt.Println("download the version from remote if not exists")
	upgradeScripts, err := p.remoteScript.RemoteGetUpgradeScripts(toVersion)
	if err != nil {

		return
	}

	fmt.Println("get constraint info from upgrade scripts")
	scriptInfo, err := upgradeScripts.UpgradeInfo()
	if err != nil {
		return
	}

	infoVersion, err := scriptInfo.UpgradeInfoCurrent()
	if err != nil {
		return
	}

	if !infoVersion.Equal(toVersion) {
		err = fmt.Errorf("the target version number for the upgrade should match the script version number")
		return
	}

	conConstraint, err := scriptInfo.UpgradeInfoConstraint()

	*targets = append([]IUpgradeScripts{upgradeScripts}, *targets...)

	if conConstraint.Check(fromVersion) {

		fmt.Println("obtaining the minimum version that satisfies constraints. set target version")
		return
	}

	fmt.Println("list versions from remote")
	list, err := p.remoteScript.RemoteVersions(conConstraint)
	if err != nil {
		return err
	}

	listChecked := []*semver.Version{}
	for _, v := range list {

		if !conConstraint.Check(v) {
			continue
		}

		if !toVersion.GreaterThan(v) {
			continue
		}

		listChecked = append(listChecked, v)
	}

	if len(listChecked) == 0 {

		return fmt.Errorf("no version man zu condition")
	}

	sort.Sort(semver.Collection(listChecked))

	return p.compareVersionFromAndTo(targets, fromVersion, listChecked[0])

}
