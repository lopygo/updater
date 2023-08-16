package demo

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/updater"
)

type RemoteLatestHandle func() (latest string, err error)

type RemoteVersionsHandle func(constraints ...*semver.Constraints) (versions []string, err error)

type RemoteGetScriptsHandle func(v *semver.Version) (scripts updater.IUpgradeScripts, err error)

var _ updater.IRemoteScripts = new(remoteDemo)

type remoteDemo struct {
	latestFn    RemoteLatestHandle
	versionsFn  RemoteVersionsHandle
	getScriptFn RemoteGetScriptsHandle
}

// RemoteLatestVersion Get the latest version from the remote server.
func (p *remoteDemo) RemoteLatestVersion() (v *semver.Version, err error) {

	if p.latestFn == nil {
		err = fmt.Errorf("handle can not set")
		return
	}

	la, err := p.latestFn()
	if err != nil {
		return
	}

	return semver.NewVersion(la)
}

// RemoteGetUpgradeScripts Fetch the upgrade script from the remote server.
func (p *remoteDemo) RemoteGetUpgradeScripts(v *semver.Version) (scripts updater.IUpgradeScripts, err error) {

	if p.getScriptFn == nil {
		err = fmt.Errorf("handle can not set")
		return
	}

	return p.getScriptFn(v)
}

// RemoteVersions Retrieve all versions that meet the criteria from the remote server.
func (p *remoteDemo) RemoteVersions(constraints ...*semver.Constraints) (versions []*semver.Version, err error) {

	if p.versionsFn == nil {
		err = fmt.Errorf("handle can not set")
		return
	}

	l, err := p.versionsFn(constraints...)
	if err != nil {
		return
	}

	versions = make([]*semver.Version, 0)
	for _, v := range l {
		vv, err := semver.NewVersion(v)
		if err != nil {
			continue
		}
		versions = append(versions, vv)
	}

	sort.Sort(semver.Collection(versions))

	return
}

func NewRemote(latestFn RemoteLatestHandle,
	versionsFn RemoteVersionsHandle,
	getScriptFn RemoteGetScriptsHandle) (p *remoteDemo, err error) {

	if latestFn == nil {
		err = fmt.Errorf("latestFn can not empty")
		return
	}

	if versionsFn == nil {
		err = fmt.Errorf("versionsFn can not empty")
		return
	}
	if getScriptFn == nil {
		err = fmt.Errorf("getScriptFn can not empty")
		return
	}

	p = &remoteDemo{
		versionsFn:  versionsFn,
		latestFn:    latestFn,
		getScriptFn: getScriptFn,
	}

	return p, nil
}
