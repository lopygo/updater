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
	log.Println("get local version:", localVer)

	latestVer, err := p.remoteScript.RemoteLatestVersion()
	if err != nil {
		return
	}
	log.Println("get latest version from remote:", latestVer)

	if !latestVer.GreaterThan(localVer) {
		err = fmt.Errorf("the local version is already up to date")
		return
	}

	log.Printf("preparing to upgrade from %s to %s.", localVer, latestVer)
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

	log.Println("################################ compare versions ################################")
	log.Printf("compare versions from [%s] to [%s]", fromVersion.String(), toVersion.String())
	if !fromVersion.LessThan(toVersion) {
		err = fmt.Errorf("local version[%s] not less than target[%s] ", fromVersion, toVersion)
		return
	}

	log.Printf("download the scripts of version[%s] from remote if not exists", toVersion.String())
	upgradeScripts, err := p.remoteScript.RemoteGetUpgradeScripts(toVersion)
	if err != nil {
		return
	}

	log.Println("get scripts info")
	scriptInfo, err := upgradeScripts.UpgradeInfo()
	if err != nil {
		return
	}

	log.Println("get version of scripts")
	infoVersion, err := scriptInfo.UpgradeInfoCurrent()
	if err != nil {
		return
	}

	log.Println("match versions of scripts and binary")
	if !infoVersion.Equal(toVersion) {
		err = fmt.Errorf("the target version number for the upgrade should match the script version number")
		return
	}

	log.Println("get constraint from this script")
	conConstraint, err := scriptInfo.UpgradeInfoConstraint()
	if err != nil {
		return
	}

	log.Printf(`the constraint is "%s" of "%s"`, conConstraint.String(), toVersion.String())

	tttt := &UpgradeArgs{
		Script: upgradeScripts,
		From:   fromVersion,
	}

	*targets = append([]*UpgradeArgs{
		tttt,
	}, *targets...)

	// this only looking for release versions
	log.Println(`constraint check self`, conConstraint.String(), fromVersion.String())
	if conConstraint.Check(fromVersion) {
		return
	}
	log.Println(`constraint check self res`, conConstraint.Check(fromVersion))

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

	sort.Sort(semver.Collection(listChecked))

	// if it is a prerelease version
	if fromVersion.GreaterThan(listChecked[0]) {
		tttt.From = fromVersion
		return nil
	}

	tttt.From = listChecked[0]
	return p.compareVersionFromAndTo(targets, fromVersion, listChecked[0])

}

type UpgradeArgs struct {
	From   *semver.Version
	Script IUpgradeScripts
}
