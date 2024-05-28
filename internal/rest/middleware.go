package rest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kavinddd/mangtoon_be/internal/role"
	"github.com/kavinddd/mangtoon_be/internal/token"
	"github.com/kavinddd/mangtoon_be/pkg/util"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware allows all roles to access the resource
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return tokenMiddleware(tokenMaker, []string{})
}

// adminMiddleware allows only admin role to access the resource
func adminMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return tokenMiddleware(tokenMaker, []string{role.Admin})
}

// readerMiddleware allows only reader role to access the resource
func readerMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return tokenMiddleware(tokenMaker, []string{role.Reader})
}

// writerMiddleware allows only reader writer to access the resource
func writerMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return tokenMiddleware(tokenMaker, []string{role.Writer})
}

func tokenMiddleware(tokenMaker token.Maker, accessibleRoles []string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		// roles
		if !hasPermission(payload.Roles, accessibleRoles) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(fmt.Errorf("permission denied")))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}

}

func hasPermission(userRoles []string, accessibleRoles []string) bool {

	if len(accessibleRoles) == 0 {
		return true
	}

	for _, accessibleRole := range accessibleRoles {
		for _, userRole := range userRoles {
			if accessibleRole == userRole {
				return true
			}
		}
	}

	return false
}

func getAuthorizationPayload(ctx *gin.Context) (*token.Payload, error) {
	authPayload, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		return &token.Payload{}, fmt.Errorf("user is not authorized")
	}

	payload, _ := authPayload.(*token.Payload)

	return payload, nil
}

func getCurrentUsername(ctx *gin.Context) (string, error) {
	payload, err := getAuthorizationPayload(ctx)
	if err != nil {
		return "", err
	}
	return payload.Username, nil
}

func getRoles(ctx *gin.Context) ([]string, error) {
	payload, err := getAuthorizationPayload(ctx)
	if err != nil {
		return nil, err
	}
	return payload.Roles, nil
}
