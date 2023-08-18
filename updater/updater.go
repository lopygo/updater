package updater

import (
	"fmt"
	"log"
	"sort"

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

func (p *Updater) GenerateUpgrades() (upgrades []*UpgradeArgs, err error) {

	localVer, err := p.localScript.LocalCurrentVersion()
	if err != nil {

		return
	}
	log.Println("get version of local", localVer)

	latestVer, err := p.remoteScript.RemoteLatestVersion()
	if err != nil {
		return
	}
	log.Println("get latest version from remote", latestVer)

	if !latestVer.GreaterThan(localVer) {
		err = fmt.Errorf("the local version is already up to date")
		return
	}

	log.Println("set target version, from current")
	upgrades = []*UpgradeArgs{}

	log.Println("loop checking")
	err = p.compareVersionFromAndTo(&upgrades, localVer, latestVer)

	return
}

func (p *Updater) Exec(upgrades []*UpgradeArgs) error {
	for _, v := range upgrades {
		err := v.Script.UpgradeExec(v.From)
		if err != nil {
			return err
		}

	}

	return nil
}

func (p *Updater) compareVersionFromAndTo(targets *[]*UpgradeArgs, fromVersion, toVersion *semver.Version) (err error) {

	if !fromVersion.LessThan(toVersion) {
		err = fmt.Errorf("local version[%s] not less than target[%s] ", fromVersion, toVersion)
		return
	}

	log.Println("download the version from remote if not exists")
	upgradeScripts, err := p.remoteScript.RemoteGetUpgradeScripts(toVersion)
	if err != nil {

		return
	}

	log.Println("get constraint info from upgrade scripts")
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

	tttt := &UpgradeArgs{
		Script: upgradeScripts,
		From:   fromVersion,
	}
	*targets = append([]*UpgradeArgs{
		tttt,
	}, *targets...)

	if conConstraint.Check(fromVersion) {
		return
	}

	log.Println("list versions from remote")
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

		return fmt.Errorf("no version that meet the criteria were found")
	}
	tttt.From = listChecked[0]

	sort.Sort(semver.Collection(listChecked))

	return p.compareVersionFromAndTo(targets, fromVersion, listChecked[0])

}

type UpgradeArgs struct {
	From   *semver.Version
	Script IUpgradeScripts
}
