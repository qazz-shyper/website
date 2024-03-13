NAME = website

VERSION=$(shell git describe --always --dirty)

LDFLAGS = -X github.com/xtls/xray-core/core.build=$(VERSION) -s -w -buildid=
PARAMS = -trimpath -ldflags "$(LDFLAGS)" -v
MAIN = ./main
PREFIX ?= $(shell go env GOPATH)
ifeq ($(GOOS),windows)
OUTPUT = $(NAME).exe
ADDITION = go build -o build_assets/w$(NAME).exe -trimpath -ldflags "-H windowsgui $(LDFLAGS)" -v $(MAIN)
else
OUTPUT = $(NAME)
endif
ifeq ($(shell echo "$(GOARCH)" | grep -Pq "(mips|mipsle)" && echo true),true) # 
ADDITION = GOMIPS=softfloat go build -o build_assets/$(NAME)_softfloat -trimpath -ldflags "$(LDFLAGS)" -v $(MAIN)
endif
.PHONY: clean

build:
	go build -o $(OUTPUT) $(PARAMS) $(MAIN)
	$(ADDITION)

install:
	go build -o $(PREFIX)/bin/$(OUTPUT) $(PARAMS) $(MAIN)

clean:
	go clean -v -i $(PWD)
	rm -f website website.exe wwebsite.exe website_softfloat
