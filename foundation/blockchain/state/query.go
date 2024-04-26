package state

func (s *State) QueryMempoolLength() int {
	return s.mempool.Count()
}
