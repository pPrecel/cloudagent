CLOUDAGENT_PATH ?= /usr/local/bin/cloudagent-dev
CURRENT_DIR = $(shell pwd)

.PHONY: build
build:
	./hack/build.sh

.PHONY: rm-from-path
rm-from-path:
	rm $(CLOUDAGENT_PATH)

.PHONY: ln-to-path
ln-to-path:
	ln -s -f "$(CURRENT_DIR)/.out/cloudagent" $(CLOUDAGENT_PATH)

.PHONY: bootstrap-config
bootstrap-config:
	./hack/config_template.sh

.PHONY: protobuf
protobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/agent/proto/route.proto

.PHONY: verify
verify:
	./hack/verify.sh

.PHONY: verify-proto
verify-proto:
	./hack/verify-proto.sh

.PHONY: verify-ci
verify-ci: verify verify-proto
