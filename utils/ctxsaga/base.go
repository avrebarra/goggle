package ctxsaga

import (
	"context"
	"sync"

	"github.com/avrebarra/goggle/utils/ctxboard"
)

var keysaga = "!ctxwall/saga"

func CreateSaga(ctx context.Context) *SagaCenter {
	sc := SagaCenter{rollbackfxs: &sync.Map{}, commitfxs: &sync.Map{}}
	ctxboard.SetData(ctx, keysaga, &sc)
	return &sc
}

func GetSaga(ctx context.Context) *SagaCenter {
	sc := ctxboard.GetData(ctx, keysaga)
	if sc != nil {
		return sc.(*SagaCenter)
	}
	CreateSaga(ctx)
	return GetSaga(ctx)
}

// ***

type SagaFunc func() error

type SagaCenter struct {
	commitfxs   *sync.Map
	rollbackfxs *sync.Map
}

func (s *SagaCenter) AddRollbackFx(fx SagaFunc) {
	count := countmap(s.rollbackfxs)
	s.rollbackfxs.Store(count, fx)
}

func (s *SagaCenter) AddCommitFx(fx SagaFunc) {
	count := countmap(s.commitfxs)
	s.commitfxs.Store(count, fx)
}

func (s *SagaCenter) Commit() (err error) {
	targmap := s.commitfxs
	for i := 0; i < countmap(targmap); i++ {
		if fx, ok := targmap.Load(i); ok {
			if err = fx.(SagaFunc)(); err != nil {
				return
			}
		}
	}
	return
}

func (s *SagaCenter) Rollback() (err error) {
	targmap := s.rollbackfxs
	for i := countmap(targmap) - 1; i >= 0; i-- {
		if fx, ok := targmap.Load(i); ok {
			if err = fx.(SagaFunc)(); err != nil {
				return
			}
		}
	}
	return
}

// ***

func countmap(m *sync.Map) (out int) {
	m.Range(func(_, _ any) bool { out++; return true })
	return
}
