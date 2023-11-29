package btcslasher

import (
	"fmt"

	bbn "github.com/babylonchain/babylon/types"
)

// isTaprootOutputSpendable checks if the taproot output of a given tx is still spendable on Bitcoin
// This function can be used to check the output of a staking tx or an undelegation tx
func (bs *BTCSlasher) isTaprootOutputSpendable(txBytes []byte, outIdx uint32) (bool, error) {
	stakingMsgTx, err := bbn.NewBTCTxFromBytes(txBytes)
	if err != nil {
		return false, fmt.Errorf(
			"failed to convert staking tx to MsgTx: %v",
			err,
		)
	}
	stakingMsgTxHash := stakingMsgTx.TxHash()

	// we make use of GetTxOut, which returns a non-nil UTXO if it's spendable
	// see https://developer.bitcoin.org/reference/rpc/gettxout.html for details
	// NOTE: we also consider mempool tx as per the last parameter
	txOut, err := bs.BTCClient.GetTxOut(&stakingMsgTxHash, outIdx, true)
	if err != nil {
		return false, fmt.Errorf(
			"failed to get the output of tx %s: %v",
			stakingMsgTxHash.String(),
			err,
		)
	}
	if txOut == nil {
		log.Debugf(
			"tx %s output is already unspendable",
			stakingMsgTxHash.String(),
		)
		return false, nil
	}
	// spendable
	return true, nil
}
