#!/bin/sh -e
# Go package directory.
godir="$GOPATH/src/github.com/rkoesters/xkcd-gtk"
# Location we are putting files.
target="$HOME/xkcd-gtk-win"

exe="$godir/xkcd-gtk.exe"
icon="$godir/xkcd-gtk.svg"

cd "$godir"

make clean
make BUILDFLAGS="-ldflags -H=windowsgui -tags gtk_3_18" xkcd-gtk

# Made sure our target exists.
mkdir -p "$target"

echo "Copying xkcd-gtk.exe to $target"
install -t "$target" "$exe"

echo "Copying DLLs to $target"
ldd "$exe" | grep mingw64 | cut -d ' ' -f 3 | xargs install -t "$target"

echo "Copying icons to $target"
mkdir -p "$target/share/icons/"
cp -R "/mingw64/share/icons/Adwaita" "$target/share/icons/"
cp "$icon" "$target/share/icons/Adwaita/scalable/apps/"

echo "Creating settings.ini in $target"
mkdir -p "$target/etc/gtk-3.0"
printf "[Settings]\ngtk-theme-name=win32" >"$target/etc/gtk-3.0/settings.ini"
