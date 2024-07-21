package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/maquinas07/golibs/ascii"
	"golang.org/x/exp/slices"
)

const (
	reverseFilename = ".reverse"
)

var videoExtension = "mkv"
var subExtension = "srt"
var language = "ja"
var reverse, verbose bool
var sortByName bool

func extractChapterInfoFromFilename(filename string) int {
	chapter := -1
	isCommonEpisodeMarker := func(rune rune) bool {
		return rune == '-' || rune == ' ' || rune == '_' || rune == '.'
	}
	isOpenEpisodeMarker := func(rune rune) bool {
		return rune == '(' || rune == '第' || rune == 'E' || rune == 'e' || rune == '['
	}
	isCloseEpisodeMarker := func(rune rune) bool {
		return rune == ')' || rune == '話' || rune == ']'
	}
	isSubtitleVersionMarker := func(rune rune) bool {
		return rune == 'v'
	}
	var common, wasOpenEpisodeMarker bool
	for i, ir := range filename {
		if (common || wasOpenEpisodeMarker) && ascii.IsDigit(byte(ir)) {
			for j, jr := range filename[i+1:] {
				if ascii.IsDigit(byte(jr)) {
					continue
				}
				if isCommonEpisodeMarker(rune(jr)) || isSubtitleVersionMarker(rune(jr)) || (wasOpenEpisodeMarker && isCloseEpisodeMarker(rune(jr))) || chapter == -1 {
					maybeChapter, err := ascii.ParseInt([]byte(filename[i : i+j+1]))
					if err == nil {
						chapter = maybeChapter
						if wasOpenEpisodeMarker {
							return chapter
						}
					}
					break
				}
			}
			common = false
			wasOpenEpisodeMarker = false
		}
		common = isCommonEpisodeMarker(ir)
		wasOpenEpisodeMarker = isOpenEpisodeMarker(ir)
	}
	return chapter
}

func globFilenamesSortedByChapter(dir, pattern string) (files map[string]string, err error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return // ignore I/O error
	}
	if !fi.IsDir() {
		return // ignore I/O error
	}
	d, err := os.Open(dir)
	if err != nil {
		return // ignore I/O error
	}
	defer d.Close()

	names, _ := d.Readdirnames(-1)
	var matchedFiles []string
	for _, n := range names {
		matched, err := filepath.Match(pattern, n)
		if err != nil {
			return files, err
		}
		if matched {
			matchedFiles = append(matchedFiles, filepath.Join(dir, n))
		}
	}

	files = make(map[string]string)
	for i := 0; i < len(matchedFiles); i++ {
		chapter := extractChapterInfoFromFilename(matchedFiles[i])
		if chapter > 0 {
			chapterStr := strconv.Itoa(chapter)
			files[chapterStr] = matchedFiles[i]
		}
	}

	return
}

func globFilenamesSortedByName(dir, pattern string) (files map[string]string, err error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return // ignore I/O error
	}

	if !fi.IsDir() {
		return // ignore I/O error
	}

	d, err := os.Open(dir)
	if err != nil {
		return // ignore I/O error
	}

	defer d.Close()

	names, _ := d.Readdirnames(-1)
	var matchedFiles []string
	for _, n := range names {
		matched, err := filepath.Match(pattern, n)
		if err != nil {
			return files, err
		}
		if matched {
			matchedFiles = append(matchedFiles, filepath.Join(dir, n))
		}
	}
	slices.Sort(matchedFiles)

	files = make(map[string]string)
	for chapter := 0; chapter < len(matchedFiles); chapter++ {
		chapterStr := strconv.Itoa(chapter)
		files[chapterStr] = matchedFiles[chapter]
	}

	return
}

func trimExtension(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func reportAndDie(err error, code int) {
	fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
	os.Exit(code)
}

func reverseRenameSubtitles(dir string) error {
	reverseFile, err := os.Open(reverseFilename)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(reverseFile)
	for scanner.Scan() {
		fileNames := strings.Split(string(scanner.Bytes()), "=")
		if len(fileNames) == 2 {
			oldName := fileNames[0]
			newName := fileNames[1]
			os.Rename(newName, oldName)
			if verbose {
				fmt.Printf("Renaming %s into %s\n", newName, oldName)
			}

		} else {
			return fmt.Errorf("Unexpected line found.\n")
		}
	}
	reverseFile.Close()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	err = os.Remove(reverseFilename)
	return err
}

func renameSubtitles(dir string) error {
	var globSorted func(string, string) (files map[string]string, err error)
	if sortByName {
		globSorted = globFilenamesSortedByName
	} else {
		globSorted = globFilenamesSortedByChapter
	}

	videoFiles, err := globSorted(dir, fmt.Sprintf("*.%s", videoExtension))
	if err != nil {
		reportAndDie(err, -1)
	}
	subFiles, err := globSorted(dir, fmt.Sprintf("*.%s", subExtension))
	if err != nil {
		reportAndDie(err, -1)
	}

	_, err = os.Stat(reverseFilename)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("renamesubtitles already executed in this directory, remove `./.reverse` or revert the renaming before using again.\n")
	}
	reverseFile, err := os.Create(reverseFilename)
	defer reverseFile.Close()
	if err != nil {
		return err
	}
	for i := range videoFiles {
		if videoFiles[i] != "" && subFiles[i] != "" {
			oldSubName := filepath.Base(subFiles[i])
			newSubName := fmt.Sprintf("%s.%s.%s", trimExtension(videoFiles[i]), language, subExtension)
			os.Rename(oldSubName, newSubName)
			reverseFile.WriteString(fmt.Sprintf("%s=%s\n", oldSubName, newSubName))
			if verbose {
				fmt.Printf("Renaming %s into %s\n", oldSubName, newSubName)
			}
		}
	}
	return nil
}
