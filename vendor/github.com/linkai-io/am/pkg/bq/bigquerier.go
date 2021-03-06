package bq

import (
	"context"
	"time"

	"github.com/linkai-io/am/am"
)

type BigQuerier interface {
	Init(config, credentials []byte) error
	QueryETLD(ctx context.Context, from time.Time, etld string) (map[string]*am.CTRecord, error)
	QuerySubdomains(ctx context.Context, from time.Time, etld string) (map[string]*am.CTSubdomain, error)
}
