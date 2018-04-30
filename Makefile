EXEC=ticker
PLUGIN=${EXEC}.so

COLLECTD_SRC?=/usr/include/collectd/core
BUILD_FLAGS=CGO_CPPFLAGS="-I${COLLECTD_SRC}/daemon -I${COLLECTD_SRC}"

PLUGIN_DIR=/usr/lib/collectd

all:	exec plugin

exec:
	${BUILD_FLAGS} go build -o ${EXEC} main.go

plugin:
	${BUILD_FLAGS} go build -buildmode=c-shared -o ${PLUGIN} main.go

install:
	install ${PLUGIN} ${PLUGIN_DIR}/

clean:
	rm -f ${EXEC} ${PLUGIN}
