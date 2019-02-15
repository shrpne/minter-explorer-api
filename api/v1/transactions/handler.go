package transactions

import (
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetTransactionsRequest struct {
	Addresses  []string `form:"addresses[]" binding:"omitempty,minterAddress"`
	Page       string   `form:"page"        binding:"omitempty,numeric"`
	StartBlock *string  `form:"startblock"  binding:"omitempty,numeric"`
	EndBlock   *string  `form:"endblock"    binding:"omitempty,numeric"`
}

type GetTransactionRequest struct {
	Hash string `uri:"hash" binding:"minterTxHash"`
}

// Get list of transactions
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetTransactionsRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// remove Minter wallet prefix from each address
	var minterAddresses []interface{}
	for _, addr := range request.Addresses {
		minterAddresses = append(minterAddresses, helpers.RemoveMinterPrefix(addr))
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxByFilter(transaction.SelectFilter{
		Addresses: minterAddresses,
	}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get transaction detail by hash
func GetTransaction(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetTransactionRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	tx := explorer.TransactionRepository.GetTxByHash(helpers.RemoveMinterPrefix(request.Hash))
	if tx == nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(transaction.Resource).Transform(*tx),
	})
}
