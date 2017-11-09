APP=com.github.rkoesters.xkcd-gtk
EXE=$(APP)
ICON=$(APP).svg
DESKTOP_FILE=$(APP).desktop
APPDATA_FILE=$(APP).appdata.xml

BUILDFLAGS=-tags gtk_3_18

prefix=$(DESTDIR)/usr
bindir=$(prefix)/bin
datadir=$(prefix)/share

all: deps $(EXE)

deps:
	go get -u $(BUILDFLAGS) github.com/rkoesters/xkcd
	go get -u $(BUILDFLAGS) github.com/rkoesters/xdg/...
	go get -u $(BUILDFLAGS) github.com/skratchdot/open-golang/open
	go get -u $(BUILDFLAGS) github.com/blevesearch/bleve/...
	go get -u $(BUILDFLAGS) github.com/gotk3/gotk3/...

$(EXE): *.go
	go build $(BUILDFLAGS) -o $@

clean:
	rm -f $(EXE)

install: $(EXE)
	mkdir -p $(bindir)
	install $(EXE) $(bindir)
	mkdir -p $(datadir)/icons/hicolor/scalable/apps
	cp $(ICON) $(datadir)/icons/hicolor/scalable/apps
	mkdir -p $(datadir)/applications
	cp $(DESKTOP_FILE) $(datadir)/applications
	mkdir -p $(datadir)/metainfo
	cp $(APPDATA_FILE) $(datadir)/metainfo

uninstall:
	rm $(bindir)/$(EXE) \
	   $(datadir)/icons/hicolor/scalable/apps/$(ICON) \
	   $(datadir)/applications/$(DESKTOP_FILE) \
	   $(datadir)/metainfo/$(APPDATA_FILE)
