go.mod: tools/tools.go
	go mod tidy
	touch go.mod

# must be global install becuse of golines (base-formatter option)
gofumpt:
	command -v gofumpt >/dev/null || go install mvdan.cc/gofumpt@latest

tools: go.mod gofumpt
	GOBIN=$(PWD)/tools/bin go install $$(go list -f '{{join .Imports " "}}' tools/tools.go)

