################################################################################
# Build Variables
################################################################################

BUILDFLAGS = -tags $(GTK_VERSION)
LDFLAGS    = -ldflags="-X main.appVersion=$(APP_VERSION)"
POTFLAGS   = --from-code=utf-8 \
             -kl \
             --package-name="$(APP)"

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
POT_NAME     = $(APP).pot

EXE_PATH     = $(EXE_NAME)
ICON_PATH    = data/$(ICON_NAME)
DESKTOP_PATH = data/$(DESKTOP_NAME)
APPDATA_PATH = data/$(APPDATA_NAME)
POT_PATH     = po/$(APP).pot

################################################################################
# Automatic Variables
################################################################################

GO_SOURCES  = $(shell find . -name '*.go' -type f)
UI_SOURCES  = $(shell find . -name '*.ui' -type f)
GEN_SOURCES = $(patsubst %,%.go,$(UI_SOURCES))
SOURCES     = $(GO_SOURCES) $(GEN_SOURCES)
DEPS        = $(shell tools/list-imports.sh ./...)

POTFILES = $(shell cat po/POTFILES)
LINGUAS  = $(shell cat po/LINGUAS)
PO       = $(shell find po -name '*.po' -type f)
MO       = $(patsubst %.po,%.mo,$(PO))

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

all: $(EXE_PATH) $(DESKTOP_PATH) $(APPDATA_PATH) $(POT_PATH) $(MO)

deps:
	go get -u $(BUILDFLAGS) $(DEPS)

$(EXE_PATH): Makefile $(SOURCES)
	go build -o $@ $(BUILDFLAGS) $(LDFLAGS) ./cmd/xkcd-gtk

$(POT_PATH): $(POTFILES)
	xgettext -o $@ $(POTFLAGS) $^

%.ui.go: %.ui
	tools/go-wrap.sh $< >$@

%.desktop: %.desktop.in $(PO)
	msgfmt --desktop -d po -c -o $@ --template $<

%.xml: %.xml.in $(PO)
	msgfmt --xml -d po -c -o $@ --template $<

%.mo: %.po
	msgfmt -c -o $@ $<

check:
	-go fmt ./...
	-go vet ./...
	-golint ./...

clean:
	-rm -f $(EXE_PATH) $(GEN_SOURCES) $(DESKTOP_PATH) $(APPDATA_PATH) $(MO)

install: $(EXE_PATH) $(DESKTOP_PATH) $(APPDATA_PATH) $(MO)
	mkdir -p $(DESTDIR)$(bindir)
	install $(EXE_PATH) $(DESTDIR)$(bindir)
	mkdir -p $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	cp $(ICON_PATH) $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	mkdir -p $(DESTDIR)$(datadir)/applications
	cp $(DESKTOP_PATH) $(DESTDIR)$(datadir)/applications
	mkdir -p $(DESTDIR)$(datadir)/metainfo
	cp $(APPDATA_PATH) $(DESTDIR)$(datadir)/metainfo
	for lang in $(LINGUAS); do \
		mkdir -p "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES"; \
		cp "po/$$lang.mo" "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES/$(APP).mo"; \
	done

uninstall:
	rm $(DESTDIR)$(bindir)/$(EXE_NAME) \
	   $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps/$(ICON_NAME) \
	   $(DESTDIR)$(datadir)/applications/$(DESKTOP_NAME) \
	   $(DESTDIR)$(datadir)/metainfo/$(APPDATA_NAME)
	for lang in $(LINGUAS); do \
		rm "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES/$(APP).mo"; \
	done

.PHONY: all check clean deps install uninstall
