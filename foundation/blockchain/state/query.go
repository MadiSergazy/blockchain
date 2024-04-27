package state

import (
	"errors"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
)

func (s *State) QueryMempoolLength() int {
	return s.mempool.Count()
}

func (s *State) QueryAccounts(account database.AccountID) (database.Account, error) {
	accounts := s.db.CopyAccounts()

	if info, exists := accounts[account]; exists {
		return info, nil
	}
	return database.Account{}, errors.New("not found")
}
