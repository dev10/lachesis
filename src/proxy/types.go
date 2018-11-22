package proxy

import "github.com/Fantom-foundation/go-lachesis/src/poset"

type CommitResponse struct {
	StateHash                    []byte
	AcceptedInternalTransactions []poset.InternalTransaction
}

type CommitCallback func(block poset.Block) (CommitResponse, error)

//DummyCommitCallback is used for testing
func DummyCommitCallback(block poset.Block) (CommitResponse, error) {
	acceptedInternalTransactions := make([]bool, len(block.InternalTransactions()))
	for i := range block.InternalTransactions() {
		acceptedInternalTransactions[i] = true
	}

	var transactions []poset.InternalTransaction
	for _, v := range block.InternalTransactions() {
		transactions = append(transactions, *v)
	}
	res := CommitResponse{
		StateHash:                    []byte{},
		AcceptedInternalTransactions: transactions,
	}

	return res, nil
}
