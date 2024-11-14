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

// ➜  goggle git:(main) ✗ gow run ./cmd/pocsaga (success case)
// 2024/11/14 12:21:18 task_c/fill_balance: doing for 13 seconds...
// 2024/11/14 12:21:18 task_d/announce: doing for 2 seconds...
// 2024/11/14 12:21:18 task_a/create_user: doing for 12 seconds... (has commit step)
// 2024/11/14 12:21:18 task_b/create_address: doing for 14 seconds... (has commit step)
// 2024/11/14 12:21:20 task_d/announce: done
// 2024/11/14 12:21:30 task_a/create_user: done
// 2024/11/14 12:21:31 task_c/fill_balance: done
// 2024/11/14 12:21:32 task_b/create_address: done
// 2024/11/14 12:21:32 committed task_a/create_user
// 2024/11/14 12:21:32 committed task_b/create_address

// ➜  goggle git:(main) ✗ gow run ./cmd/pocsaga (error case)
// 2024/11/14 12:19:53 task_b/create_address: doing for 2 seconds... (has commit step)
// 2024/11/14 12:19:53 task_a/create_user: doing for 12 seconds... (has commit step)
// 2024/11/14 12:19:53 task_d/announce: doing for 14 seconds...
// 2024/11/14 12:19:53 task_c/fill_balance: doing for 13 seconds...
// 2024/11/14 12:19:55 task_b/create_address: failed!
// 2024/11/14 12:20:05 task_a/create_user: done
// 2024/11/14 12:20:06 task_c/fill_balance: failed!
// 2024/11/14 12:20:07 task_d/announce: done
// 2024/11/14 12:20:07 parallel execution failed: 2 joint error: task_b/create_address failed; task_c/fill_balance failed
// 2024/11/14 12:20:07 rolledback task_d/announce
// 2024/11/14 12:20:07 rolledback task_a/create_user

func main() {
	ctx := context.Background()

	ctx = ctxboard.CreateWith(ctx)
	saga := ctxsaga.CreateSaga(ctx)

	errs := do.Parallel(ctx, []func() error{
		func() (err error) { return Do(ctx, "task_a/create_user", true, true) },
		func() (err error) { return Do(ctx, "task_b/create_address", true, true) },
		func() (err error) { return Do(ctx, "task_c/fill_balance", false, true) },
		func() (err error) { return Do(ctx, "task_d/announce", false, true) },
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
func Do(ctx context.Context, name string, hasCommitFx bool, shouldPass bool) (err error) {
	saga, _ := ctxsaga.GetSaga(ctx)

	dur := rand.Intn(20) + 1

	log.Printf("%s: doing for %d seconds... %s\n", name, dur, func() string {
		if hasCommitFx {
			return "(has commit step)"
		}
		return ""
	}())
	time.Sleep(time.Duration(dur) * time.Second)
	if !shouldPass {
		log.Printf("%s: failed!\n", name)
		err = errors.Errorf("%s failed", name)
		return
	}

	log.Printf("%s: done\n", name)

	if hasCommitFx {
		saga.AddCommitFx(func() (err error) { log.Println("committed", name); return })
	}
	saga.AddRollbackFx(func() (err error) { log.Println("rolledback", name); return })

	return
}
