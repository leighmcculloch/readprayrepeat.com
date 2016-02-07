package main

import (
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

const cacheFolder = "cache"

func GetBiblePassageWithCache(bible Bible, reference string) (*BiblePassage, error) {
	bp, err := loadBiblePassage(bible, reference)
	if err != nil {
		bp, err = bible.GetPassage(reference)
		if err != nil {
			return nil, err
		}
		err = saveBiblePassage(bible, reference, bp)
		if err != nil {
			return nil, err
		}
		return bp, nil
	}
	return bp, nil
}

func loadBiblePassage(bible Bible, reference string) (*BiblePassage, error) {
	f, err := os.Open(getCacheFilePath(bible, reference))
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

func saveBiblePassage(bible Bible, reference string, biblePassage *BiblePassage) error {
	f, err := os.Create(getCacheFilePath(bible, reference))
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

func getCacheFilePath(bible Bible, reference string) string {
	hasher := sha1.New()
	io.WriteString(hasher, bible.Source())
	io.WriteString(hasher, bible.NameShort())
	io.WriteString(hasher, reference)
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%s/%0x", cacheFolder, hash)
}
