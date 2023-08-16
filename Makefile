

all: 



test: fmt gotest

testp:

	go test ./...

fmt:

	go vet -buildvcs=false `go list ./... | grep -v /tmp | grep -v /temp`

	go fmt `go list ./... | grep -v /tmp | grep -v /temp`

gotest:
	go test -buildvcs=false -v --cover `go list ./... | grep -v /tmp | grep -v /temp`
	

############################################
######## tag version test start#############
############################################
tagChangeFileLatest=changelog/latest.md
tagChangeFileAll:="changelog/CHANGELOG.md"
tagCurrentVersion=$(shell git describe --tags --abbrev=0)

tagRechangelog:
	conventional-changelog -p angular -i "$(tagChangeFileAll)" -r 0 -s

	git tag -d "$(tagCurrentVersion)"

	git add $(tagChangeFileAll)

	git commit -am "generate $(tagChangeFileAll)"
	
	git tag -a "$(tagCurrentVersion)" -m "chore(release): $(tagCurrentVersion)"

# <major|minor|patch>
release-as?=patch
tag:
	mkdir -p `dirname ${tagChangeFileLatest}`
	rm -rf ${tagChangeFileLatest}
	standard-version -i ${tagChangeFileLatest} --header="" --release-as=${release-as} --scripts.posttag="make tagRechangelog"

taga:
	make tag release-as=minor

tagb:
	make tag release-as=major

############################################
######## tag version test end###############
############################################
