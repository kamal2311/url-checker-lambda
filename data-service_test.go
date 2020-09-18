package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

func Test_retrieveItem(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		shouldError bool
		want        *Item
	}{
		{
			"should retrieve an item",
			args{
				id: "66b9cb08638d49a6d3559718551d59243fa2b0eb",
			},
			true,
			&Item{
				Id:     "66b9cb08638d49a6d3559718551d59243fa2b0eb",
				Source: "Sophos",
				Score:  8,
			},
		},
		{
			"should not retrieve a missing item",
			args{
				id: "good-url",
			},
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tableName := os.Getenv("MC_TABLE_NAME")
			dataService := NewDynamoDataService(tableName)

			ans, err := dataService.GetItem(tt.args.id)
			if err != nil && !tt.shouldError {
				t.Errorf("Expected an error but got %v", err)
				return
			}

			if tt.want != nil && ans != nil && *tt.want != *ans {
				t.Errorf("expected %v, got %v", tt.want, ans)
			}

			log.Info().Interface("result", ans).Send()
		})
	}
}
