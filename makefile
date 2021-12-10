# Build the project
default:
	@echo "--> Running build"
	@sh -c "$(CURDIR)/scripts/build.sh"

dev: assets
	@echo "--> Running build"
	@DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

fmt:
	@cd $(CURDIR)
	@go fmt $$(go list ./... | grep -v /vendor/)

test: tools dev
	@echo "--> Running go test"
	go list ./... | grep -v -E '^github.com/teatak/pipe/(vendor|cmd/pipe/vendor)' | xargs -n1 go test

.PHONY: default fmt
