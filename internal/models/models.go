package models

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	Result string `json:"result"`
}
