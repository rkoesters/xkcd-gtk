package cache

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

var imageFileNameRegex = regexp.MustCompile("^[0-9][0-9]*$")

type Stat struct {
	LatestComicNumber int
	CachedCount       int
}

func (s Stat) Complete() bool {
	return s.LatestComicNumber == s.CachedCount
}

func (s Stat) Fraction() (float64, error) {
	if s.LatestComicNumber == 0 {
		return 0, errors.New("division by zero")
	}
	return float64(s.CachedCount) / float64(s.LatestComicNumber), nil
}

func (s Stat) String() string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(s.CachedCount))
	b.WriteString(" / ")
	b.WriteString(strconv.Itoa(s.LatestComicNumber))
	return b.String()
}

func StatMetadata() (Stat, error) {
	var s Stat
	latestComic, err := CheckForNewestComicInfo(time.Second)
	if err != nil {
		return s, err
	}
	s.LatestComicNumber = latestComic.Num
	s.CachedCount, err = countCachedMetadata()
	return s, err
}

func countCachedMetadata() (int, error) {
	var count int

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

func StatImages() (Stat, error) {
	var s Stat
	latestComic, err := CheckForNewestComicInfo(time.Second)
	if err != nil {
		return s, err
	}
	s.LatestComicNumber = int(latestComic.Num)
	s.CachedCount, err = countCachedImages()
	return s, err
}

func countCachedImages() (int, error) {
	var count int

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
