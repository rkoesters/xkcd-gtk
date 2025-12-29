################################################################################
# Build Variables
################################################################################

BUILDFLAGS = -v
DEVFLAGS   = -v -race
TESTFLAGS  = -race -cover
VETFLAGS   = -v
POTFLAGS   = --from-code=utf-8 --sort-output
FPBFLAGS   = --user --force-clean
ASVFLAGS   = --explain --pedantic

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

EXEC    = $(APP)
ICON    = data/$(APP).svg
DESKTOP = data/$(APP).desktop
SERVICE = data/$(APP).service
APPDATA = data/$(APP).appdata.xml

MODULE_PACKAGES = $(MODULE)/cmd/... $(MODULE)/internal/...
BUILD_PACKAGE   = $(MODULE)/internal/build
BUILD_DATA      = app-id=$(APP),version=$(shell tools/app-version.sh)
TAGS            = $(shell tools/pkg-config-tags.sh)

GO_SOURCES  = $(shell find cmd internal -name '*.go' -type f)
CSS_SOURCES = $(shell find cmd internal -name '*.css' -type f)
UI_SOURCES  = $(shell find cmd internal -name '*.ui' -type f)
SH_SOURCES  = $(shell find tools -name '*.sh' -type f)

POT      = po/$(APP).pot
LINGUAS  = $(shell cat po/LINGUAS)
PO       = $(shell find po -name '*.po' -type f)
MO       = $(patsubst %.po,%.mo,$(PO))

FLATPAK_YML_IN = $(shell find flatpak -name '*.yml.in')
FLATPAK_YML    = $(APP).yml $(patsubst %.in,%,$(FLATPAK_YML_IN))

################################################################################
# Targets
################################################################################

all: $(EXEC) $(DESKTOP) $(APPDATA) $(POT) $(MO)

$(EXEC): Makefile $(GO_SOURCES) $(CSS_SOURCES) $(UI_SOURCES) $(APPDATA)
	go build -o $@ -ldflags="-X '$(BUILD_PACKAGE).data=$(BUILD_DATA)'" -tags "$(TAGS)" $(BUILDFLAGS) $(MODULE)/cmd/xkcd-gtk

dev: $(APPDATA)
	go build -o $(EXEC)-dev -ldflags="-X $(BUILD_PACKAGE).data=$(BUILD_DATA)" -tags "$(TAGS) xkcd_gtk_debug" $(DEVFLAGS) $(MODULE)/cmd/xkcd-gtk

$(POT): tools/fill-pot-header.sh $(GO_SOURCES) $(UI_SOURCES) $(DESKTOP).in $(APPDATA).in
	xgettext -o $@ --package-name="$(APP)" --language=C -kl $(POTFLAGS) $(GO_SOURCES)
	xgettext -o $@ --package-name="$(APP)" -j $(POTFLAGS) $(UI_SOURCES)
	xgettext -o $@ --package-name="$(APP)" -j --language=Desktop $(POTFLAGS) $(DESKTOP).in
	xgettext -o $@ --package-name="$(APP)" -j --its=po/appdata.its $(POTFLAGS) $(APPDATA).in
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

appcenter-reviews: # No dependencies, manifest will be read as-is from repo.
	flatpak-builder $(FPBFLAGS) --state-dir=flatpak-build/.flatpak-builder-$@/ --install-deps-from=appcenter flatpak-build/$@/ $(APP).yml

$(APP).yml: flatpak/appcenter.yml flatpak/modules.txt
	sed -e "s/path: '..'/path: '.'/" -e 's/path: modules.txt/path: flatpak\/modules.txt/' $< >$@

fix: $(POT) $(PO) $(APP).yml
	go fix $(MODULE_PACKAGES)
	go fmt $(MODULE_PACKAGES)
	go mod tidy
	([ -d vendor ] && go mod vendor) || true
	echo $(UI_SOURCES) | xargs -n1 gtk-builder-tool simplify --replace
	(find data -name '*.in' && \
	 printf '%s\n' $(UI_SOURCES) && \
	 grep -E -l '[^A-Za-z0-9]l\(' $(GO_SOURCES)\
	) | sort >po/POTFILES
	dos2unix -q po/LINGUAS po/appdata.its $(POT) $(PO)
	for lang in $(LINGUAS); do \
		msgmerge -U --backup=none "po/$$lang.po" $(POT); \
	done

check: $(APPDATA) $(FLATPAK_YML)
	go vet -tags "$(TAGS)" $(VETFLAGS) $(MODULE_PACKAGES)
	staticcheck -tags "$(TAGS)" $(MODULE_PACKAGES)
	shellcheck $(SH_SOURCES)
	xmllint --noout $(APPDATA) $(ICON) $(UI_SOURCES)
	yamllint .github/*.yml .github/*/*.yml $(FLATPAK_YML)
	-appstreamcli validate $(ASVFLAGS) --no-net $(APPDATA)
	-appstreamcli validate $(ASVFLAGS) --strict $(APPDATA)

test: $(FLATPAK_YML) $(APPDATA)
	go test -ldflags="-X $(BUILD_PACKAGE).data=$(BUILD_DATA)" -tags "$(TAGS)" $(TESTFLAGS) $(MODULE_PACKAGES)
	tools/test-banned-imports.sh
	tools/test-flatpak-config.sh $(FLATPAK_YML)
	tools/test-install.sh

# Shorthand for all the targets that CI covers.
ci: all check test
ci-full: ci appcenter appcenter-reviews flathub

clean:
	rm -f $(EXEC)
	rm -f $(EXEC)-dev
	rm -f $(DESKTOP)
	rm -f $(APPDATA)
	rm -f $(MO)
	rm -f flatpak/*.yml
	rm -rf flatpak-build/
	rm -rf .flatpak-builder/

strip: $(EXEC)
	strip $(EXEC)

install: $(EXEC) $(DESKTOP) $(APPDATA) $(MO)
	mkdir -p $(DESTDIR)$(bindir)
	install $(EXEC) $(DESTDIR)$(bindir)
	mkdir -p $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	cp $(ICON) $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	mkdir -p $(DESTDIR)$(datadir)/applications
	cp $(DESKTOP) $(DESTDIR)$(datadir)/applications
	mkdir -p $(DESTDIR)$(datadir)/dbus-1/services
	cp $(SERVICE) $(DESTDIR)$(datadir)/dbus-1/services
	mkdir -p $(DESTDIR)$(datadir)/metainfo
	cp $(APPDATA) $(DESTDIR)$(datadir)/metainfo
	for lang in $(LINGUAS); do \
		mkdir -p "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES"; \
		cp "po/$$lang.mo" "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES/$(APP).mo"; \
	done

uninstall:
	rm $(DESTDIR)$(bindir)/$(notdir $(EXEC))
	rm $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps/$(notdir $(ICON))
	rm $(DESTDIR)$(datadir)/applications/$(notdir $(DESKTOP))
	rm $(DESTDIR)$(datadir)/dbus-1/services/$(notdir $(SERVICE))
	rm $(DESTDIR)$(datadir)/metainfo/$(notdir $(APPDATA))
	for lang in $(LINGUAS); do \
		rm "$(DESTDIR)$(datadir)/locale/$$lang/LC_MESSAGES/$(APP).mo"; \
	done

.PHONY: all appcenter appcenter-install appcenter-reviews check ci ci-full clean dev fix flathub flathub-install install strip test uninstall
