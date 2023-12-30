package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractChapter(t *testing.T) {
	assert.EqualValues(t, 41, extractChapterInfoFromFilename("+_-_41_"))
	assert.EqualValues(t, 2, extractChapterInfoFromFilename(" - 02."))
	assert.EqualValues(t, 1, extractChapterInfoFromFilename(" - S01E01."))
	assert.EqualValues(t, 3, extractChapterInfoFromFilename("_-_03."))
	assert.EqualValues(t, 1, extractChapterInfoFromFilename(" S2 - 01."))
	assert.EqualValues(t, 4, extractChapterInfoFromFilename(" - S02E04."))
	assert.EqualValues(t, 4, extractChapterInfoFromFilename("86 - 04."))
	assert.EqualValues(t, 13, extractChapterInfoFromFilename(" Lv999 (13)"))
	assert.EqualValues(t, 7, extractChapterInfoFromFilename(" 第07話 "))
}
