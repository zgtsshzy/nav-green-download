package tools

import (
	"fmt"
	"nav-green-download/pkg/tools/parser"
)

const (
	ECType     FileType = "EC"
	GFSType    FileType = "GFS"
	MFWAMType  FileType = "MFWAM"
	SEAICEType FileType = "SEAICE"
	SMOCType   FileType = "SMOC"
)

type FileType string

func NewParser(path string, fileType FileType) (parser.Parser, error) {
	switch fileType {
	case ECType:
		return nil, nil
	case GFSType:
		return nil, nil
	case MFWAMType:
		return nil, nil
	case SEAICEType:
		return nil, nil
	case SMOCType:
		return nil, nil
	default:
		return nil, fmt.Errorf("fileType 不合法")
	}
}
