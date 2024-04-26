// Package worker implements mining, peer updates, and transaction sharing for
// the blockchain.
package worker

import (
	"sync"

	"github.com/ardanlabs/blockchain/foundation/blockchain/state"
)

// Worker manages the POW workflows for the blockchain.
type Worker struct {
	state *state.State
	wg    sync.WaitGroup
	// ticker      time.Ticker
	//we are signaling by only closing chanl
	shut chan struct{} //if the only think that we are gonna do in chanel is closing it is bettter to use empty struct
	//it is better to use bool when we are signaling with data BUT data is arbitrary(data itself is irrelevant)
	startMining  chan bool
	cancelMining chan bool
	// txSharing chan database.BlockTx
	evHandler state.EventHandler
}

// Run creates a worker, registers the worker with the state package, and
// starts up all the background processes.
func Run(st *state.State, evHandler state.EventHandler) {

	w := Worker{
		state: st,
		// ticker:       *time.NewTicker(peerUpdateInterval),
		shut:         make(chan struct{}),
		startMining:  make(chan bool, 1),
		cancelMining: make(chan bool, 1), //make(chan bool, 1),
		// txSharing:    make(chan database.BlockTx, maxTxShareRequests),
		evHandler: evHandler,
	}

	w.evHandler("worker: SignalStartMining: mining signaled")
	// Update this node before starting any support G's.
	// w.Sync()

	// Select the consensus operation to run.
	// consensusOperation := w.powOperations
	// if st.Consensus() == state.ConsensusPOA {
	// 	consensusOperation = w.poaOperations
	// }

	// Register this worker with the state package.
	st.Worker = &w

	// Load the set of operations we need to run.
	operations := []func(){
		// w.peerOperations,
		// w.shareTxOperations,
		// consensusOperation,
		w.miningOperations,
	}

	// Set waitgroup to match the number of G's we need for the set
	// of operations we have.
	g := len(operations) //whoever owns worker goroutine in our case it is main goroutine becames parent of all this goroutines
	w.wg.Add(g)          //this gives us ability to wait in main untill all it's child are shutdown

	// We don't want to return until we know all the G's are up and running.
	hasStarted := make(chan bool)

	// Start all the operational G's.
	// before handling transactions and minig we need to ensure that all goroutines that essential for this operations are running
	for _, op := range operations {
		go func(op func()) {
			defer w.wg.Done()
			hasStarted <- true
			op()
		}(op)
	}

	// Wait for the G's to report they are running.
	for i := 0; i < g; i++ {
		<-hasStarted
	}

}

// =============================================================================
// These methods implement the state.Worker interface.

// Shutdown terminates the goroutine performing work.
func (w *Worker) Shutdown() {
	w.evHandler("worker: shutdown: started")
	defer w.evHandler("worker: shutdown: completed")

	// w.evHandler("worker: shutdown: stop ticker")
	// w.ticker.Stop()

	// w.evHandler("worker: shutdown: signal cancel mining")
	// w.SignalCancelMining()

	w.evHandler("worker: shutdown: signal cancel mining")

	w.evHandler("worker: shutdown: terminate goroutines")
	close(w.shut)
	w.wg.Wait()
}

// SignalStartMining starts a mining operation. If there is already a signal
// pending in the channel, just return since a mining operation will start.
func (w *Worker) SignalStartMining() {
	// if !w.state.IsMiningAllowed() {
	// 	w.evHandler("state: MinePeerBlock: accepting blocks turned off")
	// 	return
	// }

	// Only PoW requires signalling to start mining.
	// if w.state.Consensus() != state.ConsensusPOW {
	// 	return
	// }

	select {
	case w.startMining <- true:
	default: //it means do not block because there will a lot of other chanel than want to send signal but we can't block them
	}
	w.evHandler("worker: SignalStartMining: mining signaled")
}

// used to test if shutdowwn has been signaled
func (w *Worker) isShutdown() bool {
	select {
	case <-w.shut:
		return true
	default:
		return false
	}

}

// SignalCancelMining signals the G executing the runMiningOperation function
// to stop immediately.
func (w *Worker) SignalCancelMining() {

	select {
	case w.cancelMining <- true:
	default:
	}
	w.evHandler("worker: SignalCancelMining: MINING: CANCEL: signaled")

}
