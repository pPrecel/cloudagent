GARDENAGENT_PATH ?= /usr/local/bin/gardenagent
PLIST_LABEL ?= com.pPrecel.gardenagent.agent.plist
PLIST_PATH ?= ~/Library/LaunchAgents/com.pPrecel.gardenagent.agent.plist
CURRENT_DIR = $(shell pwd)

.PHONY: build
build:
	go build -o .out/gardenagent cmd/main.go

.PHONY: cp-to-path
cp-to-path:
	cp "$(CURRENT_DIR)/.out/gardenagent" $(GARDENAGENT_PATH)

.PHONY: rm-from-path
rm-from-path:
	rm $(GARDENAGENT_PATH)

.PHONY: ln-to-path
ln-to-path:
	ln -s -f "$(CURRENT_DIR)/.out/gardenagent" $(GARDENAGENT_PATH)

.PHONY: install-agent
install-agent:
	gardenagent generate plist -k $(kubeconfigPath) -n $(namespace) --cronSpec "@every 60s" $(other_flags) > $(PLIST_PATH)
	launchctl load -w $(PLIST_PATH)

.PHONY: uninstall-agent
uninstall-agent:
	launchctl remove $(PLIST_LABEL)
	rm $(PLIST_PATH)

.PHONY: protobuf
protobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/agent/proto/route.proto

.PHONY: verify
verify:
	@./hack/verify.sh
