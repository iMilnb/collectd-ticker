BINARY=ticker

all:
	COLLECTD_SRC="/usr/include/collectd/core"			\
	CGO_CPPFLAGS="-I${COLLECTD_SRC}/daemon -I${COLLECTD_SRC}"	\
	go build -o ${BINARY} main.go
