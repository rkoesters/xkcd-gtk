################################################################################
# Build Variables
################################################################################

BUILDFLAGS = -tags $(GTK_VERSION)
LDFLAGS    = -ldflags="-X main.appVersion=$(APP_VERSION)"

################################################################################
# Install Variables
################################################################################

prefix  = /usr
bindir  = $(prefix)/bin
datadir = $(prefix)/share

################################################################################
# Application Variables
################################################################################

APP = com.github.rkoesters.xkcd-gtk

EXE_NAME     = $(APP)
ICON_NAME    = $(APP).svg
DESKTOP_NAME = $(APP).desktop
APPDATA_NAME = $(APP).appdata.xml

EXE_PATH     = $(EXE_NAME)
ICON_PATH    = data/$(ICON_NAME)
DESKTOP_PATH = data/$(DESKTOP_NAME)
APPDATA_PATH = data/$(APPDATA_NAME)

################################################################################
# Automatic Variables
################################################################################

SOURCES = $(shell find . -type f -name '*.go')
DEPS    = $(shell tools/list-imports.sh ./...)

APP_VERSION = $(shell tools/app-version.sh)
GTK_VERSION = $(shell tools/gtk-version.sh)

# If GOPATH isn't set, then just use the current directory.
ifeq "$(shell go env GOPATH)" ""
export GOPATH = $(shell pwd)
endif
ifeq "$(shell go env GOPATH)" "/nonexistent/go"
export GOPATH = $(shell pwd)
endif

################################################################################
# Targets
################################################################################

all: $(EXE_PATH)

deps:
	go get -u $(BUILDFLAGS) $(DEPS)

$(EXE_PATH): Makefile $(SOURCES)
	go build -o $@ $(BUILDFLAGS) $(LDFLAGS) ./cmd/xkcd-gtk

clean:
	-go clean ./...
	-rm -f $(EXE_PATH)

install: $(EXE_PATH)
	mkdir -p $(DESTDIR)$(bindir)
	install $(EXE_PATH) $(DESTDIR)$(bindir)
	mkdir -p $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	cp $(ICON_PATH) $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	mkdir -p $(DESTDIR)$(datadir)/applications
	cp $(DESKTOP_PATH) $(DESTDIR)$(datadir)/applications
	mkdir -p $(DESTDIR)$(datadir)/metainfo
	cp $(APPDATA_PATH) $(DESTDIR)$(datadir)/metainfo

uninstall:
	rm -f $(DESTDIR)$(bindir)/$(EXE_NAME) \
	      $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps/$(ICON_NAME) \
	      $(DESTDIR)$(datadir)/applications/$(DESKTOP_NAME) \
	      $(DESTDIR)$(datadir)/metainfo/$(APPDATA_NAME)

.PHONY: all clean deps install uninstall
