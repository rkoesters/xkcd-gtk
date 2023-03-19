package cache

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	bolt "go.etcd.io/bbolt"
)

var imageFileNameRegex = regexp.MustCompile("^[0-9][0-9]*$")

type CacheStat struct {
	LatestComicNumber uint
	CachedCount       uint
}

func (cs CacheStat) Fraction() (float64, error) {
	if cs.LatestComicNumber == 0 {
		return 0, errors.New("division by zero")
	}
	return float64(cs.CachedCount) / float64(cs.LatestComicNumber), nil
}

func (cs CacheStat) String() string {
	return fmt.Sprintln(cs.CachedCount, "/", cs.LatestComicNumber)
}

func StatMetadata() (CacheStat, error) {
	var cs CacheStat
	latestComic, err := NewestComicInfoFromCache()
	if err != nil {
		return cs, err
	}
	cs.LatestComicNumber = uint(latestComic.Num)
	cs.CachedCount, err = countCachedMetadata()
	return cs, err
}

func countCachedMetadata() (uint, error) {
	var count uint

	// We are ready to display comic #404, which is an error page rather than an
	// image.
	count++

	err := cacheDB.View(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket(comicCacheMetadataBucketName)
		if bucket == nil {
			return ErrLocalFailure
		}

		err = bucket.ForEach(func(_, _ []byte) error {
			count++
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	return count, err
}

func StatImages() (CacheStat, error) {
	var cs CacheStat

	latestComic, err := NewestComicInfoFromCache()
	if err != nil {
		return cs, err
	}
	cs.LatestComicNumber = uint(latestComic.Num)

	cs.CachedCount, err = countCachedImages()
	return cs, err
}

func countCachedImages() (uint, error) {
	var count uint

	// We are ready to display comic #404, which is an error page rather than an
	// image.
	count++

	files, err := os.ReadDir(comicImageDirPath())
	if err != nil {
		return count, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !imageFileNameRegex.MatchString(file.Name()) {
			continue
		}
		count++
	}

	return count, nil
}
