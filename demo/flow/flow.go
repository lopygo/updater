package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/Masterminds/semver/v3"
)

func main() {
	fmt.Println("updater running")

	localVer, err := localVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println("get version of local", localVer)

	latestVer, err := remoteLatest()
	if err != nil {
		panic(err)
	}
	fmt.Println("get latest version from remote", latestVer)

	fmt.Println("set target version, from current")
	targets := []*semver.Version{}

	fmt.Println("loop")
	err = compareVersionFromAndTo(&targets, localVer, latestVer)
	if err != nil {
		panic(err)
	}
	// loop

	fmt.Println()
	fmt.Println()
	fmt.Println("ready exec update")
	fmt.Println(targets)
	for _, v := range targets {
		fmt.Println("exec upgrade to ", v)

	}

}

func compareVersionFromAndTo(targets *[]*semver.Version, fromVersion, toVersion *semver.Version) (err error) {
	time.Sleep(time.Second)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("compare from and to", fromVersion, toVersion)

	if !fromVersion.LessThan(toVersion) {
		err = fmt.Errorf("local version[%s] less than target[%s] ", fromVersion, toVersion)
		return
	}

	fmt.Println("if local upgrade scripts exists")

	fmt.Println("download the version from remote if not exists")

	fmt.Println("get constraint info from upgrade scripts")

	versionConstraintMp := versionConstraint()
	conVer, ok := versionConstraintMp[toVersion.String()]
	if !ok {
		return fmt.Errorf("can not get constraint from upgrade scripts")
	}
	conConstraint, err := semver.NewConstraint(conVer.Constraint)
	if err != nil {
		// Handle constraint not being parsable.
		return fmt.Errorf("can not get constraint from scripts downloaded")
	}

	*targets = append([]*semver.Version{toVersion}, *targets...)

	if !conConstraint.Check(fromVersion) {
		fmt.Println("list versions from remote")
		list := remoteVersions()

		listChecked := []*semver.Version{}
		for _, v := range list {
			vv, err := semver.NewVersion(v)
			if err != nil {
				continue
			}
			if !conConstraint.Check(vv) {
				continue
			}

			if !toVersion.GreaterThan(vv) {
				continue
			}

			listChecked = append(listChecked, vv)
		}

		if len(listChecked) == 0 {

			return fmt.Errorf("no version man zu condition")
		}

		sort.Sort(semver.Collection(listChecked))

		// *targets = append(*targets, listChecked[0])
		// *targets = append([]*semver.Version{listChecked[0]}, *targets...)
		return compareVersionFromAndTo(targets, fromVersion, listChecked[0])
	}

	fmt.Println("obtaining the minimum version that satisfies constraints. set target version")
	// *targets = append(*targets, toVersion)
	return
}

func remoteVersions() []string {
	return []string{
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
}

func versionConstraint() map[string]*Info {

	mp := map[string]*Info{
		"1.0.0": {
			Constraint: ">= v0.3.2",
		},
		"0.3.3": {
			Constraint: ">= v0.3.0",
		},
		"0.3.2": {
			Constraint: ">= v0.3.0",
		},
		"0.3.1": {
			Constraint: ">= v0.3.0",
		},
		"0.3.0": {
			Constraint: "= v0.2.2",
		},
		"0.2.2": {
			Constraint: ">= v0.2.0",
		},
		"0.2.1": {
			Constraint: ">= v0.2.0",
		},
		"0.2.0": {
			Constraint: "v0.1.3",
		},
		"0.1.3": {
			Constraint: "< v0.1.3",
		},
		"0.1.2": {
			Constraint: "< v0.1.2",
		},
	}

	for k, _ := range mp {
		mp[k].Current = k
	}
	return mp
}

func remoteLatest() (vv *semver.Version, err error) {

	v := "v1.0.0"
	return semver.NewVersion(v)
}

func localVersion() (vv *semver.Version, err error) {

	v := "v0.1.2"
	return semver.NewVersion(v)

}

type Info struct {
	Current    string `json:"current"`
	Constraint string `json:"constraint"`
}
