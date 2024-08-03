MAKEFLAGS += --check-symlink-times
MAKEFLAGS += --jobs
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
MAKEFLAGS += --shuffle
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.DELETE_ON_ERROR:
.ONESHELL:
.SHELLFLAGS := --norc --noprofile -Eeuo pipefail -O dotglob -O nullglob -O extglob -O failglob -O globstar -c

.DEFAULT_GOAL := package

.PHONY: clean clobber package ci

DIST := dist
GIT_TAG := $(shell git describe --tags --always --dirty)
VERSION := $(patsubst v%,%,$(GIT_TAG))
NAME_PREFIX := $(notdir $(PWD))_$(VERSION)
MANIFEST := $(NAME_PREFIX)_manifest.json
SHA_FILE := $(NAME_PREFIX)_SHA256SUMS
SIG_FILE := $(SHA_FILE).sig

clean:
	shopt -u failglob
	rm -v -rf -- '$(DIST)'

clobber: clean
	shopt -u failglob
	rm -v -rf --

$(DIST): .goreleaser.yml main.go
	GORELEASER_CURRENT_TAG='$(GIT_TAG)' goreleaser release --clean --skip validate,publish

$(DIST)/$(MANIFEST): terraform-registry-manifest.json $(DIST)
	cp -v -- '$<' '$@'

$(DIST)/$(SHA_FILE): $(DIST)/$(MANIFEST)
	cd -- '$(DIST)'
	sha256sum -- '$(<F)' *.zip > '$(@F)'

$(DIST)/$(SIG_FILE): $(DIST)/$(SHA_FILE)
	gpg --batch --yes --detach-sign --output '$@' -- '$<'

package: $(DIST)/$(SIG_FILE)

ci: package
	cd -- '$(DIST)'
	gh release create -- '$(GIT_TAG)' $(MANIFEST) $(SHA_FILE) $(SIG_FILE) *.zip
