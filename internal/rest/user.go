package rest

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kavinddd/mangtoon_be/internal/db"
	"github.com/kavinddd/mangtoon_be/internal/role"
	"github.com/kavinddd/mangtoon_be/pkg/util"
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

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
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

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	args := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageNo - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type findUserByIdRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}
type userResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) findUserById(ctx *gin.Context) {
	var req findUserByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	id, _ := uuid.Parse(req.Id)

	user, err := server.store.GetUserById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	payload, _ := getAuthorizationPayload(ctx)

	if !role.IsAdmin(payload.Roles) {
		if user.Username != payload.Username {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(sql.ErrNoRows))
			return
		}
	}

	response := newUserResponse(user)

	ctx.JSON(http.StatusOK, response)
	return
}
