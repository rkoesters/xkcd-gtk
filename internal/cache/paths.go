package cache

import (
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func checkForMisplacedCacheFiles() {
	// map[misplacedFile]correctPath
	misplacedCacheFiles := make(map[string]string)
	misplacedCacheFiles[filepath.Join(paths.Builder{}.CacheDir(), "cache_version")] = cacheVersionPath()
	misplacedCacheFiles[filepath.Join(paths.Builder{}.CacheDir(), "comics")] = comicCacheDBPath()
	misplacedCacheFiles[filepath.Join(paths.Builder{}.CacheDir(), "comic_image")] = comicImageDirPath()

	for misplaced, correct := range misplacedCacheFiles {
		_, err := os.Stat(misplaced)
		if !os.IsNotExist(err) {
			log.Printf("WARNING: Potentially misplaced cache file '%v'. Should be '%v'.", misplaced, correct)
		}
	}
}

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
