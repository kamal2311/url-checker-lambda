package main

import (
	"reflect"
	"testing"
)

func TestMaliciousUrlChecker_EvaluateSafety(t *testing.T) {
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
			m := &MaliciousUrlChecker{}
			if got := m.EvaluateSafety(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvaluateSafety() = %v, want %v", got, tt.want)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateId(tt.args.s); got != tt.want {
				t.Errorf("generateId() = %v, want %v", got, tt.want)
			}
		})
	}
}
