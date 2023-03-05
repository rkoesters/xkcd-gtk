################################################################################
# Build Variables
################################################################################

BUILDFLAGS =
DEVFLAGS   = -race
TESTFLAGS  = -cover
POTFLAGS   = --package-name="$(APP)" --from-code=utf-8 --sort-output

GTK_VERSION   = $(shell tools/gtk-version.sh)
PANGO_VERSION = $(shell tools/pango-version.sh)

# Space separated
TAGS     = $(GTK_VERSION) $(PANGO_VERSION)
DEV_TAGS = xkcd_gtk_debug

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

APP_VERSION = $(shell tools/app-version.sh)

EXE_NAME     = $(APP)
ICON_NAME    = $(APP).svg
DESKTOP_NAME = $(APP).desktop
SERVICE_NAME = $(APP).service
APPDATA_NAME = $(APP).appdata.xml
POT_NAME     = $(APP).pot

EXE_PATH     = $(EXE_NAME)
DEV_PATH     = $(EXE_PATH)-dev
ICON_PATH    = data/$(ICON_NAME)
DESKTOP_PATH = data/$(DESKTOP_NAME)
SERVICE_PATH = data/$(SERVICE_NAME)
APPDATA_PATH = data/$(APPDATA_NAME)
POT_PATH     = po/$(POT_NAME)

MODULE_PACKAGES = $(MODULE)/cmd/... $(MODULE)/internal/...
BUILD_PACKAGE   = $(MODULE)/internal/build
BUILD_DATA      = app-id=$(APP),version=$(APP_VERSION)

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
# Targets
################################################################################

all: $(EXE_PATH) $(DESKTOP_PATH) $(APPDATA_PATH) $(POT_PATH) $(MO)

$(EXE_PATH): Makefile $(ALL_GO_SOURCES) $(APPDATA_PATH)
	go build -o $@ -v -ldflags="-X '$(BUILD_PACKAGE).data=$(BUILD_DATA)'" -tags "$(TAGS)" $(BUILDFLAGS) $(MODULE)/cmd/xkcd-gtk

dev: $(GEN_SOURCES) $(APPDATA_PATH)
	go build -o $(DEV_PATH) -v -ldflags="-X $(BUILD_PACKAGE).data=$(BUILD_DATA)" -tags "$(TAGS) $(DEV_TAGS)" $(BUILDFLAGS) $(DEVFLAGS) $(MODULE)/cmd/xkcd-gtk

go.mod: $(ALL_GO_SOURCES)
go.sum: go.mod $(ALL_GO_SOURCES)
	go mod tidy

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
	cp $< $@.tmp
	tools/gen-flatpak-deps.sh >>$@.tmp
	mv $@.tmp $@

# Generate flatpak/modules.txt (a copy of vendor/modules.txt) without touching
# the module's vendor directory.
flatpak/modules.txt: go.mod go.sum
	go mod vendor -o flatpak-build/vendor/
	cp flatpak-build/vendor/modules.txt $@

# Generate vendor/modules.txt without network access (using the cached
# flatpak/modules.txt). If you have network access, then use 'go mod vendor'
# instead.
vendor/modules.txt:
	mkdir -p vendor
	cp flatpak/modules.txt $@

flathub: flatpak/flathub.yml
	flatpak-builder --user --install-deps-from=flathub --state-dir=flatpak-build/.flatpak-builder-$@/ --force-clean $(FPBFLAGS) flatpak-build/$@/ $<

flathub-install: flatpak/flathub.yml
	flatpak-builder --user --install --install-deps-from=flathub --state-dir=flatpak-build/.flatpak-builder-$@/ --force-clean $(FPBFLAGS) flatpak-build/$@/ $<

appcenter: flatpak/appcenter.yml
	flatpak-builder --user --install-deps-from=appcenter --state-dir=flatpak-build/.flatpak-builder-$@/ --force-clean $(FPBFLAGS) flatpak-build/$@/ $<

appcenter-install: flatpak/appcenter.yml
	flatpak-builder --user --install --install-deps-from=appcenter --state-dir=flatpak-build/.flatpak-builder-$@/ --force-clean $(FPBFLAGS) flatpak-build/$@/ $<

$(APP).yml: flatpak/appcenter.yml
	sed "s/path: '..'/path: '.'/" $< >$@

fix: $(GEN_SOURCES) $(POT_PATH) $(PO) $(APP).yml flatpak/modules.txt go.sum
	go fix $(MODULE_PACKAGES)
	go fmt $(MODULE_PACKAGES)
	([ -d vendor ] && go mod vendor) || true
	echo $(UI_SOURCES) | xargs -n1 gtk-builder-tool simplify --replace
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
ci: all appcenter check flathub test

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
	mkdir -p $(DESTDIR)$(datadir)/dbus-1/services
	cp $(SERVICE_PATH) $(DESTDIR)$(datadir)/dbus-1/services
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
	rm $(DESTDIR)$(datadir)/dbus-1/services/$(SERVICE_NAME)
	rm $(DESTDIR)$(datadir)/metainfo/$(APPDATA_NAME)
	for lang in $(LINGUAS); do \
		rm "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES/$(APP).mo"; \
	done

.PHONY: all appcenter appcenter-install check ci clean dev fix flathub flathub-install install strip test uninstall
