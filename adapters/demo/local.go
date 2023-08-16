package demo

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/lopygo/updater/updater"
)

var _ updater.ILocalScripts = new(localDemo)

type LocalGetVersionHandle func() (string, error)

func NewLocal(fn LocalGetVersionHandle) (*localDemo, error) {
	if fn == nil {
		return nil, fmt.Errorf("localHandle can not empty")
	}

	return &localDemo{
		getVersion: fn,
	}, nil
}

type localDemo struct {
	getVersion LocalGetVersionHandle
}

func (p *localDemo) LocalCurrentVersion() (*semver.Version, error) {

	if p.getVersion == nil {
		return nil, fmt.Errorf("handle can not set")
	}
	v, err := p.getVersion()

	if err != nil {
		return nil, err
	}

	return semver.NewVersion(v)
}
