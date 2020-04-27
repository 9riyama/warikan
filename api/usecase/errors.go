package usecase

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
