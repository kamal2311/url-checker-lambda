package main

import (
	"errors"
	"reflect"
	"testing"
)

type MockDataService struct{}

func (MockDataService) GetItem(id string) (*Item, error) {
	mockData := make(map[string]*Item)
	mockData["66b9cb08638d49a6d3559718551d59243fa2b0eb"] = &Item{
		Id:     "66b9cb08638d49a6d3559718551d59243fa2b0eb",
		Source: "Sophos",
		Score:  8,
	}
	mockData["4782cc39a5294f566242f9d36bccc9889916e2b6"] = &Item{
		Id:     "4782cc39a5294f566242f9d36bccc9889916e2b6",
		Source: "Malware Patrol",
		Score:  9,
	}

	if val, ok := mockData[id]; ok {
		return val, nil
	}

	return nil, errors.New(ITEM_NOT_FOUND)
}

func (MockDataService) PutItem(item Item) error {
	return nil
}

func TestMaliciousUrlService_Check(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want CheckerResponse
	}{
		{
			name: "bad-url-1 should be considered as not safe",
			args: args{
				url: "bad-url-1",
			},
			want: CheckerResponse{
				IsSafe: false,
				Score:  8,
				Source: "Sophos",
			},
		},
		{
			name: "bad-url-2 should be considered as not safe",
			args: args{
				url: "bad-url-2",
			},
			want: CheckerResponse{
				IsSafe: false,
				Score:  9,
				Source: "Malware Patrol",
			},
		},
		{
			name: "good-url-1 should be considered as safe",
			args: args{
				url: "good-url-1",
			},
			want: CheckerResponse{
				IsSafe: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaliciousUrlService{MockDataService{}}
			if got, _ := m.Check(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateId(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "bad-url-1 should have an id of 66b9cb08638d49a6d3559718551d59243fa2b0eb",
			args: args{
				s: "bad-url-1",
			},
			want: "66b9cb08638d49a6d3559718551d59243fa2b0eb",
		},
		{
			name: "bad-url-2 should have an id of 4782cc39a5294f566242f9d36bccc9889916e2b6",
			args: args{
				s: "bad-url-2",
			},
			want: "4782cc39a5294f566242f9d36bccc9889916e2b6",
		},
		{
			name: "bad-url-3 should have an id of bf889056b6f3523cedcfa5cc999b9a2df30b5e3b",
			args: args{
				s: "bad-url-3",
			},
			want: "bf889056b6f3523cedcfa5cc999b9a2df30b5e3b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateId(tt.args.s); got != tt.want {
				t.Errorf("generateId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateAndDecode(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name    string
		args    args
		want    SaveItem
		wantErr bool
	}{
		{
			name: "url must not contain back-slashes",
			args: args{body: "{ \"url\": \"some-bad-irl\\\\dfdfdfdfdf\"}"},
			want: SaveItem{},
			wantErr:true,
		},
		{
			name: "valid url should be decoded properly",
			args: args{body: "{ \"url\": \"good-url-1/123\" ,\"score\":1, \"source\":\"Sophos\"}"},
			want: SaveItem{
				Url:    "good-url-1/123",
				Score:  1,
				Source: "Sophos",
			},
			wantErr:false,
		},
		{
			name: "score must not be negative",
			args: args{body: "{ \"url\": \"good-url-1/123\" ,\"score\":-1, \"source\":\"Sophos\"}"},
			want: SaveItem{},
			wantErr:true,
		},
		{
			name: "score must not be a non-integer",
			args: args{body: "{ \"url\": \"good-url-1/123\" ,\"score\":1.5, \"source\":\"Sophos\"}"},
			want: SaveItem{},
			wantErr:true,
		},
		{
			name: "score must not be larger than 10",
			args: args{body: "{ \"url\": \"good-url-1/123\" ,\"score\":11, \"source\":\"Sophos\"}"},
			want: SaveItem{},
			wantErr:true,
		},
		{
			name: "source must not be empty",
			args: args{body: "{ \"url\": \"good-url-1/123\" ,\"score\":10, \"source\":\"\"}"},
			want: SaveItem{},
			wantErr:true,
		},
		{
			name: "source must not longer than 20 characters",
			args: args{body: "{ \"url\": \"good-url-1/123\" ,\"score\":10, \"source\":\"0102030405060708092021\"}"},
			want: SaveItem{},
			wantErr:true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateAndDecode(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateAndDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}