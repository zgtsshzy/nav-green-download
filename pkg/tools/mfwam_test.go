package tools

import (
	"testing"
	"time"
)

func TestGetMFWAMDownloadUrlByDate(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2024-06-30")

	url, err := GetMFWAMDownloadUrlByDate(date)
	if err != nil {
		t.Errorf("GetMFWAMDownloadUrlByDate: %v", err)
		return
	}

	t.Logf("download url: %s", url)
}

func TestGetMFWAMNameByDate(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2024-06-30")

	name, err := GetMFWAMNameByDate(date)
	if err != nil {
		t.Errorf("GetMFWAMNameByDate: %v", err)
		return
	}

	t.Logf("remote name: %s", name)
}
