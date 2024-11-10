package main

import (
	"context"
	"log"
	"time"

	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/avrebarra/goggle/utils/do"
	"github.com/pkg/errors"
	"golang.org/x/exp/rand"
)

// 2024/11/09 20:46:02 task_d/announce: doing for 2 seconds...
// 2024/11/09 20:46:02 task_a/create_user: doing for 2 seconds...
// 2024/11/09 20:46:02 task_b/create_address: doing for 3 seconds...
// 2024/11/09 20:46:02 task_c/fill_balance: doing for 4 seconds...
// 2024/11/09 20:46:04 task_d/announce: failed!
// 2024/11/09 20:46:04 task_a/create_user: failed!
// 2024/11/09 20:46:05 task_b/create_address: failed!
// 2024/11/09 20:46:06 task_c/fill_balance: failed!
// 2024/11/09 20:46:06 parallel execution failed: 4 joint error: task_a/create_user failed; task_b/create_address failed; task_c/fill_balance failed; task_d/announce failed

func main() {
	ctx := context.Background()

	ctx = ctxboard.CreateWith(ctx)
	saga := ctxsaga.CreateSaga(ctx)

	errs := do.Parallel(ctx, []func() error{
		func() (err error) { return DoBasic(ctx, "task_a/create_user", true) },
		func() (err error) { return DoWithCommit(ctx, "task_b/create_address", true) },
		func() (err error) { return DoBasic(ctx, "task_c/fill_balance", true) },
		func() (err error) { return DoWithCommit(ctx, "task_d/announce", true) },
	})
	err := do.JoinErrors(errs)
	shouldRollback := false
	if err != nil {
		err = errors.Wrap(err, "parallel execution failed")
		log.Println(err.Error())
		shouldRollback = true
		err = nil // discard
	}
	if shouldRollback {
		if err = saga.Rollback(); err != nil {
			err = errors.Wrap(err, "failed to rollback")
			log.Println(err.Error())
			return
		}
		return
	}

	err = saga.Commit()
	if err != nil {
		err = errors.Wrap(err, "failed to commit")
		log.Println(err.Error())
		return
	}
}

func DoBasic(ctx context.Context, name string, shouldPass bool) (err error) {
	saga, _ := ctxsaga.GetSaga(ctx)

	dur := rand.Intn(5) + 1

	log.Printf("%s: doing for %d seconds...\n", name, dur)
	time.Sleep(time.Duration(dur) * time.Second)
	if !shouldPass {
		log.Printf("%s: failed!\n", name)
		err = errors.Errorf("%s failed", name)
		return
	}

	log.Printf("%s: done\n", name)

	saga.AddRollbackFx(func() (err error) { log.Println("rolledback", name); return })

	return
}

func DoWithCommit(ctx context.Context, name string, shouldPass bool) (err error) {
	saga, _ := ctxsaga.GetSaga(ctx)

	dur := rand.Intn(5) + 1

	log.Printf("%s: doing for %d seconds...\n", name, dur)
	time.Sleep(time.Duration(dur) * time.Second)
	if !shouldPass {
		log.Printf("%s: failed!\n", name)
		err = errors.Errorf("%s failed", name)
		return
	}

	log.Printf("%s: done\n", name)

	saga.AddRollbackFx(func() (err error) { log.Println("rolledback", name); return })
	saga.AddCommitFx(func() (err error) { log.Println("committed", name); return })

	return
}
