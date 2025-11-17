package constants

const (
	StatusOK       = 200
	StatusCreated  = 201
	StatusAccepted = 202

	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusUnprocessableEntity = 422

	StatusInternalServerError = 500
	StatusServiceUnavailable  = 503
)

const (
	ErrInvalidRequestBody    = "invalid request body"
	ErrInvalidCredentials    = "invalid email or password"
	ErrEmailAlreadyExists    = "email already exists"
	ErrUserNotFound          = "user not found"
	ErrFailedToGenerateToken = "unable to generate token"
	ErrInternalServerError   = "internal server error"
	ErrFailedToHashPassword  = "failed to hash password"
	ErrFailedToSaveUser      = "failed to save user"
	ErrFailedToFetchUser     = "failed to fetch user"
	ErrNoAuthorizationToken  = "authorization token is not available"
)
