package authmiddleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/meshyampratap01/letStayInn/internal/auth"
	"github.com/meshyampratap01/letStayInn/internal/constants"
	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
)

func AuthMiddleware(fn func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		authHeader:=req.Headers["Authorization"]
		if authHeader == ""{
			return events.APIGatewayProxyResponse{
				StatusCode: constants.StatusUnauthorized,
				Body: "Unauthorized",
			},nil
		}

		tokenString := authHeader[len("Bearer "):]

		claims,err := auth.ValidateJWT(tokenString)
		if err!= nil{
			return events.APIGatewayProxyResponse{
				StatusCode: constants.StatusUnauthorized,
				Body: "Unauthorized",
			},nil
		}

		authContext := context.WithValue(ctx,contextkeys.UserIDKey, claims.UserID)
		authContext = context.WithValue(authContext,contextkeys.UserRoleKey,claims.Role)
		authContext = context.WithValue(authContext,contextkeys.UserNameKey,claims.Username)

		return fn(authContext,req)
	}
}