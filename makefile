MAINDIR = src
FRONTDIR = src/saltsecui

.PHONY: node_modules test

dev: node_modules
	cd $(MAINDIR) && go run main.go &
	cd $(FRONTDIR) && npm run serve

install: node_modules

node_modules: $(FRONTDIR)/package.json
	cd $(FRONTDIR)
	npm install

test:
	cd $(MAINDIR) && go test -v ./...
