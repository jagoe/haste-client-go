version=1.0.0-beta.1
commit=$(shell git rev-parse --short HEAD)
built_at=$(shell date +%FT%T%z)
built_by=Jakob Göbel

build:
	@go build -ldflags "-X 'github.com/jagoe/haste-client-go/cmd.version=${version}' \
		-X 'github.com/jagoe/haste-client-go/cmd.commit=${commit}' \
		-X 'github.com/jagoe/haste-client-go/cmd.builtBy=${built_by}' \
		-X 'github.com/jagoe/haste-client-go/cmd.builtAt=${built_at}'" \
	-o bin/haste ./main.go

test:
	@$(call test) | $(call color_test_output)

coverage:
	@$(call test) -cover -json | $(call jq-coverage) | $(call color_test_output)

test-summary:
	@$(call test) -json \
	| $(call jq-summary) \
	| $(call color_test_output)

### HELPERS ###

define test
	@go test -v ./...
endef

define color_test_output
	sed "/^PASS/s//$(shell printf "\033[32mPASS\033[0m")/" \
	| sed "/^ok/s//$(shell printf "\033[32mok\033[0m")/" \
	| sed "/^FAIL/s//$(shell printf "\033[31mFAIL\033[0m")/" \
	| sed "/^?/s//$(shell printf "\033[33m?\033[0m")/" \
	| sed "/.*\.go:[[:digit:]]*:.*/s//$(shell printf "\033[31m&\033[0m")/"
endef

define jq-summary
	jq -r --slurp '.[] | select(.Action == "output") | select(.Output | test("^(ok|\\?|FAIL)\\s*\\t")) | .Output[:-1]' \
	| sort -k2,2
endef

define jq-coverage
	jq -r --slurp '.[] | select(.Action == "output") | select(.Output | test("(.+coverage: \\d+\\.\\d+% of statements)|(no test files)")) | .Output[:-1]' \
	| sort -k2,2
endef
