#!/bin/sh -eu
exe="com.github.rkoesters.xkcd-gtk.exe"
icon="data/com.github.rkoesters.xkcd-gtk.svg"
app_version=$(tools/app-version.sh)

tmp_target="/tmp/xkcd-gtk-$app_version"
archive_target="xkcd-gtk-$app_version-windows-$MSYSTEM_CARCH.zip"

echo "Building executable..."
make clean EXE_PATH="$exe"
make EXE_PATH="$exe" LDFLAGS="-ldflags='-H=windowsgui -X main.appVersion=$app_version'"

mkdir -p "$tmp_target"

echo "Copying executable..."
install -t "$tmp_target" "$exe"

echo "Copying DLLs..."
ldd "$exe" |
grep "=>" |
grep "$MINGW_PREFIX" |
cut -d ' ' -f 3 |
sort |
uniq |
xargs install -t "$tmp_target"

echo "Copying gdk-pixbuf loaders..."
mkdir -p "$tmp_target/lib/gdk-pixbuf-2.0/2.10.0/loaders"
install -t "$tmp_target/lib/gdk-pixbuf-2.0/2.10.0/loaders" "$MINGW_PREFIX/lib/gdk-pixbuf-2.0/2.10.0/loaders"/*.dll
cp "$MINGW_PREFIX/lib/gdk-pixbuf-2.0/2.10.0/loaders.cache" "$tmp_target/lib/gdk-pixbuf-2.0/2.10.0"

echo "Copying gdk-pixbuf DLLs..."
ldd "$MINGW_PREFIX/lib/gdk-pixbuf-2.0/2.10.0/loaders"/*.dll |
grep "=>" |
grep "$MINGW_PREFIX" |
cut -d ' ' -f 3 |
sort |
uniq |
xargs install -t "$tmp_target"

echo "Copying icons..."
mkdir -p "$tmp_target/share/icons/"
cp -R "$MINGW_PREFIX/share/icons/hicolor" "$tmp_target/share/icons/"
cp "$icon" "$tmp_target/share/icons/hicolor/scalable/apps/"
gtk-update-icon-cache-3.0 -f "$tmp_target/share/icons/hicolor"
cp -R "$MINGW_PREFIX/share/icons/Adwaita" "$tmp_target/share/icons/"
gtk-update-icon-cache-3.0 -f "$tmp_target/share/icons/Adwaita"

echo "Creating gtk settings..."
mkdir -p "$tmp_target/etc/gtk-3.0"
{
	echo "[Settings]"
	echo "gtk-theme-name = win32"
	echo "gtk-icon-theme-name = Adwaita"
	echo ""
	echo "gtk-dialogs-use-header = true"
	echo ""
	echo "gtk-xft-antialias = 1"
	echo "gtk-xft-hinting = 1"
	echo "gtk-xft-hintstyle = hintfull"
	echo "gtk-xft-rgba = rgb"
} >"$tmp_target/etc/gtk-3.0/settings.ini"

echo "Creating zip package..."
rm -f "$archive_target"
output=$(readlink -f $archive_target)
(
	dirname=$(dirname "$tmp_target")
	basename=$(basename "$tmp_target")
	cd "$dirname"
	zip -r "$output" "$basename"
)
rm -r "$tmp_target"
