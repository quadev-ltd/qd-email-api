### Linting and formatting
To enable linting and formatting on each commit you need to install the following dependencies:
```
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
```
And run `git config core.hooksPath githooks/`

# TODOs
