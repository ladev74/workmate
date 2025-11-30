package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"link-service/internal/domain"
	filesystem "link-service/internal/repository/file_system"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name       string
		serverCtx  context.Context
		requestCtx context.Context
		links      []string
		wantRec    *domain.Record
		wantErr    error
	}{
		{
			name:       "success",
			serverCtx:  context.Background(),
			requestCtx: context.Background(),
			links: []string{
				"google.com",
				"yandex.ru",
			},
			wantRec: &domain.Record{
				Links: map[string]string{
					"google.com": statusAvailable,
					"yandex.ru":  statusAvailable,
				},
				ID: 1,
			},
			wantErr: nil,
		},
		{
			name:       "not available",
			serverCtx:  context.Background(),
			requestCtx: context.Background(),
			links: []string{
				"12dqf4wgf4.com",
				"yandex.ru",
			},
			wantRec: &domain.Record{
				Links: map[string]string{
					"12dqf4wgf4.com": statusNotAvailable,
					"yandex.ru":      statusAvailable,
				},
				ID: 1,
			},
			wantErr: nil,
		},
		{
			name: "app stoped",
			serverCtx: func() context.Context {
				canceledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				return canceledCtx
			}(),
			requestCtx: context.Background(),
			links: []string{
				"google.com",
				"yandex.ru",
			},
			wantRec: &domain.Record{
				Links: map[string]string{
					"google.com": statusUnknown,
					"yandex.ru":  statusUnknown,
				},
				ID: 1,
			},
			wantErr: ErrAppStopped,
		},
		{
			name:      "request context canceled",
			serverCtx: context.Background(),
			requestCtx: func() context.Context {
				canceledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				return canceledCtx
			}(),
			links: []string{
				"google.com",
				"yandex.ru",
			},
			wantRec: nil,
			wantErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := New(filesystem.NewMockStorage(), 30*time.Second, zap.NewNop())

			gotRec, err := srv.Process(tt.serverCtx, tt.requestCtx, tt.links)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRec, gotRec)
		})
	}
}
