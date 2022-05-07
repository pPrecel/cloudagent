CLOUDAGENT_PATH ?= /usr/local/bin/cloudagent
PLIST_LABEL ?= com.pPrecel.cloudagent.agent.plist
PLIST_PATH ?= ~/Library/LaunchAgents/com.pPrecel.cloudagent.agent.plist
CURRENT_DIR = $(shell pwd)

.PHONY: build
build:
	go build -o .out/cloudagent main.go

.PHONY: rm-from-path
rm-from-path:
	rm $(CLOUDAGENT_PATH)

.PHONY: ln-to-path
ln-to-path:
	ln -s -f "$(CURRENT_DIR)/.out/cloudagent" $(CLOUDAGENT_PATH)

.PHONY: install-agent
install-agent: bootstrap-config
	cloudagent generate plist $(other_flags) > $(PLIST_PATH)
	launchctl load -w $(PLIST_PATH)

.PHONY: bootstrap-config
bootstrap-config:
	@./hack/config_template.sh

.PHONY: uninstall-agent
uninstall-agent:
	launchctl remove $(PLIST_LABEL)
	rm $(PLIST_PATH)

.PHONY: protobuf
protobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/agent/proto/route.proto

.PHONY: verify
verify:
	@./hack/verify.sh
