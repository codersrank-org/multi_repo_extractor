generated: build

test:
	@go test ./...

build: test
	export CGO_ENABLED=0 export GOOS=linux && go build -a -tags netgo -ldflags '-w' -o multi_repo_extractor_linux
	export CGO_ENABLED=0 export GOOS=darwin && go build -a -tags netgo -ldflags '-w' -o multi_repo_extractor_osx
	export CGO_ENABLED=0 export GOOS=windows && go build -a -tags netgo -ldflags '-w' -o multi_repo_extractor_windows
	export GOOS=$GOOS_OLD