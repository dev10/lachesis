package peers

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"crypto/ecdsa"

	"reflect"

	scrypto "github.com/andrecronje/lachesis/src/crypto"
)

func TestJSONPeerSet(t *testing.T) {
	// Create a test dir
	dir, err := ioutil.TempDir("", "babble")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(dir)

	// Create the store
	store := NewJSONPeerSet(dir)

	// Try a read, should get nothing
	peerSet, err := store.PeerSet()
	if err == nil {
		t.Fatalf("store.PeerSet() should generate an error")
	}
	if peerSet != nil {
		t.Fatalf("peerSet: %v", peerSet)
	}

	keys := map[string]*ecdsa.PrivateKey{}
	peers := []*Peer{}
	for i := 0; i < 3; i++ {
		key, _ := scrypto.GenerateECDSAKey()
		peer := &Peer{
			NetAddr:   fmt.Sprintf("addr%d", i),
			PubKeyHex: fmt.Sprintf("0x%X", scrypto.FromECDSAPub(&key.PublicKey)),
		}
		peers = append(peers, peer)
		keys[peer.NetAddr] = key
	}

	newPeerSet := NewPeerSet(peers)
	newPeerSlice := newPeerSet.Peers

	if err := store.Write(newPeerSlice); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Try a read, should find 3 peers
	peerSet, err = store.PeerSet()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if peerSet.Len() != 3 {
		t.Fatalf("peers: %v", peers)
	}

	peerSlice := peerSet.Peers

	for i := 0; i < 3; i++ {
		if peerSlice[i].NetAddr != newPeerSlice[i].NetAddr {
			t.Fatalf("peers[%d] NetAddr should be %s, not %s", i,
				newPeerSlice[i].NetAddr, peerSlice[i].NetAddr)
		}
		if peerSlice[i].PubKeyHex != newPeerSlice[i].PubKeyHex {
			t.Fatalf("peers[%d] PubKeyHex should be %s, not %s", i,
				newPeerSlice[i].PubKeyHex, peerSlice[i].PubKeyHex)
		}
		pubKeyBytes, err := peerSlice[i].PubKeyBytes()
		if err != nil {
			t.Fatal(err)
		}
		pubKey := scrypto.ToECDSAPub(pubKeyBytes)
		if !reflect.DeepEqual(*pubKey, keys[peerSlice[i].NetAddr].PublicKey) {
			t.Fatalf("peers[%d] PublicKey not parsed correctly", i)
		}
	}
}