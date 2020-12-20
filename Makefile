define color_test_output
	sed "/^PASS/s//$(shell printf "\033[32mPASS\033[0m")/" \
	| sed "/^ok/s//$(shell printf "\033[32mok\033[0m")/" \
	| sed "/^FAIL/s//$(shell printf "\033[31mFAIL\033[0m")/" \
	| sed "/^?/s//$(shell printf "\033[33m?\033[0m")/"
endef

define test
	@go test -v ./...
endef

build:
	@go build -o bin/haste ./main.go

test:
	@$(call test) | $(call color_test_output)

cover:
	@$(call test) -cover | $(call color_test_output)

coverage: cover

test-summary:
	@$(call test) -json \
	| jq -r --slurp '.[] | select(.Action == "output") | select(.Output | test("^(ok|\\?|FAIL)\\s*\\t")) | .Output[:-1]' \
	| sort -k2,2 \
	| $(call color_test_output)
