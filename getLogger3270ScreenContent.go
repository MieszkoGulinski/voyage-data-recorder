package main

import "github.com/racingmars/go3270"

var screenTitles = []string{
	"GPS Position",
	"Weather",
	"Battery status",
}

func getLogger3270ScreenContent (
	currentScreenId int,
  historyStack []int64,
) (
	screenContent go3270.Screen,
	lastTimestamp int64,
	nextPageExists bool,
	err error,
) {
	// TODO read from actual database and append the appropriate lines

	// Default size is 80x24 characters - TODO handle more screen sizes if available
	screenContent = go3270.Screen{
		{Row: 0, Col: 0, Intense: true, Content: screenTitles[currentScreenId]},
		{Row: 23, Col: 0, Content: "F3 close F7 page up F8 page down F9 next table F10 prev table"},
	}

	return
}