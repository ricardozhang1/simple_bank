package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	// 需要对参数进行验证
	Username	string	`json:"username" binding:"required,alphanum"`
	Password	string	`json:"hashed_password" binding:"required,min=6"`
	FullName	string	`json:"full_name" binding:"required"`
	Email		string	`json:"email" binding:"required,email"`
}

type userResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// newResponse 构建请求返回的数据
func newResponse(user db.User) userResponse {
	return userResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}
}

// createUser POST请求创建User
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		// 提交参数验证错误，返回400 请求错误
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		// 明文密码HASH出错，返回500 服务端错误
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 构造请求参数
	arg := db.CreateUserParams{
		Username: 		req.Username,
		HashedPassword:	hashPassword,
		FullName: 		req.FullName,
		Email: 			req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				// 数据库中，插入的表中有唯一值的验证
				// 返回403 禁止提交这个请求
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		// 其他情况下返回500 归结为服务器端的错误
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	rsp := newResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}




