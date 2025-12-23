.PHONY: cluster-up
cluster-up:
	ctlptl apply -f kind/cluster.yaml

.PHONY: cluster-down
cluster-down:
	ctlptl delete -f kind/cluster.yaml

.PHONY: lint
lint:
	golangci-lint run --fix