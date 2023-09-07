shell_path=$(cd `dirname $0`; pwd)
GOOS=linux GOARCH=amd64 cd ${shell_path}/../ && go build -ldflags '-w -s' -o contrib/bin/godocs-amd64
GOOS=linux GOARCH=arm64 go build -ldflags '-w -s' -o contrib/bin/godocs-arm64