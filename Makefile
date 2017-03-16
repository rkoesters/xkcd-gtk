BUILDFLAGS=-tags gtk_3_18

xkcd-gtk: *.go
	go build $(BUILDFLAGS)

deps:
	go get -u $(BUILDFLAGS) github.com/rkoesters/xkcd
	go get -u $(BUILDFLAGS) github.com/rkoesters/xdg/...
	go get -u $(BUILDFLAGS) github.com/blevesearch/bleve/...
	go get -u $(BUILDFLAGS) github.com/gotk3/gotk3/...

clean:
	go clean

home-install: xkcd-gtk
	mkdir -p $$HOME/bin
	install xkcd-gtk $$HOME/bin
	mkdir -p $$HOME/.local/share/applications
	cp com.ryankoesters.xkcd-gtk.desktop $$HOME/.local/share/applications
	mkdir -p $$HOME/.local/share/icons/hicolor/scalable/apps
	cp xkcd-gtk.svg $$HOME/.local/share/icons/hicolor/scalable/apps

home-uninstall:
	rm $$HOME/bin/xkcd-gtk \
	   $$HOME/.local/share/applications/com.ryankoesters.xkcd-gtk.desktop \
	   $$HOME/.local/share/icons/hicolor/scalable/apps/xkcd-gtk.svg

root-install: xkcd-gtk
	mkdir -p /usr/local/bin
	install xkcd-gtk /usr/local/bin
	mkdir -p /usr/share/applications
	cp com.ryankoesters.xkcd-gtk.desktop /usr/share/applications
	mkdir -p /usr/local/share/icons/hicolor/scalable/apps
	cp xkcd-gtk.svg /usr/local/share/icons/hicolor/scalable/apps

root-uninstall:
	rm /usr/local/bin/xkcd-gtk \
	   /usr/share/applications/com.ryankoesters.xkcd-gtk.desktop \
	   /usr/local/share/icons/hicolor/scalable/apps/xkcd-gtk.svg
