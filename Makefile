.PHONY: gogenerate
gogenerate: .gogeneratedeps fmt
	@go generate -v $(ALL_APPS)

.PHONY: generate-all
api-generate-all: gogenerate fmt
