package main

import (
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"os"
)

const cacheFolder = "cache"

func GetBiblePassageWithCache(reference string) (*BiblePassage, error) {
	bp, err := loadBiblePassage(reference)
	if err != nil {
		bp, err = GetBiblePassage(reference)
		if err != nil {
			return nil, err
		}
		err = saveBiblePassage(reference, bp)
		if err != nil {
			return nil, err
		}
		return bp, nil
	}
	return bp, nil
}

func loadBiblePassage(reference string) (*BiblePassage, error) {
	f, err := os.Open(getCacheFilePath(reference))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bp BiblePassage
	dec := gob.NewDecoder(f)
	err = dec.Decode(&bp)
	if err != nil {
		return nil, err
	}

	return &bp, nil
}

func saveBiblePassage(reference string, biblePassage *BiblePassage) error {
	f, err := os.Create(getCacheFilePath(reference))
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	enc.Encode(biblePassage)
	if err != nil {
		return err
	}

	return nil
}

func getCacheFilePath(reference string) string {
	hash := sha1.Sum([]byte(reference))
	return fmt.Sprintf("%s/%0x", cacheFolder, hash)
}
