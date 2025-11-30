#!/bin/bash
set -ex

VERSION=$(git describe --tags 2>/dev/null)
COMMIT=$(git rev-parse --short HEAD)
TIME=$(date +%FT%T)
if [ -z $VERSION ]; then
    VERSION=$COMMIT
fi

VerPkg="gitlab.xiaoduoai.com/golib/xd_sdk/xdspec"
BASE=$(pwd)

go_tags="nomsgpack"

export CGO_ENABLED=0
go build -tags "${go_tags}" \
    -ldflags "-X ${VerPkg}.Version=$VERSION \
    -X ${VerPkg}.GitCommit=$COMMIT \
    -X ${VerPkg}.BuildTime=${TIME}" \
    -o ./dist/goods_center
cd component/cid/center_worker && go build  -o ${BASE}/dist/center_worker && cd ${BASE}
cd component/cid/sync_worker && go build  -o ${BASE}/dist/sync_worker && cd ${BASE}
cd component/cid/filter_worker && go build  -o ${BASE}/dist/filter_worker && cd ${BASE}
cd component/notify/notify_worker && go build  -o ${BASE}/dist/notify_worker && cd ${BASE}
cd component/notify-filter && go build  -o ${BASE}/dist/notify-filter && cd ${BASE}
cd component/repair/repair_worker && go build  -o ${BASE}/dist/repair_worker && cd ${BASE}
cd component/platform_implement/taobao && go build  -o ${BASE}/dist/tmc_worker && cd ${BASE}
cd component/platform_implement/general && go build  -o ${BASE}/dist/gmc_worker && cd ${BASE}