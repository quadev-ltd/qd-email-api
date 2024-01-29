## Hooks: Linting and formatting
To enable linting and formatting on each commit we use the following dependencies:
```
golang.org/x/tools/cmd/goimports@latest
golang.org/x/lint/golint@latest
```
To activate commit hooks use the following command:
```git config core.hooksPath githooks/```
And make `githooks/pre-commit` executable.
To avoid running hooks do `git commit --no-verify`

## GRPC
To generate the grpc code:
- Follow the steps in https://buf.build/docs/installation to install buf.
- In the root of the repository, run `git submodule update --init --recursive`.
- Then, in `/pb/`, run `buf generate` to generate the protobuf files.  
> note: Flags `-v --debug` will provide more details on the execution.


# TODOs
