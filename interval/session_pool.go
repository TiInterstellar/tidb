package interval

import (
	"sync"
	"time"

	"github.com/ngaut/pools"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/util/logutil"
	"github.com/pingcap/tidb/util/mock"
)

// sessionPool is used to new session.
type sessionPool struct {
	mu struct {
		sync.Mutex
		closed bool
	}
	resPool *pools.ResourcePool
}

func newSessionPool(resPool *pools.ResourcePool) *sessionPool {
	return &sessionPool{resPool: resPool}
}

// get gets sessionctx from context resource pool.
// Please remember to call put after you finished using sessionctx.
func (sg *sessionPool) get() (sessionctx.Context, error) {
	if sg.resPool == nil {
		return mock.NewContext(), nil
	}

	sg.mu.Lock()
	if sg.mu.closed {
		sg.mu.Unlock()
		return nil, errors.Errorf("sessionPool is closed.")
	}
	sg.mu.Unlock()

	// no need to protect sg.resPool
	resource, err := sg.resPool.Get()
	if err != nil {
		return nil, errors.Trace(err)
	}

	ctx := resource.(sessionctx.Context)
	ctx.GetSessionVars().SetStatusFlag(mysql.ServerStatusAutocommit, true)
	ctx.GetSessionVars().InRestrictedSQL = true
	ctx.GetSessionVars().TimeZone = time.Local
	ctx.GetSessionVars().StmtCtx.TimeZone = time.Local
	return ctx, nil
}

// put returns sessionctx to context resource pool.
func (sg *sessionPool) put(ctx sessionctx.Context) {
	if sg.resPool == nil {
		return
	}

	// no need to protect sg.resPool, even the sg.resPool is closed, the ctx still need to
	// put into resPool, because when resPool is closing, it will wait all the ctx returns, then resPool finish closing.
	sg.resPool.Put(ctx.(pools.Resource))
}

// close clean up the sessionPool.
func (sg *sessionPool) close() {
	sg.mu.Lock()
	defer sg.mu.Unlock()
	// prevent closing resPool twice.
	if sg.mu.closed || sg.resPool == nil {
		return
	}
	logutil.BgLogger().Info("[ddl] closing sessionPool")
	sg.resPool.Close()
	sg.mu.closed = true
}
