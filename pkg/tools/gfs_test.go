package tools

import "testing"

func TestGetGFSFirstLevel(t *testing.T) {
	folders, err := GetGFSFirstLevel()
	if err != nil {
		t.Errorf("GetGFSFirstLevel: %v", err)
		return
	}

	t.Logf("%v", folders)
}
func TestGetGFSSecondLevel(t *testing.T) {
	folders, err := GetGFSSecondLevel("gfs.20240816/")
	if err != nil {
		t.Errorf("GetGFSSecondLevel: %v", err)
		return
	}

	t.Logf("%v", folders)
}

func TestGetGFSThirdLevel(t *testing.T) {
	folders, err := GetGFSThirdLevel("gfs.20240816/00/")
	if err != nil {
		t.Errorf("GetGFSThirdLevel: %v", err)
		return
	}

	t.Logf("%v", folders)
}

func TestGetGFSFourthLevel(t *testing.T) {
	folders, err := GetGFSFourthLevel("gfs.20240816/00/atmos/")
	if err != nil {
		t.Errorf("GetGFSFourthLevel: %v", err)
		return
	}

	for _, folder := range folders {
		t.Logf("%v", folder)
	}

	folders, err = GetGFSFourthLevel("gfs.20240816/00/wave/")
	if err != nil {
		t.Errorf("GetGFSFourthLevel: %v", err)
		return
	}

	for _, folder := range folders {
		t.Logf("%v", folder)
	}
}
