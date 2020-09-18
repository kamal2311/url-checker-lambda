package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
)

type DataService interface {
	GetItem(id string) (*Item, error)
	PutItem(item Item) error
}

type SaveItem struct {
	Url    string `json:"url"`
	Score  int    `json:"score"`
	Source string `json:"source"`
}

type MaliciousUrlService struct {
	dataService DataService
}

type CheckerResponse struct {
	IsSafe bool   `json:"is_safe"`
	Source string `json:"source,omitempty"`
	Score  int    `json:"score,omitempty"`
}

func NewMaliciousUrlService(ds DataService) *MaliciousUrlService {
	return &MaliciousUrlService{
		dataService: ds,
	}
}

func (m *MaliciousUrlService) Check(url string) (CheckerResponse, error) {

	id := generateId(url)
	item, err := m.dataService.GetItem(id)
	if err != nil {

		if err.Error() != ITEM_NOT_FOUND {
			return CheckerResponse{}, err
		}

		return CheckerResponse{
			IsSafe: true,
		}, nil
	}

	return CheckerResponse{
		IsSafe: false,
		Score:  item.Score,
		Source: item.Source,
	}, nil

}

func (m *MaliciousUrlService) Save(saveItem SaveItem) error {

	id := generateId(saveItem.Url)

	item := Item{
		Id:     id,
		Score:  saveItem.Score,
		Source: saveItem.Source,
	}

	return m.dataService.PutItem(item)
}

func generateId(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	out := fmt.Sprintf("%x", bs)
	log.Info().Interface("Id", out).Msg("generated")
	return out
}

func validateAndDecode(body string) (SaveItem, error) {
	var item SaveItem
	err := json.Unmarshal([]byte(body), &item)
	if err != nil {
		log.Err(err).Send()
		return SaveItem{}, errors.New("Item validation failed")
	}

	parsedUrl, err := url.Parse(item.Url)
	if err != nil {
		log.Err(err).Send()
		return SaveItem{}, errors.New("url validation failed")
	}

	item.Url = parsedUrl.Path
	if item.Url == "" {
		return SaveItem{}, errors.New("url field must not be empty")
	}

	if item.Score <= 0 || item.Score > 10 {
		return SaveItem{}, errors.New("score must be an integer between 1 and 10")
	}

	sourceLen := len(item.Source)
	if sourceLen == 0 || sourceLen > 20 {
		return SaveItem{}, errors.New("source field is mandatory and must not be longer than 20 characters")
	}

	return item, nil
}
