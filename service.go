package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/rs/zerolog/log"
)

type DataService interface {
	RetrieveItem(id string) (*Item, error)
	InsertItem(item Item) error
}

type MaliciousUrlChecker struct {
	dataService DataService
}

type CheckerResponse struct {
	IsSafe bool `json:"is_safe"`
	Score  int `json:"score,omitempty"`
	Source string `json:"source,omitempty"`
}

func NewMaliciousUrlChecker(ds DataService) *MaliciousUrlChecker {
	return &MaliciousUrlChecker{
		dataService: ds,
	}
}

func (m *MaliciousUrlChecker) EvaluateSafety(url string) (CheckerResponse, error) {

	id := generateId(url)
	item, err := m.dataService.RetrieveItem(id)
	if err != nil  {

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

func generateId(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	out := fmt.Sprintf("%x", bs)
	log.Info().Interface("Id", out).Msg("generated")
	return out
}
