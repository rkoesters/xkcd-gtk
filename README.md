![xkcd-gtk icon](https://cdn.rawgit.com/rkoesters/xkcd-gtk/master/xkcd-gtk.svg)

`xkcd-gtk` is a simple XKCD comic viewer written in go using GTK+3.

Requirements
------------

To build this program, you will need [Go](https://golang.org/) and GTK+3
developement files (something like `libgtk-3-dev`).

Install
-------

To install for current user:

	make deps && make && make home-install

To install for all users:

	make deps && make && sudo make root-install

Uninstall
---------

To uninstall for current user:

	make home-uninstall

To uninstall for all users:

	sudo make root-uninstall

Todo
----

- Add search.
  - Will probably use blevesearch
- Save state: it seems our gtk wrapper doesn't support
  win.IsMaximized(), so we can't remember if we were maximized.

License
-------

This program comes with absolutely no warranty. See the [GNU General
Public License, version 3 or later](LICENSE) for details.
