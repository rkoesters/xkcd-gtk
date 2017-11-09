APP=com.github.rkoesters.xkcd-gtk
EXE=$(APP)
DESKTOP_FILE=$(APP).desktop
ICON=$(APP).svg

BUILDFLAGS=-tags gtk_3_18

prefix=$(DESTDIR)/usr/local
bindir=$(prefix)/bin
datadir=$(prefix)/share

all: deps $(EXE)

deps:
	go get -u $(BUILDFLAGS) github.com/golang/lint/golint
	go get -u $(BUILDFLAGS) github.com/rkoesters/xkcd
	go get -u $(BUILDFLAGS) github.com/rkoesters/xdg/...
	go get -u $(BUILDFLAGS) github.com/skratchdot/open-golang/open
	go get -u $(BUILDFLAGS) github.com/blevesearch/bleve/...
	go get -u $(BUILDFLAGS) github.com/gotk3/gotk3/...

$(EXE): *.go
	go build $(BUILDFLAGS) -o $@

clean:
	go clean
	rm -f $(EXE)

fmt:
	go fmt

lint:
	golint

install: $(EXE)
	mkdir -p $(bindir)
	install $(EXE) $(bindir)
	mkdir -p $(datadir)/applications
	cp $(DESKTOP_FILE) $(datadir)/applications
	mkdir -p $(datadir)/icons/hicolor/scalable/apps
	cp $(ICON) $(datadir)/icons/hicolor/scalable/apps

uninstall:
	rm $(bindir)/$(EXE) \
	   $(datadir)/applications/$(DESKTOP_FILE) \
	   $(datadir)/icons/hicolor/scalable/apps/$(ICON)
