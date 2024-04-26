package api

type Error struct {
	Message string `json:"error"`
}

func errorResponse(err error) Error {
	return Error{
		Message: err.Error(),
	}
}
