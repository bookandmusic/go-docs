shell_path=$(cd `dirname $0`; pwd)
gofmt -l -w ./
goimports-reviser -rm-unused -set-alias -format ./...