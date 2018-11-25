.PHONY: build
build: tools gen
	go build -o bin/fm5 cmd/fm5/main.go

.PHONY: gen
gen:
	go generate ./...

clean:
	rm bin/* e2e/fm5 e2e/*_fm.go

tools: bin/generr bin/richgo

bin/generr: vendor/github.com/akito0107/generr/cmd/generr
	go build -o bin/generr ./vendor/github.com/akito0107/generr/cmd/generr

bin/richgo: vendor/github.com/kyoh86/richgo
	go build -o bin/richgo ./vendor/github.com/kyoh86/richgo

vendor/github.com/kyoh86/richgo: vendor
vendor/github.com/akito0107/generr/cmd/generr: vendor

vendor:
	dep ensure

.PHONY: test
test: test/small test/e2e

.PHONY: test/small
test/small: tools
	bin/richgo test -v -coverprofile=out -count=1

.PHONY: test/e2e
test/e2e: e2e/fm5 tools
	cd e2e; ../bin/richgo test -count=1 -v .

e2e/fm5:
	go build -o e2e/fm5 cmd/fm5/main.go
