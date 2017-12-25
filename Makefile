BINARY=ticker

COLLECTD_SRC=/usr/include/collectd/core
BUILD_FLAGS=CGO_CPPFLAGS="-I${COLLECTD_SRC}/daemon -I${COLLECTD_SRC}"

exec:
	${BUILD_FLAGS} go build -o ${BINARY} main.go

plugin:
	${BUILD_FLAGS} go build -buildmode=c-shared -o ${BINARY}.so main.go
