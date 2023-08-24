package parser

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackageLockParser_Parse(t *testing.T) {
	parser := NewPackageLockParser()

	err := parser.Init(context.Background())
	assert.Nil(t, err)

	input := &PackageLockJsonParserInput{
		PackageLockJsonPath: "./test_data/package-lock.json/join-dev-design.json",
	}
	project, err := parser.Parse(context.Background(), input)
	assert.Nil(t, err)
	assert.NotNil(t, project)
	t.Log(project)

}
