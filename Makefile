.PHONY: all dist dist-local dist-travis start-docker build-server package build-client test travis-init build-container stop-docker clean-docker clean nuke run stop setup-mac cleandb docker-build docker-run

GOPATH ?= $(GOPATH:)
GOFLAGS ?= $(GOFLAGS:)
BUILD_NUMBER ?= $(BUILD_NUMBER)
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)

GO=$(GOPATH)/bin/godep go
ESLINT=node_modules/eslint/bin/eslint.js

ifeq ($(BUILD_NUMBER),)
	BUILD_NUMBER := dev
endif

ifeq ($(TRAVIS_BUILD_NUMBER),)
	BUILD_NUMBER := dev
else
	BUILD_NUMBER := $(TRAVIS_BUILD_NUMBER)
endif

DIST_ROOT=dist
DIST_PATH=$(DIST_ROOT)/mattermost

TESTS=.

DOCKERNAME ?= mm-dev
DOCKER_CONTAINER_NAME ?= mm-test

all: dist-local

dist: | build-server build-client go-test package
	mv ./model/version.go.bak ./model/version.go

dist-local: | start-docker dist

dist-travis: | travis-init build-container

start-docker:
	@echo Starting docker containers

	@if [ $(shell docker ps -a | grep -ci mattermost-mysql) -eq 0 ]; then \
		echo starting mattermost-mysql; \
		docker run --name mattermost-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=mostest \
		-e MYSQL_USER=mmuser -e MYSQL_PASSWORD=mostest -e MYSQL_DATABASE=mattermost_test -d mysql:5.7 > /dev/null; \
	elif [ $(shell docker ps | grep -ci mattermost-mysql) -eq 0 ]; then \
		echo restarting mattermost-mysql; \
		docker start mattermost-mysql > /dev/null; \
	fi

	@if [ $(shell docker ps -a | grep -ci mattermost-postgres) -eq 0 ]; then \
		echo starting mattermost-postgres; \
		docker run --name mattermost-postgres -p 5432:5432 -e POSTGRES_USER=mmuser -e POSTGRES_PASSWORD=mostest \
		-d postgres:9.4 > /dev/null; \
		sleep 10; \
	elif [ $(shell docker ps | grep -ci mattermost-postgres) -eq 0 ]; then \
		echo restarting mattermost-postgres; \
		docker start mattermost-postgres > /dev/null; \
		sleep 10; \
	fi

build-server:
	@echo Building mattermost server

	rm -Rf $(DIST_ROOT)
	$(GO) clean $(GOFLAGS) -i ./...

	@echo GOFMT
	$(eval GOFMT_OUTPUT := $(shell gofmt -d -s api/ model/ store/ utils/ manualtesting/ mattermost.go 2>&1))
	@echo "$(GOFMT_OUTPUT)"
	@if [ ! "$(GOFMT_OUTPUT)" ]; then \
		echo "gofmt sucess"; \
	else \
		echo "gofmt failure"; \
		exit 1; \
	fi

	cp ./model/version.go ./model/version.go.bak
	sed -i'.make_mac_work' 's|_BUILD_NUMBER_|$(BUILD_NUMBER)|g' ./model/version.go
	sed -i'.make_mac_work' 's|_BUILD_DATE_|$(BUILD_DATE)|g' ./model/version.go
	sed -i'.make_mac_work' 's|_BUILD_HASH_|$(BUILD_HASH)|g' ./model/version.go
	rm ./model/version.go.make_mac_work

	$(GO) build $(GOFLAGS) ./...
	$(GO) install $(GOFLAGS) ./...

package:
	@ echo Packaging mattermost

	mkdir -p $(DIST_PATH)/bin
	cp $(GOPATH)/bin/platform $(DIST_PATH)/bin

	cp -RL config $(DIST_PATH)/config
	touch $(DIST_PATH)/config/build.txt
	echo $(BUILD_NUMBER) | tee -a $(DIST_PATH)/config/build.txt

	mkdir -p $(DIST_PATH)/logs

	mkdir -p $(DIST_PATH)/web/static/js
	cp -L web/static/js/*.min.js $(DIST_PATH)/web/static/js/
	cp -L web/static/js/*.min.js.map $(DIST_PATH)/web/static/js/
	cp -RL web/static/config $(DIST_PATH)/web/static
	cp -RL web/static/css $(DIST_PATH)/web/static
	cp -RL web/static/fonts $(DIST_PATH)/web/static
	cp -RL web/static/help $(DIST_PATH)/web/static
	cp -RL web/static/images $(DIST_PATH)/web/static
	cp -RL web/static/js/jquery-dragster $(DIST_PATH)/web/static/js/
	cp -RL web/templates $(DIST_PATH)/web

	mkdir -p $(DIST_PATH)/api
	cp -RL api/templates $(DIST_PATH)/api

	cp build/MIT-COMPILED-LICENSE.md $(DIST_PATH)
	cp NOTICE.txt $(DIST_PATH)
	cp README.md $(DIST_PATH)

	mv $(DIST_PATH)/web/static/js/bundle.min.js $(DIST_PATH)/web/static/js/bundle-$(BUILD_NUMBER).min.js
	mv $(DIST_PATH)/web/static/js/libs.min.js $(DIST_PATH)/web/static/js/libs-$(BUILD_NUMBER).min.js

	sed -i'.bak' 's|react-0.14.3.js|react-0.14.3.min.js|g' $(DIST_PATH)/web/templates/head.html
	sed -i'.bak' 's|react-dom-0.14.3.js|react-dom-0.14.3.min.js|g' $(DIST_PATH)/web/templates/head.html
	sed -i'.bak' 's|jquery-2.1.4.js|jquery-2.1.4.min.js|g' $(DIST_PATH)/web/templates/head.html
	sed -i'.bak' 's|bootstrap-3.3.5.js|bootstrap-3.3.5.min.js|g' $(DIST_PATH)/web/templates/head.html
	sed -i'.bak' 's|react-bootstrap-0.28.1.js|react-bootstrap-0.28.1.min.js|g' $(DIST_PATH)/web/templates/head.html
	sed -i'.bak' 's|perfect-scrollbar-0.6.7.jquery.js|perfect-scrollbar-0.6.7.jquery.min.js|g' $(DIST_PATH)/web/templates/head.html
	sed -i'.bak' 's|bundle.js|bundle-$(BUILD_NUMBER).min.js|g' $(DIST_PATH)/web/templates/head.html
	sed -i'.bak' 's|libs.min.js|libs-$(BUILD_NUMBER).min.js|g' $(DIST_PATH)/web/templates/head.html
	rm $(DIST_PATH)/web/templates/*.bak

	tar -C dist -czf $(DIST_PATH).tar.gz mattermost

build-client:
	@echo Building mattermost web client

	cd web/react/ && npm install

	@echo Checking for style guide compliance

	@echo ESLint...
	cd web/react && $(ESLINT) --ext ".jsx" --ignore-pattern node_modules --quiet .

	cd web/react/ && npm run build-libs

	mkdir -p web/static/js
	cd web/react && npm run build

	cd web/sass-files && compass compile -e production --force

go-test:
	$(GO) test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=180s ./api || exit 1
	$(GO) test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=12s ./model || exit 1
	$(GO) test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=120s ./store || exit 1
	$(GO) test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=120s ./utils || exit 1
	$(GO) test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=120s ./web || exit 1

test: | start-docker .prepare-go go-test

travis-init:
	@echo Setting up enviroment for travis

	if [ "$(TRAVIS_DB)" = "postgres" ]; then \
		sed -i'.bak' 's|mysql|postgres|g' config/config.json; \
		sed -i'.bak' 's|mmuser:mostest@tcp(dockerhost:3306)/mattermost_test?charset=utf8mb4,utf8|postgres://mmuser:mostest@postgres:5432/mattermost_test?sslmode=disable\&connect_timeout=10|g' config/config.json; \
	fi

	if [ "$(TRAVIS_DB)" = "mysql" ]; then \
		sed -i'.bak' 's|mmuser:mostest@tcp(dockerhost:3306)/mattermost_test?charset=utf8mb4,utf8|mmuser:mostest@tcp(mysql:3306)/mattermost_test?charset=utf8mb4,utf8|g' config/config.json; \
	fi

build-container:
	@echo Building in container

	docker run -e TRAVIS_BUILD_NUMBER=$(TRAVIS_BUILD_NUMBER) --link mattermost-mysql:mysql --link mattermost-postgres:postgres -v `pwd`:/go/src/github.com/mattermost/platform mattermost/builder:latest

stop-docker:
	@echo Stopping docker containers

	@if [ $(shell docker ps -a | grep -ci mattermost-mysql) -eq 1 ]; then \
		echo stopping mattermost-mysql; \
		docker stop mattermost-mysql > /dev/null; \
	fi

	@if [ $(shell docker ps -a | grep -ci mattermost-postgres) -eq 1 ]; then \
		echo stopping mattermost-postgres; \
		docker stop mattermost-postgres > /dev/null; \
	fi

clean-docker:
	@echo Removing docker containers

	@if [ $(shell docker ps -a | grep -ci mattermost-mysql) -eq 1 ]; then \
		echo stopping mattermost-mysql; \
		docker stop mattermost-mysql > /dev/null; \
		docker rm -v mattermost-mysql > /dev/null; \
	fi

	@if [ $(shell docker ps -a | grep -ci mattermost-postgres) -eq 1 ]; then \
		echo stopping mattermost-postgres; \
		docker stop mattermost-postgres > /dev/null; \
		docker rm -v mattermost-postgres > /dev/null; \
	fi

clean: stop-docker
	rm -Rf $(DIST_ROOT)
	go clean $(GOFLAGS) -i ./...

	rm -rf web/react/node_modules
	rm -f web/static/js/bundle*.js
	rm -f web/static/js/bundle*.js.map
	rm -f web/static/js/libs*.js
	rm -f web/static/css/styles.css

	rm -rf api/data
	rm -rf logs
	rm -rf web/sass-files/.sass-cache

	rm -rf Godeps/_workspace/pkg/

	rm -f mattermost.log
	rm -f .prepare-go .prepare-jsx

nuke: | clean clean-docker
	rm -rf data

.prepare-go:
	@echo Preparation for running go code
	go get $(GOFLAGS) github.com/tools/godep

	touch $@

.prepare-jsx:
	@echo Preparation for compiling jsx code

	cd web/react/ && npm install
	cd web/react/ && npm run build-libs

	touch $@

run: start-docker .prepare-go .prepare-jsx
	mkdir -p web/static/js

	@echo Starting react processo
	cd web/react && npm start &

	@echo Starting go web server
	$(GO) run $(GOFLAGS) mattermost.go -config=config.json &

	@echo Starting compass watch
	cd web/sass-files && compass watch &

stop:
	@for PID in $$(ps -ef | grep [c]ompass | awk '{ print $$2 }'); do \
		echo stopping css watch $$PID; \
		kill $$PID; \
	done

	@for PID in $$(ps -ef | grep [n]pm | awk '{ print $$2 }'); do \
		echo stopping watchify $$PID; \
		kill $$PID; \
	done

	@for PID in $$(ps -ef | grep [m]atterm | awk '{ print $$2 }'); do \
		echo stopping go web $$PID; \
		kill $$PID; \
	done

	@if [ $(shell docker ps -a | grep -ci ${DOCKER_CONTAINER_NAME}) -eq 1 ]; then \
		echo removing dev docker container; \
		docker stop ${DOCKER_CONTAINER_NAME} > /dev/null; \
		docker rm -v ${DOCKER_CONTAINER_NAME} > /dev/null; \
	fi

setup-mac:
	echo $$(boot2docker ip 2> /dev/null) dockerhost | sudo tee -a /etc/hosts

cleandb:
	@if [ $(shell docker ps -a | grep -ci mattermost-mysql) -eq 1 ]; then \
		docker stop mattermost-mysql > /dev/null; \
		docker rm -v mattermost-mysql > /dev/null; \
	fi

docker-build: stop
	docker build -t ${DOCKERNAME} -f docker/local/Dockerfile .

docker-run: docker-build
	docker run --name ${DOCKER_CONTAINER_NAME} -d --publish 8065:80 ${DOCKERNAME}

zbox:
	@sed -i'.bak' 's|_BUILD_NUMBER_|$(BUILD_NUMBER)|g' ./model/version.go
	@sed -i'.bak' 's|_BUILD_DATE_|$(BUILD_DATE)|g' ./model/version.go
	@sed -i'.bak' 's|_BUILD_HASH_|$(BUILD_HASH)|g' ./model/version.go

	@echo Cleaning up
	@rm -Rf $(DIST_ROOT)
	@$(GO) clean $(GOFLAGS) -i ./...

	@echo Building mattermost for linux
	@env GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -i ./...
	@env GOOS=linux GOARCH=amd64 $(GO) install $(GOFLAGS) ./...

	@mkdir -p $(DIST_PATH)/bin
	@cp $(GOPATH)/bin/linux_amd64/platform $(DIST_PATH)/bin

	@mkdir -p $(DIST_PATH)/config
	@touch $(DIST_PATH)/config/build.txt
	@echo $(BUILD_NUMBER) | tee -a $(DIST_PATH)/config/build.txt

	@mkdir -p $(DIST_PATH)/logs

	@echo Building frontend scripts
	@mkdir -p web/static/js
	@cd web/react && npm run build-libs && npm run build

	@echo Building frontend styles
	@cd web/sass-files && compass compile -e production --force

	@echo Preparing mattermost
	@mkdir -p $(DIST_PATH)/web
	@cp -RL web/i18n $(DIST_PATH)/web
	@cp -RL web/static $(DIST_PATH)/web
	@cp -RL web/templates $(DIST_PATH)/web

	@mkdir -p $(DIST_PATH)/api
	@cp -RL api/templates $(DIST_PATH)/api

	@mkdir -p $(DIST_PATH)/bin/i18n
	@cp -RL i18n/ $(DIST_PATH)/bin/i18n
	@rm -f $(DIST_PATH)/bin/i18n/*.go

	@mv $(DIST_PATH)/web/static/js/bundle.min.js $(DIST_PATH)/web/static/js/bundle-$(BUILD_NUMBER).min.js

	@sed -i'.bak' 's|react-with-addons-0.14.1.js|react-with-addons-0.14.1.min.js|g' $(DIST_PATH)/web/templates/head.html
	@sed -i'.bak' 's|jquery-1.11.1.js|jquery-1.11.1.min.js|g' $(DIST_PATH)/web/templates/head.html
	@sed -i'.bak' 's|bootstrap-3.3.5.js|bootstrap-3.3.5.min.js|g' $(DIST_PATH)/web/templates/head.html
	@sed -i'.bak' 's|react-bootstrap-0.27.3.js|react-bootstrap-0.27.3.min.js|g' $(DIST_PATH)/web/templates/head.html
	@sed -i'.bak' 's|perfect-scrollbar-0.6.5.jquery.js|perfect-scrollbar-0.6.5.jquery.min.js|g' $(DIST_PATH)/web/templates/head.html
	@sed -i'.bak' 's|bundle.js|bundle-$(BUILD_NUMBER).min.js|g' $(DIST_PATH)/web/templates/head.html
	@rm $(DIST_PATH)/web/templates/*.bak


zbox-test: zbox
	@echo Building Docker Image
	@docker build -t ${DOCKERNAME}:$(BUILD_NUMBER) -f docker/test/Dockerfile .
	@rm -rf dist
	@docker tag -f ${DOCKERNAME}:$(BUILD_NUMBER) ${DOCKERNAME}:latest;

zbox-prod: zbox
	@echo Building Docker Image
	@docker build -t enahum/${DOCKERNAME}:$(BUILD_NUMBER) -f docker/prod/Dockerfile .
	@rm -rf dist
	@docker tag -f enahum/${DOCKERNAME}:$(BUILD_NUMBER) enahum/${DOCKERNAME}:latest; \



