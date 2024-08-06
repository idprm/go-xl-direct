package handler

type PostbackHandler struct {
}

func NewPostbackHandler() *PostbackHandler {
	return &PostbackHandler{}
}

func (h *PostbackHandler) Postback() {
}

func (h *PostbackHandler) Billable() {
}
