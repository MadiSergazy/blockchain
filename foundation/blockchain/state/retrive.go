package state

import (
	"fmt"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ardanlabs/blockchain/foundation/blockchain/peer"
)

func (s *State) RetriveMempool() []database.BlockTx {

	fmt.Println("mempoool: ", s.mempool)
	fmt.Println("mempoool PickBest: ", s.mempool)

	d := s.mempool.PickBest()

	fmt.Println("D: ", d)
	fmt.Println("mempoool: ", s.mempool)

	return s.mempool.PickBest()
}

func (s *State) RetriveAccounts() map[database.AccountID]database.Account {
	return s.db.CopyAccounts()
}

func (s *State) RetriveALatestBlock() database.Block {
	return s.db.LatestBlock()
}

func (s *State) RetriveKnownPeers() []peer.Peer {
	return s.knownPeers.Copy(s.host)
}

func (s *State) RetriveHost() string {
	return (s.host)
}
