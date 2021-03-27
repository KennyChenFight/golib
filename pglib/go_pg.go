// Package pglib is for encapsulating github.com/go-pg/pg any operations
//
// As a quick start:
// 	cfg := pglib.GOPGConfig{
//		URL:       "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
//		DebugMode: true,
//		PoolSize:  5,
//	}
//	client, err := pglib.NewDefaultGOPGClient(cfg)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(client.Ping(context.Background()))
package pglib

import (
	"context"
	"github.com/go-pg/pg/v10"
	"log"
	"time"
)

func NewDefaultGOPGClient(config GOPGConfig) (*GOPGClient, error) {
	opts, err := pg.ParseURL(config.URL)
	if err != nil {
		return nil, err
	}
	opts.PoolSize = config.PoolSize

	client := &GOPGClient{pg.Connect(opts)}
	if config.DebugMode {
		client.DB.AddQueryHook(&GOPGDebugQueryHook{})
	}
	return client, nil
}

func NewCustomizeGOPGClient(opts *pg.Options, debugMode bool) (*GOPGClient, error) {
	client := &GOPGClient{pg.Connect(opts)}
	if debugMode {
		client.DB.AddQueryHook(&GOPGDebugQueryHook{})
	}
	return client, nil
}

type GOPGConfig struct {
	URL       string
	DebugMode bool
	PoolSize  int
}

type GOPGClient struct {
	*pg.DB
}

type GOPGDebugQueryHook struct {
}

func (h *GOPGDebugQueryHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	log.Printf("BeforeQuery:\nparams:%v\nquery:%v\n", event.Params, event.Query)
	return ctx, nil
}

func (h *GOPGDebugQueryHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	//query, _ := event.FormattedQuery()
	log.Printf("AfterQuery:\nparams:%v\nquery:%v\nstartTime:%v\nduration:%vs\nerr:%v\n", event.Params, event.Query, event.StartTime, time.Now().Sub(event.StartTime).Seconds(), event.Err)
	return nil
}
