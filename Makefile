################################################################################
# Build Variables
################################################################################

BUILDFLAGS =
DEVFLAGS   = -race
TESTFLAGS  = -cover
POTFLAGS   = --package-name="$(APP)" --from-code=utf-8 --sort-output

APP_VERSION   = $(shell tools/app-version.sh)
GTK_VERSION   = $(shell tools/gtk-version.sh)
PANGO_VERSION = $(shell tools/pango-version.sh)

# Comma separated
BUILD_DATA     = version=$(APP_VERSION)
DEV_BUILD_DATA = debug=true
# Space separated
TAGS           = $(GTK_VERSION) $(PANGO_VERSION)
DEV_TAGS       = xkcd_gtk_debug

################################################################################
# Install Variables
################################################################################

# The `DESTDIR` variable is supported by the `install`/`uninstall` targets and
# should be overridden instead of `prefix` when building packages unless a
# custom install location, e.g. `/app` for flatpaks, is required.

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
DEV_PATH     = $(EXE_PATH)-dev
ICON_PATH    = data/$(ICON_NAME)
DESKTOP_PATH = data/$(DESKTOP_NAME)
APPDATA_PATH = data/$(APPDATA_NAME)
POT_PATH     = po/$(POT_NAME)

MODULE_PACKAGES = $(MODULE)/cmd/... $(MODULE)/internal/...
BUILD_PACKAGE   = $(MODULE)/internal/build

GO_SOURCES     = $(shell find cmd internal -name '*.go' -type f)
CSS_SOURCES    = $(shell find cmd internal -name '*.css' -type f)
UI_SOURCES     = $(shell find cmd internal -name '*.ui' -type f)
GEN_SOURCES    = $(patsubst %,%.go,$(CSS_SOURCES) $(UI_SOURCES))
ALL_GO_SOURCES = $(GO_SOURCES) $(GEN_SOURCES)
SH_SOURCES     = $(shell find tools -name '*.sh' -type f)

POTFILES         = $(shell cat po/POTFILES)
POTFILES_GO      = $(filter %.go,$(POTFILES))
POTFILES_UI      = $(filter %.ui,$(POTFILES))
POTFILES_DESKTOP = $(filter %.desktop.in,$(POTFILES))
POTFILES_APPDATA = $(filter %.xml.in,$(POTFILES))

LINGUAS = $(shell cat po/LINGUAS)
PO      = $(shell find po -name '*.po' -type f)
MO      = $(patsubst %.po,%.mo,$(PO))

FLATPAK_YML_IN  = $(shell find flatpak -name '*.yml.in')
GEN_FLATPAK_YML = $(patsubst %.in,%,$(FLATPAK_YML_IN))
FLATPAK_YML     = $(APP).yml $(GEN_FLATPAK_YML)

################################################################################
# Local Customizations (not tracked by source control)
################################################################################

-include .config.mk

################################################################################
# Targets
################################################################################

all: $(EXE_PATH) $(DESKTOP_PATH) $(APPDATA_PATH) $(POT_PATH) $(MO)

$(EXE_PATH): Makefile $(ALL_GO_SOURCES) $(APPDATA_PATH)
	go build -o $@ -v -ldflags="-X '$(BUILD_PACKAGE).data=$(BUILD_DATA)'" -tags "$(TAGS)" $(BUILDFLAGS) $(MODULE)/cmd/xkcd-gtk

dev: $(GEN_SOURCES) $(APPDATA_PATH)
	go build -o $(DEV_PATH) -v -ldflags="-X $(BUILD_PACKAGE).data=$(BUILD_DATA),$(DEV_BUILD_DATA)" -tags "$(TAGS) $(DEV_TAGS)" $(BUILDFLAGS) $(DEVFLAGS) $(MODULE)/cmd/xkcd-gtk

%.css.go: %.css tools/go-wrap.sh
	tools/go-wrap.sh $< >$@

%.ui.go: %.ui tools/go-wrap.sh
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

flatpak/%.yml: flatpak/%.yml.in go.mod go.sum tools/gen-flatpak-deps.sh $(ALL_GO_SOURCES)
	cp $< $@
	tools/gen-flatpak-deps.sh >>$@

flathub: flatpak/flathub.yml
	flatpak-builder --user --install-deps-from=flathub --force-clean $(FPBFLAGS) flatpak-build/flathub/ $<

appcenter: flatpak/appcenter.yml
	flatpak-builder --user --install-deps-from=appcenter --force-clean $(FPBFLAGS) flatpak-build/appcenter/ $<

$(APP).yml: flatpak/appcenter.yml
	sed "s/path: '..'/path: '.'/" $< >$@

fix: $(GEN_SOURCES) $(POT_PATH) $(PO) $(APP).yml
	go fix $(MODULE_PACKAGES)
	go fmt $(MODULE_PACKAGES)
	go mod tidy
	([ -d vendor ] && go mod vendor) || true
	dos2unix -q po/LINGUAS po/POTFILES po/appdata.its $(POT_PATH) $(PO)
	for lang in $(LINGUAS); do \
		msgmerge -U --backup=none "po/$$lang.po" $(POT_PATH); \
	done

check: $(GEN_SOURCES) $(APPDATA_PATH) $(FLATPAK_YML)
	go vet -tags "$(TAGS)" $(BUILDFLAGS) $(MODULE_PACKAGES)
	shellcheck $(SH_SOURCES)
	xmllint --noout $(APPDATA_PATH) $(ICON_PATH) $(UI_SOURCES)
	yamllint .github/workflows/*.yml $(FLATPAK_YML)
	appstream-util --nonet validate-relax $(APPDATA_PATH)
	-appstream-util validate-strict $(APPDATA_PATH)

test: $(GEN_SOURCES) $(FLATPAK_YML) $(APPDATA_PATH)
	go test -ldflags="-X $(BUILD_PACKAGE).data=$(BUILD_DATA)" -tags "$(TAGS)" $(BUILDFLAGS) $(DEVFLAGS) $(TESTFLAGS) $(MODULE_PACKAGES)
	tools/test-flatpak-config.sh $(FLATPAK_YML)
	tools/test-install.sh

# Shorthand for all the targets that CI covers.
ci: all check test

clean:
	rm -f $(EXE_PATH)
	rm -f $(DEV_PATH)
	rm -f $(GEN_SOURCES)
	rm -f $(DESKTOP_PATH)
	rm -f $(APPDATA_PATH)
	rm -f $(MO)
	rm -f $(GEN_FLATPAK_YML)
	rm -rf flatpak-build/
	rm -rf .flatpak-builder/

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
	rm $(DESTDIR)$(bindir)/$(EXE_NAME)
	rm $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps/$(ICON_NAME)
	rm $(DESTDIR)$(datadir)/applications/$(DESKTOP_NAME)
	rm $(DESTDIR)$(datadir)/metainfo/$(APPDATA_NAME)
	for lang in $(LINGUAS); do \
		rm "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES/$(APP).mo"; \
	done

.PHONY: all appcenter check ci clean dev fix flathub install strip test uninstall
