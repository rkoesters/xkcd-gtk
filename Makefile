################################################################################
# Build Variables
################################################################################

BUILDFLAGS = -tags "$(GTK_VERSION) $(PANGO_VERSION)" -mod=vendor
DEVFLAGS   = -race
TESTFLAGS  = -cover
LDFLAGS    = -ldflags="-X main.appVersion=$(APP_VERSION)"
POTFLAGS   = --package-name="$(APP)" --from-code=utf-8 --sort-output

################################################################################
# Install Variables
################################################################################

prefix  = /usr
bindir  = $(prefix)/bin
datadir = $(prefix)/share

################################################################################
# Application Variables
################################################################################

APP    = com.github.rkoesters.xkcd-gtk
MODULE = github.com/rkoesters/xkcd-gtk

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

DEV_PATH = $(EXE_PATH)-dev

POTFILES         = $(shell cat po/POTFILES)
POTFILES_GO      = $(filter %.go,$(POTFILES))
POTFILES_UI      = $(filter %.ui,$(POTFILES))
POTFILES_DESKTOP = $(filter %.desktop.in,$(POTFILES))
POTFILES_APPDATA = $(filter %.xml.in,$(POTFILES))

LINGUAS = $(shell cat po/LINGUAS)
PO      = $(shell find po -name '*.po' -type f)
MO      = $(patsubst %.po,%.mo,$(PO))

APP_VERSION   = $(shell tools/app-version.sh)
GTK_VERSION   = $(shell tools/gtk-version.sh)
PANGO_VERSION = $(shell tools/pango-version.sh)

################################################################################
# Targets
################################################################################

all: $(EXE_PATH) $(DESKTOP_PATH) $(APPDATA_PATH) $(POT_PATH) $(MO)

$(EXE_PATH): Makefile $(SOURCES)
	go build -o $@ -v $(LDFLAGS) $(BUILDFLAGS) $(MODULE)/cmd/xkcd-gtk

dev: $(GEN_SOURCES)
	go build -o $(DEV_PATH) -v $(LDFLAGS) $(BUILDFLAGS) $(DEVFLAGS) $(MODULE)/cmd/xkcd-gtk

vendor:
	go mod vendor
	rm vendor/github.com/gotk3/gotk3/gtk/gtk_since_3_24.go

%.css.go: %.css
	tools/go-wrap.sh $< >$@

%.ui.go: %.ui
	tools/go-wrap.sh $< >$@

$(POT_PATH): $(POTFILES) tools/fill-pot-header.sh
	xgettext -o $@ -LC -kl $(POTFLAGS) $(POTFILES_GO)
	xgettext -o $@ -j $(POTFLAGS) $(POTFILES_UI)
	xgettext -o $@ -j -k -kName -kGenericName -kComment -kKeywords $(POTFLAGS) $(POTFILES_DESKTOP)
	xgettext -o $@ -j --its=po/appdata.its $(POTFLAGS) $(POTFILES_APPDATA)
	tools/fill-pot-header.sh <$@ >$@.out
	mv $@.out $@

%.desktop: %.desktop.in $(PO)
	msgfmt --desktop -d po -c -o $@ --template $<

%.xml: %.xml.in $(PO)
	msgfmt --xml -d po -c -o $@ --template $<

%.mo: %.po
	msgfmt -c -o $@ $<

fix: $(GEN_SOURCES)
	go fix $(MODULE)/...
	go fmt $(MODULE)/...

check: $(GEN_SOURCES) $(APPDATA_PATH)
	go vet $(BUILDFLAGS) $(MODULE)/...
	golint -set_exit_status $(MODULE)/...
	xmllint --noout $(APPDATA_PATH) $(ICON_PATH) $(UI_SOURCES)
	yamllint .github/workflows/*.yml
	-appstream-util validate-relax $(APPDATA_PATH)

test: $(GEN_SOURCES)
	go test $(BUILDFLAGS) $(DEVFLAGS) $(TESTFLAGS) $(MODULE)/...
	tools/test-install.sh

clean:
	rm -f $(EXE_PATH) $(DEV_PATH) $(GEN_SOURCES) $(DESKTOP_PATH) $(APPDATA_PATH) $(MO)

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

.PHONY: all check clean dev fix install strip test uninstall vendor
