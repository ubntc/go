package arts

import "strings"

const (
	Block_LeftHalf      = "▌"
	Block_RightHalf     = "▐"
	Block_RightOneEight = "▕"

	Block_LeftEights_1 = "▏"
	Block_LeftEights_2 = "▎"
	Block_LeftEights_3 = "▍"
	Block_LeftEights_4 = "▌"
	Block_LeftEights_5 = "▋"
	Block_LeftEights_6 = "▊"
	Block_LeftEights_7 = "▉"
	Block_LeftEights_8 = "█"
)

var Blocks_LeftEights = strings.Split("▏▎▍▌▋▊▉█", "")
