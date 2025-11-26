package server

import (
	"context"
	"net/url"
	"time"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/mailer"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/oauth"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store/dbstore"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *Server) initTrans() map[string]*com.TranslationElements {
	trans, err := s.commonClient.TranslationsGet()
	if err != nil {
		s.errors <- err
	}

	lang := s.config.Localization.GetDefaultClientLocale()
	if err := models.TranslationsInit(trans, lang); err != nil {
		path := "user.server.initTrans"
		err := &models.InternalError{Err: err, Msg: "failed to init translations", Path: path}
		s.errors <- err
	}

	return trans
}

func (s *Server) initDB() {
	pool, err := pgxpool.New(context.Background(), s.config.Sql.GetDataSource())
	if err != nil {
		path := "user.server.initDB"
		err := &models.InternalError{Err: err, Msg: "failed to init db pool", Path: path}
		s.errors <- err
	}
	s.dbConn = pool
}

func (s *Server) initStore() {
	store := dbstore.NewDBStore(s.dbConn)
	s.dbStore = store
}

func (s *Server) initMailer() {
	dir := "./internal/mailer/templates"
	watcher, errCh, err := mailer.NewTemplateContainerWatcher(dir)
	if err != nil {
		s.errors <- &models.InternalError{
			Err:  err,
			Msg:  "failed to initialize the template watcher",
			Path: "user.server.initMailer",
		}
	}

	go func() {
		for e := range errCh {
			s.log.Warnf("templates watcher error: %v", e)
		}
	}()

	m := mailer.NewMailer(&mailer.MailerArgs{ConfigFn: s.configFn, Store: s.dbStore, TemplateContainer: watcher})
	s.mailer = m
}

func (s *Server) initWorker() {
	path := "user.server.initWorker"
	u, err := url.Parse(s.config.GetCache().GetRedisAddress())
	if err != nil {
		s.errors <- &models.InternalError{Err: err, Msg: "failed to parse redis connection URL", Path: path}
	}

	options := &asynq.RedisClientOpt{
		Addr:         u.Host,
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
			s.errors <- &models.InternalError{Err: err, Msg: "failed to start worker", Path: path}
		}
	}()
}

func (s *Server) initOauthServer() {
	path := "user.server.initOauthServer"

	oauth := oauth.NewOauth(oauth.OAuthArgs{
		Config: s.configFn,
		Log:    s.log,
		ErrCh:  make(chan *models.InternalError),
	})

	if err := oauth.Run(); err != nil {
		s.errors <- &models.InternalError{Err: err, Msg: "failed to start oauth server", Path: path}
	}

	go func() {
		for err := range oauth.ErrorChannel() {
			isUnRecoverable := utils.IsUnrecoverableHTTPServerError(err)
			if isUnRecoverable {
				s.errors <- err
			}
		}
	}()
}
