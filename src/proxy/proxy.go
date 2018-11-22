package proxy

import (
	"github.com/andrecronje/lachesis/src/poset"
	"github.com/andrecronje/lachesis/src/proxy/proto"
)

// AppProxy provides an interface for lachesis to communicate
// with the application.
type AppProxy interface {
	SubmitCh() chan []byte
	SubmitInternalCh() chan poset.InternalTransaction
	CommitBlock(block poset.Block) (CommitResponse, error)
	GetSnapshot(blockIndex int64) ([]byte, error)
	Restore(snapshot []byte) error
}

// LachesisProxy provides an interface for the application to
// submit transactions to the lachesis node.
type LachesisProxy interface {
	CommitCh() chan proto.Commit
	SnapshotRequestCh() chan proto.SnapshotRequest
	RestoreCh() chan proto.RestoreRequest
	SubmitTx(tx []byte) error
}
