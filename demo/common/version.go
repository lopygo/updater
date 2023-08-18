package common

import "math/rand"

func VersionAll() []string {
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

	return l
}

func VersionRand() string {
	l := VersionAll()
	n := rand.Intn(len(l))

	return l[n]
}
