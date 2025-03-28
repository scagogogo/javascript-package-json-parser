package parser

import (
	"context"
	"os"
)

// YarnLockParserInput 解析器的输入，指定要解析的yarn.lock文件路径
type YarnLockParserInput struct {
	// YarnLockPath yarn.lock文件的路径
	YarnLockPath string
}

// Read 读取yarn.lock文件内容
func (x *YarnLockParserInput) Read(ctx context.Context) ([]byte, error) {
	// 读取文件内容
	return os.ReadFile(x.YarnLockPath)
}
