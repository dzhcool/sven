#!/bin/bash
export PATH=.:/sbin:/usr/sbin:/usr/local/sbin:/usr/local/bin:/bin:/usr/bin:/usr/local/bin

USER=${USER/-/_}
PRJ_ROOT=$(cd `dirname $0`; pwd)
PRJ_NAME=`basename ${PRJ_ROOT}`

# 编译节点注入变量
buildTime=`date +'%Y-%m-%d %H:%M:%S'`
buildVersion=
buildGoVersion=`go version`
LDFlags=" \
    -X 'github.com/dzhcool/sven/buildinfo.build_time=${buildTime}' \
    -X 'github.com/dzhcool/sven/buildinfo.build_version=' \
    -X 'github.com/dzhcool/sven/buildinfo.build_go_version=${buildGoVersion}' \
"
if [ -f "${PRJ_ROOT}/go.mod" ]; then
    go mod tidy
fi

go build -ldflags "$LDFlags" .