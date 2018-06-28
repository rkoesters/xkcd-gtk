# Comic Sticks

`xkcd-gtk` is a simple xkcd comic viewer written in Go using GTK+3.

![screenshot](screenshots/screenshot-1.png)

## Requirements

To build this program, you will need [Go](https://golang.org/) and GTK+3
development files (something like `libgtk-3-dev`).

## Installing

### On elementaryOS

[![Get it on AppCenter](https://appcenter.elementary.io/badge.svg)](https://appcenter.elementary.io/com.github.rkoesters.xkcd-gtk)

### On Debian-based distros

Download the .deb file from
[releases](https://github.com/rkoesters/xkcd-gtk/releases) and install
it using your preferred installation tool.

### From source

To install for current user:

	make deps
	make prefix=$HOME/.local install

To uninstall for current user:

	make prefix=$HOME/.local uninstall

To install for all users:

	make deps
	make
	sudo make install

To uninstall for all users:

	sudo make uninstall

## License

This program comes with absolutely no warranty. See the [GNU General
Public License, version 3 or later](LICENSE) for details.
