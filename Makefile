################################################################################
# Build Variables
################################################################################

BUILDFLAGS =
DEVFLAGS   = -race
TESTFLAGS  = -cover
LDFLAGS    = -ldflags="-X main.appVersion=$(APP_VERSION)"
POTFLAGS   = --from-code=utf-8 -kl --package-name="$(APP)"

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
POT_PATH     = po/$(POT_NAME)

################################################################################
# Automatic Variables
################################################################################

GO_SOURCES  = $(shell find . -name '*.go' -type f)
CSS_SOURCES = $(shell find . -name '*.css' -type f)
UI_SOURCES  = $(shell find . -name '*.ui' -type f)
GEN_SOURCES = $(patsubst %,%.go,$(CSS_SOURCES) $(UI_SOURCES))
SOURCES     = $(GO_SOURCES) $(GEN_SOURCES)
IMPORTS     = $(shell tools/list-imports.sh ./...)

POTFILES = $(shell cat po/POTFILES)
LINGUAS  = $(shell cat po/LINGUAS)
PO       = $(shell find po -name '*.po' -type f)
MO       = $(patsubst %.po,%.mo,$(PO))

APP_VERSION = $(shell tools/app-version.sh)
GTK_VERSION = $(shell tools/gtk-version.sh)

################################################################################
# Targets
################################################################################

all: $(EXE_PATH) $(DESKTOP_PATH) $(APPDATA_PATH) $(POT_PATH) $(MO)

deps:
	go get -u $(BUILDFLAGS) $(IMPORTS)

$(EXE_PATH): Makefile $(SOURCES)
	go build -o $@ $(BUILDFLAGS) $(LDFLAGS) ./cmd/xkcd-gtk

dev: $(GEN_SOURCES)
	go build -o $(EXE_PATH)-dev $(BUILDFLAGS) $(LDFLAGS) $(DEVFLAGS) ./cmd/xkcd-gtk

$(POT_PATH): $(POTFILES)
	xgettext -o $@ $(POTFLAGS) $^

%.css.go: %.css
	tools/go-wrap.sh $< >$@

%.ui.go: %.ui
	tools/go-wrap.sh $< >$@

%.desktop: %.desktop.in $(PO)
	msgfmt --desktop -d po -c -o $@ --template $<

%.xml: %.xml.in $(PO)
	msgfmt --xml -d po -c -o $@ --template $<

%.mo: %.po
	msgfmt -c -o $@ $<

fix: $(GEN_SOURCES)
	go fix ./...
	go fmt ./...

check: $(GEN_SOURCES) $(APPDATA_PATH)
	go vet ./...
	golint -set_exit_status ./...
	xmllint --noout $(APPDATA_PATH) $(ICON_PATH) $(UI_SOURCES)
	yamllint .
	appstream-util validate-relax $(APPDATA_PATH)

test: $(GEN_SOURCES)
	go test $(BUILDFLAGS) $(DEVFLAGS) $(TESTFLAGS) ./...
	tools/test-install.sh

clean:
	rm -f $(EXE_PATH) $(EXE_PATH)-dev $(GEN_SOURCES) $(DESKTOP_PATH) $(APPDATA_PATH) $(MO)

strip: $(EXE_PATH)
	strip $(EXE_PATH)

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

.PHONY: all check clean deps dev fix install strip test uninstall
