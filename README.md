`xkcd-gtk` is a simple XKCD comic viewer written in go using GTK+3.

Install
-------

To install for current user:

	make && make home-install

To install for all users:

	make && sudo make root-install

Uninstall
---------

To uninstall for current user:

	make && make home-uninstall

To uninstall for all users:

	make && sudo make root-uninstall

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
