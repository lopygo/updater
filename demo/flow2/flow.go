package main

import (
	"log"

	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/adapters/demo"
	"github.com/lopygo/updater/demo/common"
	"github.com/lopygo/updater/updater"
)

func main() {
	local, err := demo.NewLocal(localHandle)
	if err != nil {
		log.Println("init local error: ", err)
		panic(err)
	}

	remote, err := demo.NewRemote(latestHandle, versionsHandle, getScriptsHandle)
	if err != nil {
		log.Println("init remote error: ", err)
		panic(err)
	}

	updater := updater.NewUpdater(local, remote)

	up, err := updater.GenerateUpgrades()
	if err != nil {
		log.Println("generate upgrades error: ", err)
		panic(err)
	}

	log.Println("Display update steps")
	for k, v := range up {
		i, err := v.Script.UpgradeInfo()
		if err != nil {
			panic(err)
		}

		ver, err := i.UpgradeInfoCurrent()
		if err != nil {
			panic(err)
		}

		log.Printf("  -\t%d.\t from [%s] to [%s]", k+1, v.From, ver)

	}

	err = updater.Exec(up)
	if err != nil {
		log.Println("upgrade error: ", err)
		return
	}
	log.Println("upgrade done")
}

var localHandle demo.LocalGetVersionHandle = func() (string, error) {
	return common.VersionRand(), nil
}

var latestHandle demo.RemoteLatestHandle = func() (latest string, err error) {

	return common.VersionRand(), nil
}

var versionsHandle demo.RemoteVersionsHandle = func(constraints ...*semver.Constraints) (versions []string, err error) {
	l := common.VersionAll()

	return l, nil
}

var getScriptsHandle demo.RemoteGetScriptsHandle = func(v *semver.Version) (scripts updater.IUpgradeScripts, err error) {
	scripts = common.NewDemoScript(v)

	return
}
