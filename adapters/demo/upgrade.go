package demo

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/updater"
)

type GetUpgradeInfoHandle func(to *semver.Version) (cmdOutput []byte, err error)

type UpgradeHandle func(to *semver.Version) (err error)

func NewUpgrade(to *semver.Version, fnInfo GetUpgradeInfoHandle, fnUpgrade UpgradeHandle) (p *upgradeScriptDemo, err error) {

	if to == nil {
		err = fmt.Errorf("toVersion can not empty")
		return
	}

	if fnInfo == nil {
		err = fmt.Errorf("infoHandle can not empty")
		return
	}
	if fnUpgrade == nil {
		err = fmt.Errorf("upgradeHandle can not empty")
		return
	}

	p = &upgradeScriptDemo{
		to:        to,
		fnInfo:    fnInfo,
		fnUpgrade: fnUpgrade,
	}

	return p, nil
}

var _ updater.IUpgradeScripts = new(upgradeScriptDemo)

type upgradeScriptDemo struct {
	to        *semver.Version
	fnInfo    GetUpgradeInfoHandle
	fnUpgrade UpgradeHandle
}

func (p *upgradeScriptDemo) UpgradeInfo() (updater.IUpgradeScriptsInfo, error) {
	if p.fnInfo == nil {
		return nil, fmt.Errorf("handle can not set")
	}

	buf, err := p.fnInfo(p.to)
	if err != nil {
		return nil, err
	}

	var info upgradeScriptInfoDemo

	err = json.Unmarshal(buf, &info)
	if err != nil {
		return nil, fmt.Errorf("parse upgrade info error: %v", err)
	}

	return &info, nil
}

func (p *upgradeScriptDemo) UpgradeExec() error {

	if p.fnUpgrade == nil {
		return fmt.Errorf("handle can not set")
	}

	return p.fnUpgrade(p.to)
}
