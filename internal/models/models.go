package models

type CreateRequest struct {
	Url string `json:"url"`
}

type CreateResponse struct {
	Result string `json:"result"`
}
