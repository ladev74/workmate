package filesystem

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	storage := &Storage{
		logger: zap.NewNop(),
	}

	tests := []struct {
		name     string
		path     string
		fileName string
		wantErr  error
	}{
		{
			name:     "success",
			path:     "test",
			fileName: "test.txt",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Init(tt.path, tt.fileName)
			fmt.Println(err)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}
