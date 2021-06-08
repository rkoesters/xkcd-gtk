// Package cache provides a cached interface to the xkcd server.
package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	bolt "go.etcd.io/bbolt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// cacheVersionCurrent should be incremented every time a release breaks
	// compatibility with the previous release's cache (although breaking
	// compatibility should be avoided).
	cacheVersionCurrent = 2
)

var (
	cacheDB *bolt.DB

	comicCacheMetadataBucketName = []byte("comic_metadata")
	comicCacheImageBucketName    = []byte("comic_image")

	recvCachedNewestComic <-chan *xkcd.Comic
	sendCachedNewestComic chan<- *xkcd.Comic

	// addToSearchIndex is a callback to insert the given comic into the
	// search index.
	addToSearchIndex func(comic *xkcd.Comic) error

	// Error messages to be shown in the window title. Initialized in Init
	// to provide translations to system language.
	cacheDatabaseError    string
	comicNotFound         string
	couldNotDownloadComic string
	noComicsFound         string
)

// Init initializes the comic cache. Function index is called each time a comic
// is inserted into the comic cache.
func Init(index func(comic *xkcd.Comic) error) error {
	addToSearchIndex = index

	cacheDatabaseError = l("Error reading local comic database")
	comicNotFound = l("Comic Not Found")
	couldNotDownloadComic = l("Couldn't Get Comic")
	noComicsFound = l("Connect to the internet to download some comics!")

	err := paths.EnsureCacheDir()
	if err != nil {
		return err
	}

	// If the user's cache isn't compatible with our binary's cache
	// implementation, then we need to start over (we will move the old
	// cache to .bak just in case).
	if existingCacheVersion() != currentCacheVersion() {
		os.Rename(comicCacheDBPath(), comicCacheDBPath()+".bak")
	}

	cacheDB, err = bolt.Open(comicCacheDBPath(), 0644, nil)
	if err != nil {
		return err
	}

	err = cacheDB.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(comicCacheMetadataBucketName)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(comicCacheImageBucketName)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(comicImageDirPath()), 0755)
	if err != nil {
		return err
	}

	cachedNewestComicOut := make(chan *xkcd.Comic)
	cachedNewestComicIn := make(chan *xkcd.Comic)

	recvCachedNewestComic = cachedNewestComicOut
	sendCachedNewestComic = cachedNewestComicIn

	// Start cachedNewestComic manager.
	go func() {
		var cachedNewestComic *xkcd.Comic

		for {
			select {
			case newest := <-cachedNewestComicIn:
				cachedNewestComic = newest
			case cachedNewestComicOut <- cachedNewestComic:
				// Sending the comic was all we wanted to do.
			}
		}
	}()

	return nil
}

// Close closes the comic cache.
func Close() error {
	err := cacheDB.Close()
	if err != nil {
		return err
	}

	f, err := os.Create(cacheVersionPath())
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, currentCacheVersion())
	return err
}

// ComicInfo always returns a valid *xkcd.Comic that can be used, and err will
// be set if any errors were encountered, however these errors can be ignored
// safely.
func ComicInfo(n int) (*xkcd.Comic, error) {
	var comic *xkcd.Comic

	// Don't bother asking the server for comic 404, it will always return a
	// 404 error.
	if n == 404 {
		return &xkcd.Comic{
			Num:       n,
			SafeTitle: comicNotFound,
			Title:     comicNotFound,
		}, xkcd.ErrNotFound
	}

	// First, check if we have the file.
	err := cacheDB.View(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			log.Print("error trying to access metadata cache")
			comic = &xkcd.Comic{
				Num:       n,
				SafeTitle: cacheDatabaseError,
			}
			return ErrLocalFailure
		}

		data := bucket.Get(intToBytes(n))
		if data == nil {
			// The comic metadata isn't in our cache yet, we will
			// try to download it.
			return ErrMiss
		}

		comic, err = xkcd.New(bytes.NewReader(data))
		if err != nil {
			log.Print("error parsing comic metadata from cache: ", err)
			comic = &xkcd.Comic{
				Num:       n,
				SafeTitle: cacheDatabaseError,
			}
			return err
		}

		return nil
	})
	if err == ErrMiss {
		comic, err = downloadComicInfo(n)
		if err == xkcd.ErrNotFound {
			return &xkcd.Comic{
				Num:       n,
				SafeTitle: comicNotFound,
			}, err
		} else if err != nil {
			return &xkcd.Comic{
				Num:       n,
				SafeTitle: couldNotDownloadComic,
			}, err
		}
	}

	return comic, err
}

// NewestComicInfo always returns a valid *xkcd.Comic that appears to be newest,
// and err will be set if any errors were encountered, however these errors can
// be ignored safely.
func NewestComicInfo() (*xkcd.Comic, error) {
	var err error

	// Check in-memory cache.
	newest := <-recvCachedNewestComic
	if newest != nil {
		return newest, nil
	}

	// Check on-disk cache.
	newest, err = NewestComicInfoFromCache()
	if err == nil {
		return newest, nil
	}

	// Check internet.
	newest, err = NewestComicInfoFromInternet()
	if err == nil {
		return newest, nil
	}

	return &xkcd.Comic{
		Num:       1,
		SafeTitle: noComicsFound,
	}, ErrNoComicsFound
}

// CheckForNewestComicInfo fetches the latest comic info from the internet. If
// it can not connect, then it fetches the latest comic from the cache. The
// returned error can be safely ignored.
func CheckForNewestComicInfo() (*xkcd.Comic, error) {
	// Check internet.
	c, err := NewestComicInfoFromInternet()
	if err == nil {
		return c, nil
	}

	// Check in-memory cache.
	newest := <-recvCachedNewestComic
	if newest != nil {
		return newest, nil
	}

	// Check on-disk cache.
	newest, err = NewestComicInfoFromCache()
	if err == nil {
		return newest, nil
	}

	return &xkcd.Comic{
		Num:       1,
		SafeTitle: noComicsFound,
	}, ErrNoComicsFound
}

// NewestComicInfoFromCache returns the newest comic info available in the
// cache. The function will not use the internet. The returned error can be
// safely ignored.
func NewestComicInfoFromCache() (*xkcd.Comic, error) {
	newest := &xkcd.Comic{}

	err := cacheDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			return ErrLocalFailure
		}

		return bucket.ForEach(func(k, v []byte) error {
			comic, err := xkcd.New(bytes.NewReader(v))
			if err != nil {
				return err
			}

			if comic.Num > newest.Num {
				newest = comic
			}

			return nil
		})
	})
	if err != nil {
		log.Print("error reading comic from cache: ", err)
	}
	if newest.Num <= 0 {
		return &xkcd.Comic{
			Num:       1,
			SafeTitle: noComicsFound,
		}, ErrNoComicsFound
	}

	sendCachedNewestComic <- newest
	return newest, nil
}

// NewestComicInfoFromInternet fetches the latest comic info from the internet.
// May return nil, the returned error should be checked.
func NewestComicInfoFromInternet() (*xkcd.Comic, error) {
	c, err := xkcd.GetCurrent()
	if err != nil {
		return nil, ErrOffline
	}

	sendCachedNewestComic <- c
	return c, nil
}

func downloadComicInfo(n int) (*xkcd.Comic, error) {
	comic, err := xkcd.Get(n)
	if err != nil {
		return nil, err
	}

	err = cacheDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			return ErrLocalFailure
		}

		var buf bytes.Buffer
		e := json.NewEncoder(&buf)
		err = e.Encode(comic)
		if err != nil {
			return err
		}

		return bucket.Put(intToBytes(n), buf.Bytes())
	})
	if err != nil {
		return nil, err
	}

	err = addToSearchIndex(comic)

	return comic, err
}

// DownloadComicImage tries to add a comic image to our local cache. If
// successful, the image can be found at the path returned by ComicImagePath.
func DownloadComicImage(n int) error {
	comic, err := ComicInfo(n)
	if err != nil {
		return err
	}

	resp, err := http.Get(comic.Img)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(ComicImagePath(n))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// currentCacheVersion returns the cache version for this binary.
func currentCacheVersion() int { return cacheVersionCurrent }

// existingCacheVersion returns the cache version for the user's existing cache.
func existingCacheVersion() int {
	b, err := ioutil.ReadFile(cacheVersionPath())
	if err != nil {
		return 0
	}

	num, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		return 0
	}
	return num
}

func intToBytes(i int) []byte {
	buf := make([]byte, binary.MaxVarintLen64)

	n := binary.PutVarint(buf, int64(i))

	return buf[:n]
}
