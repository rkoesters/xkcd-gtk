#!/bin/sh -eu
# Location we are putting files.
target="$HOME/xkcd-gtk-win"

exe="com.github.rkoesters.xkcd-gtk.exe"
icon="data/com.github.rkoesters.xkcd-gtk.svg"

make clean
make LDFLAGS="-ldflags='-H=windowsgui -X main.appVersion=$(tools/app-version.sh)'"

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
