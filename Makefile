# shows help message defaultly
.DEFAULT_GOAL := help

#
# update
#
.PHONY: update.credits update.mocks update.swagger

# update `./CREDITS`
update.credits:
	gocredits -skip-missing . > ./CREDITS

# update mocks
update.mocks:
	# ./app/domain
	mockgen -source=./app/domain/repository.go -destination=./app/domain/repository_mock.go  -package=domain
	# ./pkg/proxy
	mockgen -source=./pkg/proxy/bubbles.go -destination=./pkg/proxy/bubbles_mock.go -package=proxy
	mockgen -source=./pkg/proxy/bubbletea.go -destination=./pkg/proxy/bubbletea_mock.go -package=proxy
	mockgen -source=./pkg/proxy/cobra.go -destination=./pkg/proxy/cobra_mock.go -package=proxy
	mockgen -source=./pkg/proxy/color.go -destination=./pkg/proxy/color_mock.go -package=proxy
	mockgen -source=./pkg/proxy/envconfig.go -destination=./pkg/proxy/envconfig_mock.go -package=proxy
	mockgen -source=./pkg/proxy/errors.go -destination=./pkg/proxy/errors_mock.go -package=proxy
	mockgen -source=./pkg/proxy/filepath.go -destination=./pkg/proxy/filepath_mock.go -package=proxy
	mockgen -source=./pkg/proxy/fmt.go -destination=./pkg/proxy/fmt_mock.go -package=proxy
	mockgen -source=./pkg/proxy/io.go -destination=./pkg/proxy/io_mock.go -package=proxy
	mockgen -source=./pkg/proxy/json.go -destination=./pkg/proxy/json_mock.go -package=proxy
	mockgen -source=./pkg/proxy/lipgloss.go -destination=./pkg/proxy/lipgloss_mock.go -package=proxy
	mockgen -source=./pkg/proxy/os.go -destination=./pkg/proxy/os_mock.go -package=proxy
	mockgen -source=./pkg/proxy/strconv.go -destination=./pkg/proxy/strconv_mock.go -package=proxy
	mockgen -source=./pkg/proxy/strings.go -destination=./pkg/proxy/strings_mock.go -package=proxy
	mockgen -source=./pkg/proxy/time.go -destination=./pkg/proxy/time_mock.go -package=proxy

#
# container
#
.PHONY: container.build container.down

# build container
container.build:
	@set -e; \
	if [ -f "./container.exist" ]; then \
		echo "container already exist"; \
		exit 1; \
	fi; \
	docker-compose -f docker-compose.yml build --no-cache; \
	touch ./container.exist

# down container
container.down:
	@set -e; \
	docker-compose down; \
	docker image prune -af; \
	if [ -f "./container.exist" ]; then \
		rm ./container.exist; \
	fi

#
# test
#
.PHONY: test.local test.container test.container.once

# execute tests in local
test.local:
	@set -e; \
	if [ -f "./test.run" ]; then \
		echo "test already running"; \
		exit 1; \
	fi; \
	touch test.run; \
	go test -v -p 1 ./... -cover -coverprofile=./cover.out; \
	grep -v -E "(_mock\.go|/mock/|/proxy/|/docs/docs\.go)" ./cover.out > ./cover.out.tmp && mv ./cover.out.tmp ./cover.out; \
	go tool cover -html=./cover.out -o ./docs/coverage.html; \
	rm ./cover.out; \
	if [ -f "./test.run" ]; then \
		rm ./test.run; \
	fi

# execute tests in container
test.container:
	@set -e; \
	if ! [ -f "./container.exist" ]; then \
		echo "container not exist"; \
		exit 1; \
	fi; \
	if [ -f "./test.run" ]; then \
		echo "test already running"; \
		exit 1; \
	fi; \
	touch test.run; \
	docker-compose -f docker-compose.yml up --abort-on-container-exit jrp-test-container; \
	CONTAINER_ID=$$(docker ps -a -q --filter "name=jrp-test-container" --filter "status=exited"); \
	docker cp $${CONTAINER_ID}:/jrp/docs/coverage.html ./docs/coverage.html; \
	rm ./test.run

# execute tests in container (once)
test.container.once:
	@set -e; \
	if [ -f "./container.exist" ]; then \
		echo "container already exist"; \
		exit 1; \
	fi; \
	if [ -f "./test.run" ]; then \
		echo "test already running"; \
		exit 1; \
	fi; \
	touch ./container.exist; \
	touch test.run; \
	docker-compose -f docker-compose.yml build --no-cache; \
	docker-compose -f docker-compose.yml up --abort-on-container-exit jrp-test-container; \
	CONTAINER_ID=$$(docker ps -a -q --filter "name=jrp-test-container" --filter "status=exited"); \
	docker cp $${CONTAINER_ID}:/jrp/docs/coverage.html ./docs/coverage.html; \
	docker-compose down; \
	docker image prune -af; \
	rm ./test.run; \
	rm ./container.exist

# required phony targets for standards
all: help
clean:
	@rm -f ./cover.out ./co
	@rm -f ./test.run ./con
	@docker-compose down
	@docker image prune -af
test: test.local

# help
.PHONY: help
help:
	@echo ""
	@echo "available targets:"
	@echo ""
	@echo "  [update]"
	@echo "    update.credits       - update ./CREDITS file"
	@echo "    update.mocks         - update all mocks"
	@echo ""
	@echo "  [container]"
	@echo "    container.build      - build container for testing"
	@echo "    container.down       - down container and remove images"
	@echo ""
	@echo "  [test]"
	@echo "    test.local           - execute all tests in local"
	@echo "    test.container       - execute all tests in container"
	@echo "    test.container.once  - build container and execute all tests in container once, then remove container and images"
	@echo ""
