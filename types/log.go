package types

import (
	"fmt"
	"math/big"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Log represents a log entry from a transaction receipt.
type Log struct {
	Address          common.Address `json:"address"`
	Topics           []common.Hash  `json:"topics"`
	Data             []byte         `json:"data"`
	BlockNumber      uint64         `json:"blockNumber"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	TransactionIndex uint64         `json:"transactionIndex"`
	BlockHash        common.Hash    `json:"blockHash"`
	LogIndex         uint64         `json:"logIndex"`
	Removed          bool           `json:"removed"`
}

// UnmarshalJSON implements json.Unmarshaler for Log.
func (l *Log) UnmarshalJSON(data []byte) error {
	type logJSON struct {
		Address          common.Address `json:"address"`
		Topics           []common.Hash  `json:"topics"`
		Data             string         `json:"data"`
		BlockNumber      string         `json:"blockNumber"`
		TransactionHash  common.Hash    `json:"transactionHash"`
		TransactionIndex string         `json:"transactionIndex"`
		BlockHash        common.Hash    `json:"blockHash"`
		LogIndex         string         `json:"logIndex"`
		Removed          bool           `json:"removed"`
	}

	var raw logJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.Address = raw.Address
	l.Topics = raw.Topics
	l.TransactionHash = raw.TransactionHash
	l.BlockHash = raw.BlockHash
	l.Removed = raw.Removed

	if raw.Data != "" {
		d, err := hexutil.Decode(raw.Data)
		if err != nil {
			return fmt.Errorf("invalid log data: %w", err)
		}
		l.Data = d
	}

	if raw.BlockNumber != "" {
		bn, err := hexutil.DecodeUint64(raw.BlockNumber)
		if err != nil {
			return fmt.Errorf("invalid block number: %w", err)
		}
		l.BlockNumber = bn
	}

	if raw.TransactionIndex != "" {
		ti, err := hexutil.DecodeUint64(raw.TransactionIndex)
		if err != nil {
			return fmt.Errorf("invalid transaction index: %w", err)
		}
		l.TransactionIndex = ti
	}

	if raw.LogIndex != "" {
		li, err := hexutil.DecodeUint64(raw.LogIndex)
		if err != nil {
			return fmt.Errorf("invalid log index: %w", err)
		}
		l.LogIndex = li
	}

	return nil
}

// Receipt represents a transaction receipt.
type Receipt struct {
	TransactionHash   common.Hash     `json:"transactionHash"`
	TransactionIndex  uint64          `json:"transactionIndex"`
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       uint64          `json:"blockNumber"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	GasUsed           uint64          `json:"gasUsed"`
	ContractAddress   *common.Address `json:"contractAddress"`
	Logs              []Log           `json:"logs"`
	Status            uint64          `json:"status"`
	LogsBloom         []byte          `json:"logsBloom"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"`
	Type              uint64          `json:"type"`
	// EIP-4844 fields
	BlobGasUsed  *uint64  `json:"blobGasUsed,omitempty"`
	BlobGasPrice *big.Int `json:"blobGasPrice,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler for Receipt.
func (r *Receipt) UnmarshalJSON(data []byte) error {
	type receiptJSON struct {
		TransactionHash   common.Hash     `json:"transactionHash"`
		TransactionIndex  string          `json:"transactionIndex"`
		BlockHash         common.Hash     `json:"blockHash"`
		BlockNumber       string          `json:"blockNumber"`
		From              common.Address  `json:"from"`
		To                *common.Address `json:"to"`
		CumulativeGasUsed string          `json:"cumulativeGasUsed"`
		GasUsed           string          `json:"gasUsed"`
		ContractAddress   *common.Address `json:"contractAddress"`
		Logs              []Log           `json:"logs"`
		Status            string          `json:"status"`
		LogsBloom         string          `json:"logsBloom"`
		EffectiveGasPrice string          `json:"effectiveGasPrice"`
		Type              string          `json:"type"`
		BlobGasUsed       string          `json:"blobGasUsed,omitempty"`
		BlobGasPrice      string          `json:"blobGasPrice,omitempty"`
	}

	var raw receiptJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.TransactionHash = raw.TransactionHash
	r.BlockHash = raw.BlockHash
	r.From = raw.From
	r.To = raw.To
	r.ContractAddress = raw.ContractAddress
	r.Logs = raw.Logs

	if raw.TransactionIndex != "" {
		idx, err := hexutil.DecodeUint64(raw.TransactionIndex)
		if err != nil {
			return fmt.Errorf("invalid transaction index: %w", err)
		}
		r.TransactionIndex = idx
	}

	if raw.BlockNumber != "" {
		bn, err := hexutil.DecodeUint64(raw.BlockNumber)
		if err != nil {
			return fmt.Errorf("invalid block number: %w", err)
		}
		r.BlockNumber = bn
	}

	if raw.CumulativeGasUsed != "" {
		cgu, err := hexutil.DecodeUint64(raw.CumulativeGasUsed)
		if err != nil {
			return fmt.Errorf("invalid cumulative gas used: %w", err)
		}
		r.CumulativeGasUsed = cgu
	}

	if raw.GasUsed != "" {
		gu, err := hexutil.DecodeUint64(raw.GasUsed)
		if err != nil {
			return fmt.Errorf("invalid gas used: %w", err)
		}
		r.GasUsed = gu
	}

	if raw.Status != "" {
		status, err := hexutil.DecodeUint64(raw.Status)
		if err != nil {
			return fmt.Errorf("invalid status: %w", err)
		}
		r.Status = status
	}

	if raw.LogsBloom != "" {
		bloom, err := hexutil.Decode(raw.LogsBloom)
		if err != nil {
			return fmt.Errorf("invalid logs bloom: %w", err)
		}
		r.LogsBloom = bloom
	}

	if raw.EffectiveGasPrice != "" {
		egp, err := hexutil.DecodeBig(raw.EffectiveGasPrice)
		if err != nil {
			return fmt.Errorf("invalid effective gas price: %w", err)
		}
		r.EffectiveGasPrice = egp
	}

	if raw.Type != "" {
		t, err := hexutil.DecodeUint64(raw.Type)
		if err != nil {
			return fmt.Errorf("invalid type: %w", err)
		}
		r.Type = t
	}

	if raw.BlobGasUsed != "" {
		bgu, err := hexutil.DecodeUint64(raw.BlobGasUsed)
		if err != nil {
			return fmt.Errorf("invalid blob gas used: %w", err)
		}
		r.BlobGasUsed = &bgu
	}

	if raw.BlobGasPrice != "" {
		bgp, err := hexutil.DecodeBig(raw.BlobGasPrice)
		if err != nil {
			return fmt.Errorf("invalid blob gas price: %w", err)
		}
		r.BlobGasPrice = bgp
	}

	return nil
}

// IsSuccess returns true if the transaction was successful.
func (r *Receipt) IsSuccess() bool {
	return r.Status == 1
}

// IsFailed returns true if the transaction failed.
func (r *Receipt) IsFailed() bool {
	return r.Status == 0
}
