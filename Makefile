##################
###  Tests  ###
##################

TEST_PACKAGES ?= ./...
UNIT_TEST_TAGS = norace
TEST_RACE_TAGS = ""

tests: ARGS=-timeout=10m -tags='$(UNIT_TEST_TAGS)' -p=4
tests:
ifneq (,$(shell which tparse 2>/dev/null))
	@echo "--> Running tests"
	@go test -mod=readonly -json $(ARGS) $(TEST_PACKAGES) | tparse
else
	@echo "--> Running tests"
	@go test -mod=readonly $(ARGS) $(TEST_PACKAGES)
endif