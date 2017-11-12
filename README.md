Comic Sticks
============

`xkcd-gtk` is a simple xkcd comic viewer written in Go using GTK+3.

[![Get it on AppCenter](https://appcenter.elementary.io/badge.svg)](https://appcenter.elementary.io/com.github.rkoesters.xkcd-gtk)


![screenshot](screenshots/screenshot-1.png)

Requirements
------------

To build this program, you will need [Go](https://golang.org/) and GTK+3
development files (something like `libgtk-3-dev`).

Install
-------

To install for current user:

	make prefix=$HOME/.local install

To install for all users:

	make && sudo make install

Uninstall
---------

To uninstall for current user:

	make prefix=$HOME/.local uninstall

To uninstall for all users:

	sudo make uninstall

License
-------

This program comes with absolutely no warranty. See the [GNU General
Public License, version 3 or later](LICENSE) for details.
