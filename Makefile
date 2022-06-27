GO_PACKAGES=. ./airer ./auth ./container ./find ./githelpers ./kubernetes ./localexec ./localfile ./remotedata ./replacers ./sanitize ./version ./winservice
GO_MODULE_NAME=github.com/mcdonaldseanp/charlie
GO_BIN_NAME=charlie
RELEASE_ARTIFACTS=./kubernetes/kind_config.yaml

# Make the build dir, and remove any go bins already there
setup:
	mkdir -p output/
	cd output && \
	rm -f $(GO_BIN_NAME) && \
	for ATFC in $(RELEASE_ARTIFACTS); do \
		rm -f $$(basename $$ATFC); \
	done

# Actually build the thing
build: setup
	go mod tidy
	go build -o output/ $(GO_MODULE_NAME)
	for ATFC in $(RELEASE_ARTIFACTS); do \
		cp $$ATFC output/; \
	done

install:
	go mod tidy
	go install $(GO_MODULE_NAME)

# Build it before publishing to make sure this publication won't be broken
#
# This also ensures that the charlie command is available for the version
# command
#
# If NEW_VERSION is set by the user, it will set the new charlie version
# to that value. Otherwise charlie will bump the Z version
publish: install format
	NEW_VERSION=$$(charlie update version ./version/version.go --version="$(NEW_VERSION)") && \
	echo "Tagging and publishing new version $$NEW_VERSION" && \
	git add --all && \
	git commit -m "(release) Update to new version $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "Version $$NEW_VERSION"
	git push
	git push --tags

format:
	go fmt $(GO_PACKAGES)