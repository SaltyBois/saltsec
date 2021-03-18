MAINDIR = src
FRONTDIR = src/saltsecui

.PHONY: node_modules test tidy

dev: node_modules
	cd $(MAINDIR) && go run main.go &
	cd $(FRONTDIR) && npm run serve

install: node_modules

node_modules: $(FRONTDIR)/package.json
	cd $(FRONTDIR)
	npm install

test:
	cd $(MAINDIR) && go test -v ./...

tidy:
	cd $(MAINDIR) && go fmt

# TODO(Jovan): Build npm
build:
	go build -v
