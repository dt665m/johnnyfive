echo "Lanching BUILD STEP .drone.sh"
go version

LDFLAGS="-w -s -X main.Version=$VERSION -X main.Build=$DRONE_BUILD_NUMBER"

# build multiplatform
TARGET=artifacts/bin
echo "Building $DRONE_REPO_NAME Darwin_amd64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o $TARGET/$DRONE_REPO_NAME.darwin.amd64 github.com/uniontsai/johnnyfive/
echo "Building $DRONE_REPO_NAME Windows_amd64"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o $TARGET/$DRONE_REPO_NAME.windows.amd64.exe github.com/uniontsai/johnnyfive
echo "Building $DRONE_REPO_NAME Linux64_amd64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o $TARGET/$DRONE_REPO_NAME.linux.amd64 github.com/uniontsai/johnnyfive
