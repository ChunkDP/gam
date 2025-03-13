package services

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// 使用免费的翻译API
func TranslateText(text string) (string, error) {
	// 这里使用 MyMemory 免费翻译API作为示例
	// 实际使用时可以替换为其他翻译服务
	apiURL := "https://api.mymemory.translated.net/get"

	params := url.Values{}
	params.Add("q", text)
	params.Add("langpair", "en|zh")

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ResponseData.TranslatedText, nil
}
