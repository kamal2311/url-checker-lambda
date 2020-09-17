package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/rs/zerolog/log"
)

type MaliciousUrlChecker struct { }

type CheckerResponse struct {
	IsSafe bool `json:"is_safe"`
	Score  int `json:"score,omitempty"`
	Source string `json:"source,omitempty"`
}

func NewMaliciousUrlChecker() *MaliciousUrlChecker {
	return &MaliciousUrlChecker{}
}

func (m *MaliciousUrlChecker) EvaluateSafety(url string) CheckerResponse {

	id := generateId(url)
	item, err := retrieveItem(id)
	if err != nil {
		return CheckerResponse{
			IsSafe: true,
		}
	}

	return CheckerResponse{
		IsSafe: false,
		Score:  item.Score,
		Source: item.Source,
	}

}

func generateId(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	out := fmt.Sprintf("%x", bs)
	log.Info().Interface("Id", out).Msg("generated")
	return out
}
