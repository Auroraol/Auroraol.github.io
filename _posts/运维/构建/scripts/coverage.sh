#!/bin/bash
PKG_LIST=$(go list ./... | grep -v /vendor/)
for package in ${PKG_LIST}; do
    go test -gcflags="all=-N -l" -covermode=count -coverprofile "scripts/cover/${package##*/}.cov" "$package";
done
echo "mode: count" > scripts/cover/coverage.txt
tail -q -n +2 scripts/cover/*.cov >> scripts/cover/coverage.txt
go tool cover -func=scripts/cover/coverage.txt
