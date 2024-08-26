package parser

type ContourLineConf struct {
	MinValue float32 // 最大值
	MaxValue float32 // 最小值
	Step     float32 // 步长
	Num      int     // 分段个数
}

// 气象文件解析器
type Parser interface {
	ParseData() error                             // 文件解析
	DrawGrayscale() error                         // 画灰度图
	DrawContourLine(config ContourLineConf) error // 画等值线
	Save2DB(tableName string) error               // 数据保存到数据库
}
