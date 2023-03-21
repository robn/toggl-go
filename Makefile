TOGGL ?= toggl
TAGS ?= jira
CMD_DIR := ./cmd
OUT_DIR := ./out
INTERNAL_DIR := ./internal
BIN_DIR := $(HOME)/bin

INTERNALS := $(shell find $(INTERNAL_DIR) -type f -name '*.go')
GO_TARGETS := $(addprefix $(OUT_DIR)/, $(TOGGL))

INSTALL_TARGETS := $(addprefix $(BIN_DIR)/, $(TOGGL))

all: $(GO_TARGETS)

$(OUT_DIR)/%: $(CMD_DIR)/*.go $(INTERNALS)
	go build -o $@ -tags $(TAGS) .

install: $(INSTALL_TARGETS)

$(BIN_DIR)/%: $(OUT_DIR)/%
	cp $< $@
