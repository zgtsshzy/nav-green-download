package download

import (
	"testing"
	"time"
)

func TestGetSeaIceDownloadUrlByDate(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2024-06-30")

	url, err := GetSeaIceDownloadUrlByDate(date)
	if err != nil {
		t.Errorf("GetSeaIceDownloadUrlByDate: %v", err)
		return
	}

	t.Logf("seaIce url: %s", url)
}

func TestGetSeaIceNameByDate(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2024-06-30")

	name, err := GetSeaIceNameByDate(date)
	if err != nil {
		t.Errorf("GetSeaIceNameByDate: %v", err)
		return
	}

	t.Logf("seaIce name: %s", name)
}
