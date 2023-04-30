package pkg

type StatusResponse struct {
	Status int `json:"status"`
}

func NewStatusResponse(status int) StatusResponse {
	return StatusResponse{
		Status: status,
	}
}
