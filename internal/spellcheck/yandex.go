package spellcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Spellchecker interface {
	CheckSpelling(text string) (string, error)
}

// YandexSpellchecker предоставляет методы для проверки орфографии с помощью API Яндекс.Спеллер
type YandexSpellchecker struct {
	apiURL string
}

// NewYandexSpellchecker создает новый экземпляр YandexSpellchecker
func NewYandexSpellchecker(apiURL string) *YandexSpellchecker {
	return &YandexSpellchecker{
		apiURL: apiURL,
	}
}

// SpellCheckResult представляет результат проверки орфографии для одного слова
type SpellCheckResult struct {
	Word string   `json:"word"`
	S    []string `json:"s"`
	Code int      `json:"code"`
}

// CheckSpelling проверяет орфографию в тексте
func (y *YandexSpellchecker) CheckSpelling(text string) (string, error) {
	params := url.Values{}
	params.Add("text", text)

	resp, err := http.Get(y.apiURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("failed to send request to Yandex.Speller: %w", err)
	}
	defer resp.Body.Close()

	var results []SpellCheckResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", fmt.Errorf("failed to decode Yandex.Speller response: %w", err)
	}

	correctedText := text
	for _, result := range results {
		if len(result.S) > 0 {
			correctedText = strings.Replace(correctedText, result.Word, result.S[0], 1)
		}
	}

	return correctedText, nil
}
