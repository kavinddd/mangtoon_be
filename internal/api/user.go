package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kavinddd/mangtoon_be/internal/db"
	"github.com/kavinddd/mangtoon_be/internal/responseUtils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createUserResponse struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseUtils.ErrorResponse(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseUtils.ErrorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseUtils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createUserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})

}

//func (server *Server) ListActiveUsers(ctx *gin.Context) {
//	args := db.ListActiveUsersParams{
//		Limit:  10,
//		Offset: 0,
//	}
//	users, err := server.store.ListActiveUsers(ctx, args)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, responseUtils.ErrorResponse(err))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, users)
//}

type listUsersRequest struct {
	PageNo   int32 `form:"page_no" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseUtils.ErrorResponse(err))
		return
	}

	args := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageNo - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseUtils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type findUserByIdRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}
type findUserByIdResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) FindUserById(ctx *gin.Context) {
	var req findUserByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseUtils.ErrorResponse(err))
		return
	}

	id, _ := uuid.Parse(req.Id)

	user, err := server.store.GetUser(ctx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responseUtils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, responseUtils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}
