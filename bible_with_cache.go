package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
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
	err := os.MkdirAll(getCacheFileDir(bible, reference), os.ModePerm)
	if err != nil {
		return err
	}
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

func getCacheFileDir(bible Bible, reference string) string {
	return filepath.Join(cacheFolder, bible.Source(), bible.NameShort())
}

func getCacheFileName(bible Bible, reference string) string {
	return fmt.Sprintf("%s.biblepassage", reference)
}

func getCacheFilePath(bible Bible, reference string) string {
	return filepath.Join(getCacheFileDir(bible, reference), getCacheFileName(bible, reference))
}
