prefix=/usr/local
bindir=$(prefix)/bin
desktopdir=$(prefix)/share/applications
icondir=$(prefix)/share/icons/hicolor/scalable/apps

xkcd-gtk: *.go
	go build

clean:
	go clean
	rm -f com.ryankoesters.xkcd-gtk.desktop

install: xkcd-gtk
	install xkcd-gtk $(prefix)/bin
	mkdir -p $(desktopdir)
	sed -e 's#BINDIR#$(bindir)#g' -e 's#ICONDIR#$(icondir)#g' <com.ryankoesters.xkcd-gtk.desktop.in >$(desktopdir)/com.ryankoesters.xkcd-gtk.desktop
	mkdir -p $(icondir)
	cp xkcd-gtk.svg $(icondir)

uninstall:
	rm $(bindir)/xkcd-gtk \
	   $(desktopdir)/com.ryankoesters.xkcd-gtk.desktop \
	   $(icondir)/xkcd-gtk.svg
