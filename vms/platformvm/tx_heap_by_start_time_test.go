// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/vms/platformvm/transactions/txheap"
	"github.com/ava-labs/avalanchego/vms/platformvm/transactions/unsigned"
)

// TODO: while heap has been moved to a different package, these UTs are still here since
// tx factory is still in platformvm package. Will be cleaned up in a subsequent PR

func TestTxHeapByStartTime(t *testing.T) {
	vm, _, _ := defaultVM()
	vm.ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(); err != nil {
			t.Fatal(err)
		}
		vm.ctx.Lock.Unlock()
	}()

	txHeap := txheap.NewTxHeapByStartTime()

	validator0, err := vm.newAddValidatorTx(
		vm.MinValidatorStake,                                               // stake amount
		uint64(defaultGenesisTime.Unix()+1),                                // startTime
		uint64(defaultGenesisTime.Add(defaultMinStakingDuration).Unix()+1), // endTime
		ids.NodeID{},                            // node ID
		ids.ShortID{1, 2, 3, 4, 5, 6, 7},        // reward address
		0,                                       // shares
		[]*crypto.PrivateKeySECP256K1R{keys[0]}, // key
		ids.ShortEmpty,                          // change addr
	)
	if err != nil {
		t.Fatal(err)
	}
	vdr0Tx := validator0.Unsigned.(*unsigned.AddValidatorTx)

	validator1, err := vm.newAddValidatorTx(
		vm.MinValidatorStake,                                               // stake amount
		uint64(defaultGenesisTime.Unix()+2),                                // startTime
		uint64(defaultGenesisTime.Add(defaultMinStakingDuration).Unix()+2), // endTime
		ids.NodeID{1},                           // node ID
		ids.ShortID{1, 2, 3, 4, 5, 6, 7},        // reward address
		0,                                       // shares
		[]*crypto.PrivateKeySECP256K1R{keys[0]}, // key
		ids.ShortEmpty,                          // change addr
	)
	if err != nil {
		t.Fatal(err)
	}
	vdr1Tx := validator1.Unsigned.(*unsigned.AddValidatorTx)

	validator2, err := vm.newAddValidatorTx(
		vm.MinValidatorStake,                                               // stake amount
		uint64(defaultGenesisTime.Unix()+3),                                // startTime
		uint64(defaultGenesisTime.Add(defaultMinStakingDuration).Unix()+3), // endTime
		ids.NodeID{},                            // node ID
		ids.ShortID{1, 2, 3, 4, 5, 6, 7},        // reward address
		0,                                       // shares
		[]*crypto.PrivateKeySECP256K1R{keys[0]}, // key
		ids.ShortEmpty,                          // change addr
	)
	if err != nil {
		t.Fatal(err)
	}
	vdr2Tx := validator2.Unsigned.(*unsigned.AddValidatorTx)

	txHeap.Add(validator2)
	if timestamp := txHeap.Timestamp(); !timestamp.Equal(vdr2Tx.StartTime()) {
		t.Fatalf("TxHeap.Timestamp returned %s, expected %s", timestamp, vdr2Tx.StartTime())
	}

	txHeap.Add(validator1)
	if timestamp := txHeap.Timestamp(); !timestamp.Equal(vdr1Tx.StartTime()) {
		t.Fatalf("TxHeap.Timestamp returned %s, expected %s", timestamp, vdr1Tx.StartTime())
	}

	txHeap.Add(validator0)
	if timestamp := txHeap.Timestamp(); !timestamp.Equal(vdr0Tx.StartTime()) {
		t.Fatalf("TxHeap.Timestamp returned %s, expected %s", timestamp, vdr0Tx.StartTime())
	} else if top := txHeap.Peek(); top.Unsigned.ID() != validator0.Unsigned.ID() {
		t.Fatalf("TxHeap prioritized %s, expected %s", top.Unsigned.ID(), validator0.Unsigned.ID())
	}
}
