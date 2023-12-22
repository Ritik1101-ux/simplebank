package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ritik1101-ux/simplebank/token"
	"github.com/gin-gonic/gin"
)


const (
	authorizationHeaderKey ="authorization"
	authorizationBearer ="bearer"
	authorizationloadKey="authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc{

	return func(ctx *gin.Context) {
		authorizationHeader:=ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader)==0{
			err:=errors.New("Authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return
		}

		fields:=strings.Fields(authorizationHeader)

		if len(fields)<2{
			err:=errors.New("Invalid Authorization header Format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return 
		}

		authorizationType:=strings.ToLower(fields[0])

		if authorizationType!=authorizationBearer{
			err:=errors.New("Unsupported Authorization Type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return 
		}

		accessToken:=fields[1]

		payload ,err:=tokenMaker.VerifyToken(accessToken)

		if err!=nil{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,errorResponse(err))
			return 
		}


		ctx.Set(authorizationloadKey,payload)

		ctx.Next()
		
	}
}