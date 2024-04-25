// Package state is the core API for the blockchain and implements all the
// business rules and processing.
package state

import (
	"sync"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ardanlabs/blockchain/foundation/blockchain/genesis"
	"github.com/ardanlabs/blockchain/foundation/blockchain/mempool"
	"github.com/ardanlabs/blockchain/foundation/blockchain/storage/disk"
)

// =============================================================================

// EventHandler defines a function that is called when events
// occur in the processing of persisting blocks.
type EventHandler func(v string, args ...any)

// =============================================================================

// State manages the blockchain database.
type State struct {
	mu sync.RWMutex
	// resyncWG    sync.WaitGroup
	// allowMining bool

	beneficiaryID database.AccountID
	host          string
	evHandler     EventHandler
	consensus     string

	// knownPeers *peer.PeerSet
	// storage database.Storage
	genesis genesis.Genesis
	mempool *mempool.Mempool
	db      *database.Database

	// Worker Worker
}

// =============================================================================

// Config represents the configuration required to start
// the blockchain node.
type Config struct {
	BeneficiaryID database.AccountID
	Host          string
	DbPath        string
	// Storage        database.Storage
	// Genesis        genesis.Genesis
	// SelectStrategy string
	// KnownPeers     *peer.PeerSet
	EvHandler EventHandler
	// Consensus string
}

// New constructs a new blockchain for data management.
func New(cfg Config) (*State, error) {

	// Build a safe event handler function for use.
	ev := func(v string, args ...any) {
		if cfg.EvHandler != nil {
			cfg.EvHandler(v, args...)
		}
	}

	genesis, err := genesis.Load()
	if err != nil {
		return nil, err
	}

	storage, err := disk.New(cfg.DbPath)
	if err != nil {
		return nil, err
	}

	// Access the storage for the blockchain.
	db, err := database.New(genesis, storage, ev)
	if err != nil {
		return nil, err
	}

	// Construct a mempool with the specified sort strategy.
	mempool, err := mempool.New()
	if err != nil {
		return nil, err
	}

	// Create the State to provide support for managing the blockchain.
	state := State{
		beneficiaryID: cfg.BeneficiaryID,
		host:          cfg.Host,
		// storage:       cfg.Storage,
		evHandler: ev,
		// consensus:     cfg.Consensus,
		// allowMining:   true,

		// knownPeers: cfg.KnownPeers,
		genesis: genesis,
		mempool: mempool,
		db:      db,
	}

	// The Worker is not set here. The call to worker.Run will assign itself
	// and start everything up and running for the node.

	return &state, nil

}
