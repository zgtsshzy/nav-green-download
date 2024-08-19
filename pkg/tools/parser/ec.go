package parser

var _ Parser = (*ECParser)(nil)

type ECParser struct {
}

func NewECParser() (*ECParser, error) {

	return nil, nil
}

func (parser *ECParser) ParseData() error {
	return nil
}

func (parser *ECParser) DrawGrayscale() error {
	return nil
}

func (parser *ECParser) DrawContourLine(config ContourLineConf) error {
	return nil
}

func (parser *ECParser) Save2DB(tableName string) error {
	return nil
}
