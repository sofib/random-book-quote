TOOLS_PATH ?= $(shell pwd)/tools/

.PHONY: cluster-up
cluster-up:
	${TOOLS_PATH}ctlptl apply -f kind/cluster.yaml
	${TOOLS_PATH}ctlptl apply -f kind/registry.yaml

.PHONY: cluster-down
cluster-down:
	${TOOLS_PATH}ctlptl delete -f kind/cluster.yaml

.PHONY: lint
lint:
	golangci-lint run --fix