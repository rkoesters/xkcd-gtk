prefix=/usr/local
bindir=$(prefix)/bin
desktopdir=$(prefix)/share/applications
icondir=$(prefix)/share/icons/hicolor/scalable/apps

xkcd-gtk:
	go build

clean:
	go clean

install: xkcd-gtk
	install xkcd-gtk $(prefix)/bin
	mkdir -p $(desktopdir)
	cp xkcd-gtk.desktop $(desktopdir)
	mkdir -p $(icondir)
	cp xkcd-gtk.svg $(icondir)

uninstall:
	rm $(bindir)/xkcd-gtk \
	   $(desktopdir)/xkcd-gtk.desktop \
	   $(icondir)/xkcd-gtk.svg
