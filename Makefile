.PHONY: proto
proto:
	@./scripts/proto.sh catalogwriteservice
	@./scripts/proto.sh orderservice 