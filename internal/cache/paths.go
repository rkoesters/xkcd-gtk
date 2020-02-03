package cache

import (
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"path/filepath"
	"strconv"
)

func cacheVersionPath() string {
	return filepath.Join(paths.CacheDir(), "cache_version")
}

func comicCacheDBPath() string {
	return filepath.Join(paths.CacheDir(), "comics")
}

func comicImageDirPath() string {
	return filepath.Join(paths.CacheDir(), "comic_image")
}

// ComicImagePath returns the path to the specified comic image within the
// cache. The file at the returned path may or may not exist. If it does not
// exist, call DownloadComicImage to fetch the file from the internet.
func ComicImagePath(n int) string {
	return filepath.Join(comicImageDirPath(), strconv.Itoa(n))
}
