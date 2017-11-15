APP=com.github.rkoesters.xkcd-gtk
EXE=$(APP)
ICON=$(APP).svg
DESKTOP_FILE=$(APP).desktop
APPDATA_FILE=$(APP).appdata.xml

# BUILD VARIABLES
BUILDFLAGS=-tags $(shell tools/gtk-version.sh)
DEPS=$(shell tools/list-deps.sh)

# INSTALL VARIABLES
prefix=/usr
bindir=$(prefix)/bin
datadir=$(prefix)/share

# If GOPATH isn't set, then just use the current directory.
ifeq "$(shell go env GOPATH)" ""
export GOPATH=$(shell pwd)
endif

all: deps $(EXE)

deps:
	go get -u $(BUILDFLAGS) $(DEPS)

$(EXE): *.go
	go build $(BUILDFLAGS) -o $@

clean:
	rm -f $(EXE)

install: $(EXE)
	mkdir -p $(DESTDIR)$(bindir)
	install $(EXE) $(DESTDIR)$(bindir)
	mkdir -p $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	cp data/$(ICON) $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps
	mkdir -p $(DESTDIR)$(datadir)/applications
	cp data/$(DESKTOP_FILE) $(DESTDIR)$(datadir)/applications
	mkdir -p $(DESTDIR)$(datadir)/metainfo
	cp data/$(APPDATA_FILE) $(DESTDIR)$(datadir)/metainfo

uninstall:
	rm $(DESTDIR)$(bindir)/$(EXE) \
	   $(DESTDIR)$(datadir)/icons/hicolor/scalable/apps/$(ICON) \
	   $(DESTDIR)$(datadir)/applications/$(DESKTOP_FILE) \
	   $(DESTDIR)$(datadir)/metainfo/$(APPDATA_FILE)
