VERSION ?= $(shell git tag --sort=committerdate | tail -n 1)
NEW_VERSION ?= v0.0.1

ifeq ($(VERSION),)
	BUMPED_VERSION := $(NEW_VERSION)
else
	BUMPED_VERSION := $(shell echo $(VERSION) | awk -F. -v OFS=. '{$$NF += 1; print}')
endif

.PHONY: tag
tag:
	@git add .
	@git commit -m "Release $(BUMPED_VERSION)" || true
	@git tag $(BUMPED_VERSION)
	@git push origin main
	@git push origin $(BUMPED_VERSION)
	@echo "✅ Tagged and pushed version $(BUMPED_VERSION)"
