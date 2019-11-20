package util

import (
	"reflect"
	"testing"
)


func Test_parsePid(t *testing.T) {
	type args struct {
		pids []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"123",args{[]byte{'1','2','3'}},123},
		{"123",args{[]byte{'2','2','3'}},223},
		{"123",args{[]byte{'6','2','6'}},626},
		{"123",args{[]byte{'1','4','3'}},143},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePid(tt.args.pids); got != tt.want {
				t.Errorf("parsePid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPidFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"first test",args{"D:/test.txt"},[]byte{'1','2','3'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPidFromFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPidFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"files",args{"D:/test.txt"},true},
		{"not exist",args{"D:/test.pid"},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.path); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}