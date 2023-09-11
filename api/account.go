package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/nphattai/go-simple-bank/db/sqlc"
	"github.com/nphattai/go-simple-bank/token"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (sever *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := sever.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code.Name() == "unique_violation" || pgErr.Code.Name() == "foreign_key_violation" {
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (sever *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := sever.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getListAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

func (sever *Server) getListAccounts(ctx *gin.Context) {
	var req getListAccountsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := sever.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
		Owner:  authPayload.Username,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
