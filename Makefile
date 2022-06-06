GO_PACKAGES=. ./airer ./auth ./cli ./container ./cygnus ./find ./gcloud ./githelpers ./localexec ./localfile ./sanitize ./validator ./version ./winservice
GO_MODULE_NAME=github.com/mcdonaldseanp/charlie
GO_BIN_NAME=charlie

# Make the build dir, and remove any go bins already there
setup:
	mkdir -p output/
	rm -rf output/$(GO_BIN_NAME)

# Actually build the thing
build: setup
	go mod tidy
	go build -o output/ $(GO_MODULE_NAME)

install:
	go mod tidy
	go install $(GO_MODULE_NAME)

# Build it before publishing to make sure this publication won't be broken
#
# This also ensures that the charlie command is available for the version
# command
publish: install
ifndef NEW_VERSION
	echo "Cannot publish, no tag provided. Set NEW_VERSION to new tag"
else
	charlie update version ./version/version.go $(NEW_VERSION)
	charlie new commit --message "(release) Update to new version $(NEW_VERSION)"
	git tag -a $(NEW_VERSION) -m "Version $(NEW_VERSION)";
	git push
	git push --tags
endif

format:
	go fmt $(REGULATOR_GO_PACKAGES)