xkcd-gtk:
	go build

clean:
	go clean

install:
	go install
	cp xkcd-gtk.desktop ~/.local/share/applications/
	cp xkcd-gtk.svg ~/.local/share/icons/hicolor/scalable/apps/
