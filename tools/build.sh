shell_path=$(cd `dirname $0`; pwd)
CGO_ENABLED=1 GOOS=linux cd ${shell_path}/../ && go build -ldflags="-s -w" -o contrib/bin/godocs