package state

import (
	"fmt"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
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
