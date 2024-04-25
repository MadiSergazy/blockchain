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
