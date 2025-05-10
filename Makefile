VERSION ?= $(shell git describe --tags --abbrev=0)
NEW_VERSION ?= $(shell echo $(VERSION) | awk -F. '{$$NF+=1; OFS="."; print $$0}')

.PHONY: tag
tag:
	@git add .
	@git commit -m "Release $(NEW_VERSION)" || true
	@git tag $(NEW_VERSION)
	@git push origin main
	@git push origin $(NEW_VERSION)
	@echo "✅ Tagged and pushed version $(NEW_VERSION)"
