package http

const (
	internalErrorMessage = "INTERNAL_SERVER_ERROR"
	userIdNotProvided    = "USER_ID_NOT_PROVIDED"
	tokensNotProvided    = "TOKENS_NOT_PROVIDED"
)

type ErrorResponse struct {
	Message string `json:"message"`
}
