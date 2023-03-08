################################################################################
# Build Variables
################################################################################

BUILDFLAGS = -v
DEVFLAGS   = $(BUILDFLAGS) -race
TESTFLAGS  = $(DEVFLAGS) -cover
VETFLAGS   = $(BUILDFLAGS)
POTFLAGS   = --package-name="$(APP)" --from-code=utf-8 --sort-output
FPBFLAGS   = --user --force-clean

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

EXE_PATH     = $(APP)
DEV_PATH     = $(APP)-dev
ICON_PATH    = data/$(APP).svg
DESKTOP_PATH = data/$(APP).desktop
SERVICE_PATH = data/$(APP).service
APPDATA_PATH = data/$(APP).appdata.xml
POT_PATH     = po/$(APP).pot

MODULE_PACKAGES = $(MODULE)/cmd/... $(MODULE)/internal/...
BUILD_PACKAGE   = $(MODULE)/internal/build
BUILD_DATA      = app-id=$(APP),version=$(APP_VERSION)
TAGS            = $(shell tools/gtk-version.sh) $(shell tools/pango-version.sh)

GO_SOURCES  = $(shell find cmd internal -name '*.go' -type f)
CSS_SOURCES = $(shell find cmd internal -name '*.css' -type f)
UI_SOURCES  = $(shell find cmd internal -name '*.ui' -type f)
SH_SOURCES  = $(shell find tools -name '*.sh' -type f)

POTFILES = $(shell cat po/POTFILES)
LINGUAS  = $(shell cat po/LINGUAS)
PO       = $(shell find po -name '*.po' -type f)
MO       = $(patsubst %.po,%.mo,$(PO))

FLATPAK_YML_IN = $(shell find flatpak -name '*.yml.in')
FLATPAK_YML    = $(APP).yml $(patsubst %.in,%,$(FLATPAK_YML_IN))

################################################################################
# Targets
################################################################################

all: $(EXE_PATH) $(DESKTOP_PATH) $(APPDATA_PATH) $(POT_PATH) $(MO)

$(EXE_PATH): Makefile $(GO_SOURCES) $(CSS_SOURCES) $(UI_SOURCES) $(APPDATA_PATH)
	go build -o $@ -ldflags="-X '$(BUILD_PACKAGE).data=$(BUILD_DATA)'" -tags "$(TAGS)" $(BUILDFLAGS) $(MODULE)/cmd/xkcd-gtk

dev: $(APPDATA_PATH)
	go build -o $(DEV_PATH) -ldflags="-X $(BUILD_PACKAGE).data=$(BUILD_DATA)" -tags "$(TAGS) xkcd_gtk_debug" $(DEVFLAGS) $(MODULE)/cmd/xkcd-gtk

$(POT_PATH): $(POTFILES) tools/fill-pot-header.sh
	xgettext -o $@ -LC -kl $(POTFLAGS) $(filter %.go,$(POTFILES))
	xgettext -o $@ -j $(POTFLAGS) $(filter %.ui,$(POTFILES))
	xgettext -o $@ -j -k -kName -kGenericName -kComment -kKeywords $(POTFLAGS) $(filter %.desktop.in,$(POTFILES))
	xgettext -o $@ -j --its=po/appdata.its $(POTFLAGS) $(filter %.xml.in,$(POTFILES))
	tools/fill-pot-header.sh <$@ >$@.out
	mv $@.out $@

%.desktop: %.desktop.in $(PO)
	msgfmt --desktop -d po -c -o $@ --template $<

%.xml: %.xml.in $(PO)
	msgfmt --xml -d po -c -o $@ --template $<

%.mo: %.po
	msgfmt -c -o $@ $<

flatpak/%.yml: flatpak/%.yml.in go.mod go.sum tools/gen-flatpak-deps.sh $(GO_SOURCES)
	cp $< $@.tmp
	tools/gen-flatpak-deps.sh >>$@.tmp
	mv $@.tmp $@

flatpak/modules.txt: go.mod go.sum $(GO_SOURCES)
	go mod vendor -o flatpak-build/vendor
	cp flatpak-build/vendor/modules.txt $@

flathub: flatpak/flathub.yml flatpak/modules.txt
	flatpak-builder $(FPBFLAGS) --state-dir=flatpak-build/.flatpak-builder-$@/ --install-deps-from=flathub flatpak-build/$@/ $<

flathub-install: flatpak/flathub.yml flatpak/modules.txt
	flatpak-builder $(FPBFLAGS) --state-dir=flatpak-build/.flatpak-builder-$@/ --install-deps-from=flathub --install flatpak-build/$@/ $<

appcenter: flatpak/appcenter.yml flatpak/modules.txt
	flatpak-builder $(FPBFLAGS) --state-dir=flatpak-build/.flatpak-builder-$@/ --install-deps-from=appcenter flatpak-build/$@/ $<

appcenter-install: flatpak/appcenter.yml flatpak/modules.txt
	flatpak-builder $(FPBFLAGS) --state-dir=flatpak-build/.flatpak-builder-$@/ --install-deps-from=appcenter --install flatpak-build/$@/ $<

$(APP).yml: flatpak/appcenter.yml
	sed "s/path: '..'/path: '.'/" $< >$@

fix: $(POT_PATH) $(PO) $(APP).yml
	go fix $(MODULE_PACKAGES)
	go fmt $(MODULE_PACKAGES)
	go mod tidy
	([ -d vendor ] && go mod vendor) || true
	echo $(UI_SOURCES) | xargs -n1 gtk-builder-tool simplify --replace
	dos2unix -q po/LINGUAS po/POTFILES po/appdata.its $(POT_PATH) $(PO)
	for lang in $(LINGUAS); do \
		msgmerge -U --backup=none "po/$$lang.po" $(POT_PATH); \
	done

check: $(APPDATA_PATH) $(FLATPAK_YML)
	go vet -tags "$(TAGS)" $(VETFLAGS) $(MODULE_PACKAGES)
	shellcheck $(SH_SOURCES)
	xmllint --noout $(APPDATA_PATH) $(ICON_PATH) $(UI_SOURCES)
	yamllint .github/workflows/*.yml $(FLATPAK_YML)
	appstream-util --nonet validate-relax $(APPDATA_PATH)
	-appstream-util validate-strict $(APPDATA_PATH)

test: $(FLATPAK_YML) $(APPDATA_PATH)
	go test -ldflags="-X $(BUILD_PACKAGE).data=$(BUILD_DATA)" -tags "$(TAGS)" $(TESTFLAGS) $(MODULE_PACKAGES)
	tools/test-flatpak-config.sh $(FLATPAK_YML)
	tools/test-install.sh

# Shorthand for all the targets that CI covers.
ci: all check test
ci-full: ci appcenter flathub

clean:
	rm -f $(EXE_PATH)
	rm -f $(DEV_PATH)
	rm -f $(DESKTOP_PATH)
	rm -f $(APPDATA_PATH)
	rm -f $(MO)
	rm -f flatpak/*.yml
	rm -f flatpak/modules.txt
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
	rm $(DESTDIR)$(bindir)/$(notdir $(EXE_PATH))
	rm $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps/$(notdir $(ICON_PATH))
	rm $(DESTDIR)$(datadir)/applications/$(notdir $(DESKTOP_PATH))
	rm $(DESTDIR)$(datadir)/dbus-1/services/$(notdir $(SERVICE_PATH))
	rm $(DESTDIR)$(datadir)/metainfo/$(notdir $(APPDATA_PATH))
	for lang in $(LINGUAS); do \
		rm "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES/$(APP).mo"; \
	done

.PHONY: all appcenter appcenter-install check ci ci-full clean dev fix flathub flathub-install install strip test uninstall
