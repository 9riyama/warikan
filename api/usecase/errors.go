package usecase

type BadRequestError struct{}

func (err BadRequestError) Error() string { return "Bad Request" }

type NotFoundError struct{}

func (err NotFoundError) Error() string { return "Not Found" }

type InvalidParamError struct{}

func (err InvalidParamError) Error() string { return "Invalid Parameter" }

type InternalServerError struct{}

func (err InternalServerError) Error() string { return "Internal Server Error" }

type UnauthorizedError struct{}

func (_ UnauthorizedError) Error() string { return "Unauthorized" }

type ConflictError struct{}

func (ConflictError) Error() string { return "Conflict" }

type ServiceUnavailableError struct{}

func (err ServiceUnavailableError) Error() string { return "Service Unavailable" }

type TokenExpiredError struct{}

func (err TokenExpiredError) Error() string { return "Token Expired" }

type InvalidPasscodeError struct{}

func (err InvalidPasscodeError) Error() string { return "Invalid Passcode" }

type CSRFTokenError struct{}

func (_ CSRFTokenError) Error() string { return "CSRF Token Error" }

type SessionExpiredError struct{}

func (_ SessionExpiredError) Error() string { return "session expired" }

type MailAccountLimitError struct{}

func (_ MailAccountLimitError) Error() string { return "Mail accounts limit error" }
