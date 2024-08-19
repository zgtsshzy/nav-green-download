package download

import (
	"testing"
)

func TestGetECFirstLevel(t *testing.T) {
	folders, err := GetECFirstLevel()
	if err != nil {
		t.Errorf("GetECFirstLevel: %v", err)
		return
	}

	t.Logf("folders: %v", folders)
}

func TestGetECSecondLevel(t *testing.T) {
	firstFolders, err := GetECFirstLevel()
	if err != nil {
		t.Errorf("GetECFirstLevel: %v", err)
		return
	}

	for _, level := range firstFolders {
		folders, err := GetECSecondLevel(level)
		if err != nil {
			t.Errorf("GetECSecondLevel: %v", err)
			return
		}

		t.Logf("folders: %v", folders)
	}
}

func TestGetECFifthLevel(t *testing.T) {
	firstFolders, err := GetECFirstLevel()
	if err != nil {
		t.Errorf("GetECFirstLevel: %v", err)
		return
	}

	for _, first := range firstFolders {
		secondFolders, err := GetECSecondLevel(first)
		if err != nil {
			t.Errorf("GetECSecondLevel: %v", err)
			return
		}

		for _, second := range secondFolders {
			fifth := first + second + "ifs/" + "0p25/"
			fifthFolders, err := GetECFifthLevel(fifth)
			if err != nil {
				t.Errorf("GetECSecondLevel: %v", err)
				return
			}

			t.Logf("fifthFolders: %v", fifthFolders)
		}
	}
}

func TestGetECSixthFiles(t *testing.T) {
	firstFolders, err := GetECFirstLevel()
	if err != nil {
		t.Errorf("GetECFirstLevel: %v", err)
		return
	}

	for _, first := range firstFolders {
		secondFolders, err := GetECSecondLevel(first)
		if err != nil {
			t.Errorf("GetECSecondLevel: %v", err)
			return
		}

		for _, second := range secondFolders {
			fifthLevel := first + second + "ifs/" + "0p25/"
			fifthFolders, err := GetECFifthLevel(fifthLevel)
			if err != nil {
				t.Errorf("GetECSecondLevel: %v", err)
				return
			}

			for _, fifth := range fifthFolders {
				sixthLevel := first + second + "ifs/" + "0p25/" + fifth
				sixthFolders, err := GetECSixthFiles(sixthLevel)
				if err != nil {
					t.Errorf("GetECSixthFiles: %v", err)
					return
				}

				t.Logf("sixthFolders: %v", sixthFolders)
			}
		}
	}
}
