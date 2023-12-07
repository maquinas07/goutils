package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractChapter(t *testing.T) {
	assert.EqualValues(t, 41, extractChapterInfoFromFilename("_Blood+_-_41_.mkv"))
	assert.EqualValues(t, 2, extractChapterInfoFromFilename("Babylon - 02.mkv"))
	assert.EqualValues(t, 1, extractChapterInfoFromFilename("Occult Academy - S01E01.mkv"))
	assert.EqualValues(t, 3, extractChapterInfoFromFilename("Seikon_no_Qwaser_II_-_03.mkv"))
	assert.EqualValues(t, 1, extractChapterInfoFromFilename("Kage no Jitsuryokusha ni Naritakute! S2 - 01.srt"))
	assert.EqualValues(t, 4, extractChapterInfoFromFilename("The Eminence in Shadow - S02E04.mkv"))
	assert.EqualValues(t, 4, extractChapterInfoFromFilename("86 - 04.mkv"))
	assert.EqualValues(t, 13, extractChapterInfoFromFilename("Yamada-kun to Lv999 no Koi wo Suru (13)"))
	assert.EqualValues(t, 7, extractChapterInfoFromFilename("[AT-X fix] 第07話 古老、曰く.ja.srt"))
}
