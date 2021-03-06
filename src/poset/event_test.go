package poset

import (
	"reflect"
	"testing"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
	"github.com/Fantom-foundation/go-lachesis/src/inter/wire"
)

func createDummyEventBody() EventBody {
	body := EventBody{}
	body.Transactions = [][]byte{[]byte("abc"), []byte("def")}
	body.InternalTransactions = []*wire.InternalTransaction{}
	body.Parents = [][]byte{[]byte("self"), []byte("other")}
	body.Creator = []byte("public key")
	body.BlockSignatures = []*BlockSignature{
		{
			Validator: body.Creator,
			Index:     0,
			Signature: "r|s",
		},
	}
	return body
}

func TestMarshallBody(t *testing.T) {
	body := createDummyEventBody()

	raw, err := body.ProtoMarshal()
	if err != nil {
		t.Fatalf("Error marshalling EventBody: %s", err)
	}

	newBody := new(EventBody)
	if err := newBody.ProtoUnmarshal(raw); err != nil {
		t.Fatalf("Error unmarshalling EventBody: %s", err)
	}

	if !reflect.DeepEqual(body.Transactions, newBody.Transactions) {
		t.Fatalf("Transactions do not match. Expected %#v, got %#v", body.Transactions, newBody.Transactions)
	}
	if !InternalTransactionListEquals(body.InternalTransactions, newBody.InternalTransactions) {
		t.Fatalf("Internal Transactions do not match. Expected %#v, got %#v", body.InternalTransactions, newBody.InternalTransactions)
	}
	if !BlockSignatureListEquals(body.BlockSignatures, newBody.BlockSignatures) {
		t.Fatalf("BlockSignatures do not match. Expected %#v, got %#v", body.BlockSignatures, newBody.BlockSignatures)
	}
	if !reflect.DeepEqual(body.Parents, newBody.Parents) {
		t.Fatalf("Parents do not match. Expected %#v, got %#v", body.Parents, newBody.Parents)
	}
	if !reflect.DeepEqual(body.Creator, newBody.Creator) {
		t.Fatalf("Creators do not match. Expected %#v, got %#v", body.Creator, newBody.Creator)
	}

}

func TestSignEvent(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	publicKeyBytes := privateKey.Public().Bytes()

	body := createDummyEventBody()
	body.Creator = publicKeyBytes

	event := Event{Message: &EventMessage{Body: &body}}
	if err := event.Sign(privateKey); err != nil {
		t.Fatalf("Error signing Event: %s", err)
	}

	res, err := event.Verify()
	if err != nil {
		t.Fatalf("Error verifying signature: %s", err)
	}
	if !res {
		t.Fatalf("Verify returned false")
	}
}

func TestMarshallEvent(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	publicKeyBytes := privateKey.Public().Bytes()

	body := createDummyEventBody()
	body.Creator = publicKeyBytes

	event := Event{Message: &EventMessage{Body: &body}}
	if err := event.Sign(privateKey); err != nil {
		t.Fatalf("Error signing Event: %s", err)
	}

	raw, err := event.ProtoMarshal()
	if err != nil {
		t.Fatalf("Error marshalling Event: %s", err)
	}

	newEvent := new(Event)
	if err := newEvent.ProtoUnmarshal(raw); err != nil {
		t.Fatalf("Error unmarshalling Event: %s", err)
	}

	if !newEvent.Message.Equals(event.Message) {
		t.Fatalf("Events are not deeply equal")
	}
}

func TestWireEvent(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	publicKeyBytes := privateKey.Public().Bytes()

	body := createDummyEventBody()
	body.Creator = publicKeyBytes

	event := Event{Message: &EventMessage{Body: &body}}
	if err := event.Sign(privateKey); err != nil {
		t.Fatalf("Error signing Event: %s", err)
	}

	event.SetWireInfo(1, 66, 2, 67)

	expectedWireEvent := WireEvent{
		Body: WireBody{
			Transactions:         event.Message.Body.Transactions,
			InternalTransactions: inter.WireToInternalTransactions(event.Message.Body.InternalTransactions),
			SelfParentIndex:      1,
			OtherParentCreatorID: 66,
			OtherParentIndex:     2,
			CreatorID:            67,
			Index:                event.Message.Body.Index,
			BlockSignatures:      event.WireBlockSignatures(),
		},
		Signature: event.Message.Signature,
	}

	wireEvent := event.ToWire()

	if !reflect.DeepEqual(expectedWireEvent, wireEvent) {
		t.Fatalf("WireEvent should be %#v, not %#v", expectedWireEvent, wireEvent)
	}
}

func TestIsLoaded(t *testing.T) {
	//nil payload

	event := NewEvent(nil, nil, nil, make(EventHashes, 2), []byte("creator"), 1, nil)
	if event.IsLoaded() {
		t.Fatalf("IsLoaded() should return false for nil Body.Transactions and Body.BlockSignatures")
	}

	//empty payload
	event.Message.Body.Transactions = [][]byte{}
	if event.IsLoaded() {
		t.Fatalf("IsLoaded() should return false for empty Body.Transactions")
	}

	event.Message.Body.BlockSignatures = []*BlockSignature{}
	if event.IsLoaded() {
		t.Fatalf("IsLoaded() should return false for empty Body.BlockSignatures")
	}

	//initial event
	event.Message.Body.Index = 0
	if !event.IsLoaded() {
		t.Fatalf("IsLoaded() should return true for initial event")
	}

	//non-empty tx payload
	event.Message.Body.Transactions = [][]byte{[]byte("abc")}
	if !event.IsLoaded() {
		t.Fatalf("IsLoaded() should return true for non-empty transaction payload")
	}

	//non-empty signature payload
	event.Message.Body.Transactions = nil
	event.Message.Body.BlockSignatures = []*BlockSignature{
		{Validator: []byte("validator"), Index: 0, Signature: "r|s"},
	}
	if !event.IsLoaded() {
		t.Fatalf("IsLoaded() should return true for non-empty signature payload")
	}
}

func TestEventFlagTable(t *testing.T) {
	exp := FlagTable{
		fakeEventHash("x"): 1,
		fakeEventHash("y"): 0,
		fakeEventHash("z"): 2,
	}

	event := NewEvent(nil, nil, nil, make(EventHashes, 2), []byte("creator"), 1, exp)
	if event.IsLoaded() {
		t.Fatalf("IsLoaded() should return false for nil Body.Transactions and Body.BlockSignatures")
	}

	if len(event.Message.FlagTable) == 0 {
		t.Fatal("FlagTable is nil")
	}

	res, err := event.GetFlagTable()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, exp) {
		t.Fatalf("expected flag table: %+v, got: %+v", exp, res)
	}
}

func TestMergeFlagTable(t *testing.T) {
	exp := FlagTable{
		fakeEventHash("x"): 1,
		fakeEventHash("y"): 1,
		fakeEventHash("z"): 1,
	}

	syncData := []FlagTable{
		{
			fakeEventHash("x"): 0,
			fakeEventHash("y"): 1,
			fakeEventHash("z"): 0,
		},
		{
			fakeEventHash("x"): 0,
			fakeEventHash("y"): 0,
			fakeEventHash("z"): 1,
		},
	}

	start := FlagTable{
		fakeEventHash("x"): 1,
		fakeEventHash("y"): 0,
		fakeEventHash("z"): 0,
	}

	ft := start.Marshal()
	event := Event{Message: &EventMessage{FlagTable: ft}}

	for _, v := range syncData {
		flagTable, err := event.MergeFlagTable(v)
		if err != nil {
			t.Fatal(err)
		}
		event.Message.FlagTable = flagTable.Marshal()
	}

	res := FlagTable{}
	err := res.Unmarshal(event.Message.FlagTable)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(exp, res) {
		t.Fatalf("expected flag table: %+v, got: %+v", exp, res)
	}
}

/*
 * stuff
 */

func fakeEventHash(s string) (hash EventHash) {
	hash.Set([]byte(s))
	return
}
