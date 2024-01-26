shell_path=$(cd `dirname $0`; pwd)

go get golang.org/x/tools/cmd/goimports

go mod tidy
gofmt -l -w ./
goimports-reviser -rm-unused -set-alias -format ./...