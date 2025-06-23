package server

import (
	"context"
	"time"

	commonPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/mailer"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store/dbstore"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *Server) initTrans() map[string]*commonPb.TranslationElements {
	trans, err := a.commonClient.TranslationsGet()
	if err != nil {
		a.errors <- err
	}

	if err := models.TranslationsInit(trans); err != nil {
		err := &models.InternalError{
			Err:  err,
			Msg:  "failed to init translations",
			Path: "user.server.initTrans",
		}
		a.errors <- err
	}

	return trans
}

func (s *Server) initDB() {
	pool, err := pgxpool.New(context.Background(), s.config.Sql.GetDataSource())
	if err != nil {
		err := &models.InternalError{
			Err:  err,
			Msg:  "failed to init db pool",
			Path: "user.server.initDB",
		}
		s.errors <- err
	}
	s.dbConn = pool
}

func (s *Server) initStore() {
	store := dbstore.NewDBStore(s.dbConn)
	s.dbStore = store
}

func (s *Server) initMailer() {
	m := mailer.NewMailer(&mailer.MailerArgs{ConfigFn: s.configFn, Store: s.dbStore})
	s.mailer = m
}

func (s *Server) initWorker() {
	options := &asynq.RedisClientOpt{
		Addr:         s.config.GetCache().GetRedisAddress(),
		Password:     s.config.Cache.GetRedisPassword(),
		DB:           int(s.config.Cache.GetRedisDb()),
		DialTimeout:  time.Second * 10, // default 5
		ReadTimeout:  time.Second * 5,  // default 3
		WriteTimeout: time.Second * 5,  // default is write timeout
	}

	tasker := worker.NewAsynqTaksDistributor(&worker.TaskDistributorArgs{
		Log:     s.log,
		Config:  s.configFn,
		Options: options,
	})

	w := worker.NewAsynqTaskProcessor(&worker.TaskProcessorArgs{
		Store:   s.dbStore,
		Config:  s.configFn,
		Mailer:  s.mailer,
		Log:     s.log,
		Options: options,
	})

	s.tasker = tasker
	go func() {
		err := w.Start()
		if err != nil {
			s.errors <- &models.InternalError{
				Err:  err,
				Msg:  "failed to start worker",
				Path: "user.server.initWorker",
			}
		}
	}()
}
