bin:
	rm -rf bin
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/dev-tools-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/dev-tools-arm64 .

build: bin
	rm -rf dist
	mkdir -p dist/bin
	cp -r bin/* dist/bin/
	cp run.sh dist/
	chmod +x dist/run.sh
	cp info.plist dist/
	pushd dist && zip -r ../DevTools.alfredworkflow . && popd

.PHONY: bin build