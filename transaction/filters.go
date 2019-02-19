package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/go-pg/pg/orm"
)

// TODO: replace string in StartBlock, EndBlock to int
type BlockFilter struct {
	BlockId         *uint64
}

func(f BlockFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.BlockId != nil {
		q = q.Where("transaction.block_id = ?", f.BlockId)
	}

	return q, nil
}

// TODO: replace string in StartBlock, EndBlock to int
type ValidatorFilter struct {
	ValidatorPubKey *string
	StartBlock      *string
	EndBlock        *string
}

func(f ValidatorFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.ValidatorPubKey != nil {
		q = q.Join("LEFT JOIN transaction_validator").
			JoinOn("transaction_validator.transaction_id = transaction.id").
			Join("LEFT JOIN validators").
			JoinOn("validators.id = transaction_validator.validator_id").
			Where("validators.public_key = ?", f.ValidatorPubKey)
	}

	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q = q.Apply(blocksRange.Filter)

	return q, nil
}
