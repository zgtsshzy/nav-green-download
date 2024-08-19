package download

import (
	"testing"
	"time"
)

func TestGetSMOCDownloadUrlByDate(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2024-06-30")

	url, err := GetSMOCDownloadUrlByDate(date)
	if err != nil {
		t.Errorf("GetSMOCDownloadUrlByDate: %v", err)
		return
	}

	t.Logf("download url: %s", url)
}

func TestGetSMOCNameByDate(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2024-06-30")

	name, err := GetSMOCNameByDate(date)
	if err != nil {
		t.Errorf("GetSMOCNameByDate: %v", err)
		return
	}

	t.Logf("remote name: %s", name)
}
