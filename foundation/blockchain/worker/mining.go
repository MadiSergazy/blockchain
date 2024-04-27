package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ardanlabs/blockchain/foundation/blockchain/state"
)

// this function is keeping goroutines alive
func (w *Worker) miningOperations() {
	w.evHandler("Worker mining operation G started")
	defer w.evHandler("Worker mining operation G is complited")

	for {
		select {
		case <-w.startMining:
			if !w.isShutdown() {
				w.runMiningOperaion()
			}
		case <-w.shut:
			w.evHandler("Worker mining operations: received shut signal")
			return
		}
	}
}

func (w *Worker) runMiningOperaion() {
	w.evHandler("Worker runMining operations: MINING: started")
	defer w.evHandler("Worker runMining operations: MINING: completed")

	length := w.state.QueryMempoolLength()
	if length == 0 {
		w.evHandler("Worker runMining operations: MINING: no transactions to mine Txs[%d]", length)
		return
	}

	//after mining operations check if a new operations shoud be signaled again
	defer func() {
		length := w.state.QueryMempoolLength()
		if length > 0 {
			w.evHandler("Worker runMining operations: MINING: signal new mining opertion: Txs[%d]", length)
			w.SignalStartMining()
		}
	}()

	select {
	case <-w.cancelMining:
		w.evHandler("Worker: runMining operations: MAINING: drained cancel")

	default:
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	//goroutine for handling cancel
	//G canceling mining operation
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		select {
		case <-w.cancelMining:
			w.evHandler("Worker: runMinigOperation: MINIMG: CANCEL: requested")

		case <-ctx.Done():
		}
	}()

	//G perfoeming minig operation
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		t := time.Now()
		_, err := w.state.MineNewBlock(ctx)
		duration := time.Since(t)

		w.evHandler("Worker: runMiningOperation: MINING: mining duration:[%v]", duration)

		if err != nil {
			switch {
			case errors.Is(err, state.ErrNoTransactions):
				w.evHandler("Worker: runMiningOperation: MINING: WARNING: No transactions in mempool")
			case ctx.Err() != nil:
				w.evHandler("Worker: runMiningOperation: MINING: CANCELL: complete")

			default:
				w.evHandler("Worker: runMiningOperation: MINING: ERROR: %s", err)

			}
			return
		}
	}()
	wg.Wait()
}
