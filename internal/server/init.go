package server

import (
	"context"

	commonPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store/dbstore"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *Server) initTrans() map[string]*commonPb.TranslationElements {
	trans, err := a.commonClient.TranslationsGet()
	if err != nil {
		a.done <- err
	}

	if err := models.TranslationsInit(trans); err != nil {
		err := &models.InternalError{
			Msg:  "failed to init translations",
			Err:  err,
			Path: "user.server.initTrans",
		}
		a.done <- err
	}

	return trans
}

func (s *Server) initDB() {
	pool, err := pgxpool.New(context.Background(), s.config.Sql.GetDataSource())
	if err != nil {
		err := &models.InternalError{
			Msg:  "failed to init db pool",
			Err:  err,
			Path: "user.server.initDB",
		}
		s.done <- err
	}
	s.db = pool
}

func (s *Server) initStore() {
	store := dbstore.NewDBStore(s.db)
	s.dbStore = store
}
