package kurrent

import (
	"context"
	"fmt"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"

	"github.com/marcelofabianov/glitchbuster-order-api/config"
)

type DBClient interface {
	AppendToStream(
		ctx context.Context,
		streamName string,
		opts kurrentdb.AppendToStreamOptions,
		events ...kurrentdb.EventData,
	) (*kurrentdb.WriteResult, error)

	Close() error
}

type client struct {
	*kurrentdb.Client
}

func NewClient(cfg config.KurrentDBConfig) (DBClient, error) {
	if cfg.ConnectionString == "" {
		return nil, fmt.Errorf("connection string cannot be empty")
	}

	settings, err := kurrentdb.ParseConnectionString(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	db, err := kurrentdb.NewClient(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kurrentdb client: %w", err)
	}

	return &client{Client: db}, nil
}
