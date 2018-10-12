package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/rkoesters/xkcd"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	comicCacheName         = "comics"
	comicCacheMetadataName = "comic_metadata"
	comicCacheImageName    = "comic_image"
)

var (
	// ErrCacheMiss means that value we are looking for wasn't in
	// the cache.
	ErrCacheMiss = errors.New("cache miss")
	// ErrCache means that there was an error while trying to access
	// the local cache.
	ErrCache = errors.New("error accessing local xkcd cache")
	// ErrOffline means that there was an error trying to access the
	// xkcd server.
	ErrOffline = errors.New("error accessing xkcd server")

	cacheDB *bolt.DB

	comicCacheMetadataBucketName = []byte(comicCacheMetadataName)
	comicCacheImageBucketName    = []byte(comicCacheImageName)

	getCachedNewestComic <-chan *xkcd.Comic
	setCachedNewestComic chan<- *xkcd.Comic
)

func initComicCache() error {
	err := os.MkdirAll(CacheDir(), 0755)
	if err != nil {
		return err
	}

	cacheDB, err = bolt.Open(getComicCachePath(), 0644, nil)
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

	err = os.MkdirAll(filepath.Join(CacheDir(), comicCacheImageName), 0755)
	if err != nil {
		return err
	}

	cachedNewestComicOut := make(chan *xkcd.Comic)
	cachedNewestComicIn := make(chan *xkcd.Comic)

	getCachedNewestComic = cachedNewestComicOut
	setCachedNewestComic = cachedNewestComicIn

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

func closeComicCache() error {
	return cacheDB.Close()
}

// GetComicInfo always returns a valid *xkcd.Comic that can be used, and
// err will be set if any errors were encountered, however these errors
// can be ignored safely.
func GetComicInfo(n int) (*xkcd.Comic, error) {
	var c *xkcd.Comic

	// First, check if we have the file.
	err := cacheDB.View(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			c = &xkcd.Comic{
				Num:       n,
				SafeTitle: "Error trying to access metadata cache",
			}
			return ErrCache
		}

		data := bucket.Get(intToBytes(n))
		if data == nil {
			// The comic metadata isn't in our cache yet, we
			// will try to download it.
			return ErrCacheMiss
		}

		c, err = xkcd.New(bytes.NewReader(data))
		if err != nil {
			c = &xkcd.Comic{
				Num:       n,
				SafeTitle: "Error parsing comic metadata from cache",
			}
			return err
		}

		return nil
	})
	if err == ErrCacheMiss {
		c, err = downloadComicInfo(n)
		if err == xkcd.ErrNotFound {
			return &xkcd.Comic{
				Num:       n,
				SafeTitle: "Comic Not Found",
			}, err
		} else if err != nil {
			return &xkcd.Comic{
				Num:       n,
				SafeTitle: "Couldn't Get Comic",
			}, err
		}
	}

	return c, err
}

// GetNewestComicInfo always returns a valid *xkcd.Comic that appears to
// be newest, and err will be set if any errors were encountered,
// however these errors can be ignored safely.
func GetNewestComicInfo() (*xkcd.Comic, error) {
	var err error

	newest := <-getCachedNewestComic

	if newest == nil {
		newest, err = xkcd.GetCurrent()
		if err != nil {
			newest, err = getNewestComicInfoFromCache()
			if err != nil {
				return newest, err
			}
			return newest, ErrOffline
		}

		setCachedNewestComic <- newest
	}
	return newest, nil
}

// GetNewestComicInfoAsync always returns a valid *xkcd.Comic that
// appears to be newest, and err will be set if any errors were
// encountered, however these errors can be ignored safely. This
// function will return the newest comic info based on the cache, but
// then asynchronously checks for the newest comic from the internet and
// calls callback when the asynchronous call completes.
func GetNewestComicInfoAsync(callback func(*xkcd.Comic, error)) (*xkcd.Comic, error) {
	newest := <-getCachedNewestComic

	if newest == nil {
		newestFromCache, err := getNewestComicInfoFromCache()

		go func() {
			newestFromInternet, err := xkcd.GetCurrent()

			if newestFromInternet != nil && err == nil {
				setCachedNewestComic <- newestFromInternet
			} else {
				setCachedNewestComic <- newestFromCache
			}

			callback(<-getCachedNewestComic, err)
		}()

		return newestFromCache, err
	}
	return newest, nil
}

func getNewestComicInfoFromCache() (*xkcd.Comic, error) {
	newest := &xkcd.Comic{
		SafeTitle: "Connect to the internet to download some comics!",
	}

	err := cacheDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			return ErrCache
		}

		return bucket.ForEach(func(k, v []byte) error {
			c, err := xkcd.New(bytes.NewReader(v))
			if err != nil {
				return err
			}

			if c.Num > newest.Num {
				newest = c
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
			return ErrCache
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

	// Now add the new file to the searchIndex.
	err = searchIndex.Index(strconv.Itoa(comic.Num), comic)

	return comic, err
}

// DownloadComicImage tries to add a comic image to our local cache. Any
// errors are indicated by err.
func DownloadComicImage(n int) error {
	c, err := GetComicInfo(n)
	if err != nil {
		return err
	}

	resp, err := http.Get(c.Img)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(getComicImagePath(n))
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

func getComicCachePath() string {
	return filepath.Join(CacheDir(), comicCacheName)
}

func getComicImagePath(n int) string {
	return filepath.Join(CacheDir(), comicCacheImageName, strconv.Itoa(n))
}

func intToBytes(i int) []byte {
	buf := make([]byte, binary.MaxVarintLen64)

	n := binary.PutVarint(buf, int64(i))

	return buf[:n]
}
