// Package cache provides a cached interface to the xkcd server.
package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	bolt "github.com/etcd-io/bbolt"
	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"io"
	"io/ioutil"
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
)

// Init initializes the comic cache. Function index is called each time a comic
// is inserted into the comic cache.
func Init(index func(comic *xkcd.Comic) error) error {
	addToSearchIndex = index

	err := os.MkdirAll(paths.CacheDir(), 0755)
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
			SafeTitle: l("Comic Not Found"),
			Title:     l("Comic Not Found"),
		}, xkcd.ErrNotFound
	}

	// First, check if we have the file.
	err := cacheDB.View(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			comic = &xkcd.Comic{
				Num:       n,
				SafeTitle: l("Error trying to access metadata cache"),
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
			comic = &xkcd.Comic{
				Num:       n,
				SafeTitle: l("Error parsing comic metadata from cache"),
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
				SafeTitle: l("Comic Not Found"),
			}, err
		} else if err != nil {
			return &xkcd.Comic{
				Num:       n,
				SafeTitle: l("Couldn't Get Comic"),
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

	newest := <-recvCachedNewestComic

	if newest == nil {
		newest, err = xkcd.GetCurrent()
		if err != nil {
			newest, err = newestComicInfoFromCache()
			if err != nil {
				return newest, err
			}
			return newest, ErrOffline
		}

		sendCachedNewestComic <- newest
	}
	return newest, nil
}

// NewestComicInfoSkipCache is equivalent to NewestComicInfo except it always
// queries the internet to check for a new comic.
func NewestComicInfoSkipCache() (*xkcd.Comic, error) {
	sendCachedNewestComic <- nil
	return NewestComicInfo()
}

// NewestComicInfoAsync always returns a valid *xkcd.Comic that appears to be
// newest, and err will be set if any errors were encountered, however these
// errors can be ignored safely. This function will return the newest comic info
// based on the cache, but then asynchronously checks for the newest comic from
// the internet and calls callback when the asynchronous call completes.
func NewestComicInfoAsync(callback func(*xkcd.Comic, error)) (*xkcd.Comic, error) {
	newest := <-recvCachedNewestComic

	if newest == nil {
		newestFromCache, err := newestComicInfoFromCache()

		go func() {
			newestFromInternet, err := xkcd.GetCurrent()

			if newestFromInternet != nil && err == nil {
				sendCachedNewestComic <- newestFromInternet
			} else {
				sendCachedNewestComic <- newestFromCache
			}

			callback(<-recvCachedNewestComic, err)
		}()

		return newestFromCache, err
	}
	return newest, nil
}

func newestComicInfoFromCache() (*xkcd.Comic, error) {
	newest := &xkcd.Comic{
		SafeTitle: l("Connect to the internet to download some comics!"),
	}

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

	return newest, err
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
