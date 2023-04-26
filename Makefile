.PHONY: build

SERVICE_NAME=wtpl
IMAGE_VERSION=v1.0.0
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD)
BUILD_BRANCH    := $(shell git symbolic-ref --short HEAD)
OS_ARCH         := `go env GOOS`/`go env GOARCH`
GO_VERSION      := `go env GOVERSION`

build:
	go mod tidy
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags  \
	"														 \
	-X 'wtpl/cmd.Version=${IMAGE_VERSION}'			     \
	-X 'wtpl/cmd.BuildTime=${BUILD_TIME}'				 \
	-X 'wtpl/cmd.CommitID=${COMMIT_SHA1}'       			 \
	-X 'wtpl/cmd.BuildBranch=${BUILD_BRANCH}'       	     \
	-X 'wtpl/cmd.OsArch=${OS_ARCH}'       	     \
	-X 'wtpl/cmd.GoVersion=${GO_VERSION}'       	     \
	" 														 \
	-o ${SERVICE_NAME} main.go