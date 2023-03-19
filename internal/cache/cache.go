// Package cache provides a cached interface to the xkcd server.
package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	bolt "go.etcd.io/bbolt"
)

var (
	offlineMode = flag.Bool("offline", false, "Do not use network, rely only on local data.")
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

	recvCachedNewestComic          <-chan *xkcd.Comic
	sendCachedNewestComic          chan<- *xkcd.Comic
	recvCachedNewestComicUpdatedAt <-chan time.Time

	// addToSearchIndex is a callback to insert the given comic into the search
	// index.
	addToSearchIndex func(comic *xkcd.Comic) error

	// Error messages to be shown in the window title. Initialized in Init to
	// provide translations to system language.
	cacheDatabaseError    string
	comicNotFound         string
	couldNotDownloadComic string
	noComicsFound         string
)

// Init initializes the comic cache. Function index is called each time a comic
// is inserted into the comic cache.
func Init(index func(comic *xkcd.Comic) error) error {
	checkForMisplacedCacheFiles()

	addToSearchIndex = index

	// Initialize localized error strings.
	cacheDatabaseError = l("Error reading local comic database")
	comicNotFound = l("Comic not found")
	couldNotDownloadComic = l("Couldn't get comic")
	noComicsFound = l("Connect to the internet to download some comics!")

	err := paths.EnsureCacheDir()
	if err != nil {
		return err
	}

	// If the user's cache isn't compatible with our binary's cache
	// implementation, then we need to start over (we will move the old cache to
	// .bak just in case).
	if existingCacheVersion() != currentCacheVersion() {
		log.Debug("incompatible cache database found, backing up and rebuilding cache database...")
		os.Rename(comicCacheDBPath(), comicCacheDBPath()+".bak")
	}

	// Open comic cache database.
	log.Debug("openning cache database: ", comicCacheDBPath())
	cacheDB, err = bolt.Open(comicCacheDBPath(), 0644, nil)
	if err != nil {
		return err
	}

	// Create comic cache buckets, if they do not exist.
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

	// Create comic image cache directory, if it does not exist.
	err = os.MkdirAll(filepath.Join(comicImageDirPath()), 0755)
	if err != nil {
		return err
	}

	cachedNewestComicOut := make(chan *xkcd.Comic)
	cachedNewestComicIn := make(chan *xkcd.Comic)
	cachedNewestComicUpdatedAtOut := make(chan time.Time)

	recvCachedNewestComic = cachedNewestComicOut
	sendCachedNewestComic = cachedNewestComicIn
	recvCachedNewestComicUpdatedAt = cachedNewestComicUpdatedAtOut

	// Start cachedNewestComic manager.
	go func() {
		var (
			cachedNewestComic          *xkcd.Comic
			cachedNewestComicUpdatedAt time.Time
		)

		for {
			select {
			case newest := <-cachedNewestComicIn:
				cachedNewestComicUpdatedAt = time.Now()
				cachedNewestComic = newest
				log.Debugf("newest cached comic set to %v at %v", newest.Num, cachedNewestComicUpdatedAt)
			case cachedNewestComicOut <- cachedNewestComic:
				// Sending the comic was all we wanted to do.
			case cachedNewestComicUpdatedAtOut <- cachedNewestComicUpdatedAt:
				// Sending the time stamp was all we wanted to do.
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

// DownloadAllComicMetadata asynchronously fills the comic metadata cache and
// search index via the internet. Status can be checked with Stat().
func DownloadAllComicMetadata(cacheWindow ViewRefresherGetter) {
	// Make sure all comic metadata is cached and indexed.
	go func() {
		newest, err := NewestComicInfoFromInternet()
		if err != nil {
			log.Print(err)
			return
		}
		for i := 1; i <= newest.Num; i++ {
			ComicInfo(i)
			cacheWindow().RefreshMetadataWith(Stat{
				LatestComicNumber: newest.Num,
				CachedCount:       i,
			})
		}
	}()
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
			// The comic metadata isn't in our cache yet, we will try to
			// download it.
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

// CheckForNewestComicInfo returns the latest xkcd.Comic. May query xkcd API if
// latest comic in the cache has not been updated since freshnessThreshold. The
// returned error can be safely ignored. Should not be used on UI event loop
// (TODO: stop breaking this rule).
func CheckForNewestComicInfo(freshnessThreshold time.Duration) (*xkcd.Comic, error) {
	if time.Since(<-recvCachedNewestComicUpdatedAt) < freshnessThreshold {
		return NewestComicInfoFromCache()
	}

	c, err := NewestComicInfoFromInternet()
	if err != nil {
		return NewestComicInfoFromCache()
	}
	return c, nil
}

// NewestComicInfoFromCache returns the newest comic info available in the
// cache. The function will not use the internet. The returned error can be
// safely ignored. Can be used on UI event loop.
func NewestComicInfoFromCache() (*xkcd.Comic, error) {
	// Check in-memory cache.
	newest := <-recvCachedNewestComic
	if newest != nil {
		return newest, nil
	}

	// Check on-disk cache.
	newest = &xkcd.Comic{}
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
// May return nil, the returned error should be checked. Should not be used on
// UI event loop.
func NewestComicInfoFromInternet() (*xkcd.Comic, error) {
	if *offlineMode {
		return nil, ErrOffline
	}

	log.Debug("NewestComicInfoFromInternet start")
	defer log.Debug("NewestComicInfoFromInternet end")

	c, err := xkcd.GetCurrent()
	if err != nil {
		return nil, ErrOffline
	}

	sendCachedNewestComic <- c
	return c, putComicInfo(c)
}

func downloadComicInfo(n int) (*xkcd.Comic, error) {
	if *offlineMode {
		return nil, ErrOffline
	}

	log.Debugf("downloadComicInfo(%v) start", n)
	defer log.Debugf("downloadComicInfo(%v) end", n)

	comic, err := xkcd.Get(n)
	if err != nil {
		return nil, err
	}
	return comic, putComicInfo(comic)
}

// putComicInfo adds the given xkcd.Comic to the cache database.
func putComicInfo(comic *xkcd.Comic) error {
	err := cacheDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			return ErrLocalFailure
		}

		var buf bytes.Buffer
		e := json.NewEncoder(&buf)
		err := e.Encode(comic)
		if err != nil {
			return err
		}

		return bucket.Put(intToBytes(comic.Num), buf.Bytes())
	})
	if err != nil {
		return err
	}

	return addToSearchIndex(comic)
}

// DownloadComicImage tries to add a comic image to our local cache. If
// successful, the image can be found at the path returned by ComicImagePath.
func DownloadComicImage(n int, cacheWindow ViewRefresherGetter) error {
	if *offlineMode {
		return ErrOffline
	}

	defer func() { go cacheWindow().RefreshImages() }()

	log.Debugf("DownloadComicImage(%v) start", n)
	defer log.Debugf("DownloadComicImage(%v) end", n)

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

// DownloadAllComicImages tries to add all comic images to our local cache. If
// successful, the images can be found at the path returned by ComicImagePath.
func DownloadAllComicImages(cacheWindow ViewRefresherGetter) {
	newest, err := NewestComicInfoFromCache()
	if err != nil {
		log.Print(err)
		return
	}
	for i := 1; i <= newest.Num; i++ {
		_, err = os.Stat(ComicImagePath(i))
		if os.IsNotExist(err) {
			DownloadComicImage(i, newNullRefresher)
		}

		cacheWindow().RefreshImagesWith(Stat{
			LatestComicNumber: newest.Num,
			CachedCount:       i,
		})
	}
}

// currentCacheVersion returns the cache version for this binary.
func currentCacheVersion() int { return cacheVersionCurrent }

// existingCacheVersion returns the cache version for the user's existing cache.
func existingCacheVersion() int {
	b, err := os.ReadFile(cacheVersionPath())
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
