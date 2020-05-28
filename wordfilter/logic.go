package wordfilter

import (
	"github.com/gin-gonic/gin"
)

type InterfaceHttp struct {
}

type FindRequest struct {
	Content string `json:"content"`
	Rank    uint8  `json:"rank"`
}

type FindResponse struct {
	Status     uint8                   `json:"status"`
	NewContent string                  `json:"new_content"`
	ErrMsg     string                  `json:"err_msg"`
	BadWords   map[uint8][]*SearchItem `json:"bad_words"`
}

func (http *InterfaceHttp) reload(context *gin.Context) {
	LoadWords()
	result := make(map[string]string)
	result["status"] = "1"
	result["msg"] = "has push cmd :reload"
	context.JSON(200, result)

}

func (http *InterfaceHttp) search(context *gin.Context) {
	returnResult := &FindResponse{
		Status:     1,
		NewContent: "",
		BadWords:   make(map[uint8][]*SearchItem),
	}

	var req FindRequest
	err := context.ShouldBind(&req)
	if err != nil {
		returnResult.ErrMsg = err.Error()
		context.JSON(200, returnResult)
		return
	}

	ac := getAc()
	if ac == nil {
		returnResult.ErrMsg = "ac is nil"
		context.JSON(200, returnResult)
		return
	}

	result := ac.Search(req.Content)
	contentBuff := []rune(req.Content)
	for _, item := range result {
		if item.Rank > req.Rank && req.Rank != 0 {
			continue
		}
		for i := item.StartP; i <= item.EndP; i++ {
			contentBuff[i] = '*'
		}
		returnResult.BadWords[item.Rank] = append(returnResult.BadWords[item.Rank], item)
	}
	returnResult.Status = 0
	returnResult.NewContent = string(contentBuff)
	context.JSON(200, returnResult)
}
