# !!!MAKE SURE YOUR GOPATH ENVIRONMENT VARIABLE IS SET FIRST!!!
# Any issues with this file, let me know / make a PR, I haven't tested it completely but it should be close enough.

AGENT=DaaC2Agent
SERVER=DaaC2Server
DIR=out
LDFLAGS=-ldflags "-s -w"
D=Darwin-x64
L=Linux-x64
W=Windows-x64

# Make Directory to store executables
$(shell mkdir -p ${DIR})

# Change default to just make for the host OS and add MAKE ALL to do this
default: server-linux server-darwin agent-windows agent-linux

all: default

# Compile Darwin binaries
darwin: server-darwin agent-darwin

# Compile Linux Binaries
linux: server-linux agent-linux

windows: agent-windows

# Compile Agent - Linux x64
agent-linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${AGENT}-${L} cmd/agent/main.go

# Compile Agent - Windows x64     REPLACE LDFLAGS + WINAGENTLDFLAGS for actual release!!
agent-windows:
	export GOOS=windows GOARCH=amd64;go build ${WINAGENTLDFLAGS} ${LDFLAGS} -o ${DIR}/${AGENT}-${W}.exe cmd/agent/main.go

# Compile Agent - MacOS
agent-darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${AGENT}-${D} cmd/agent/main.go

# Compile Server - MacOS
server-darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${SERVER}-${D} cmd/server/main.go

# Compile Server - Linux x64
server-linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${SERVER}-${L} cmd/server/main.go

# Compile Server - Windows x64
server-windows:
	export GOOS=windows;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${SERVER}-${W}.exe cmd/server/main.go

clean:
	rm -rf ${DIR}*