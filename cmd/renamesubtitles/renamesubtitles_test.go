package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractChapter(t *testing.T) {
	assert.EqualValues(t, extractChapterInfoFromFilename("_Blood+_-_41_.mkv"), 41)
	assert.EqualValues(t, extractChapterInfoFromFilename("Babylon - 02.mkv"), 2)
	assert.EqualValues(t, extractChapterInfoFromFilename("Occult Academy - S01E01.mkv"), 1)
	assert.EqualValues(t, extractChapterInfoFromFilename("Seikon_no_Qwaser_II_-_03.mkv"), 3)
	assert.EqualValues(t, extractChapterInfoFromFilename("Kage no Jitsuryokusha ni Naritakute! S2 - 01.srt"), 1)
	assert.EqualValues(t, extractChapterInfoFromFilename("The Eminence in Shadow - S02E04.mkv"), 4)
	assert.EqualValues(t, extractChapterInfoFromFilename("86 - 04.mkv"), 4)
}
