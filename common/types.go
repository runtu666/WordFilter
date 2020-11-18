package common

type (
	FindResponse struct {
		Status     uint8                   `json:"status"`
		NewContent string                  `json:"new_content"`
		ErrMsg     string                  `json:"err_msg"`
		BadWords   map[uint8][]*SearchItem `json:"bad_words"`
	}
	SearchItem struct {
		StartP int    `json:"start_p"`
		EndP   int    `json:"end_p"`
		Word   string `json:"word"`
		Rank   uint8  `json:"rank"`
	}
	SensitiveWords struct {
		Word string `json:"word"`
		Rank uint8  `json:"rank"`
	}
)
