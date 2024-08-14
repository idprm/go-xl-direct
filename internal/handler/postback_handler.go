package handler

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/providers/postback"
)

type PostbackHandler struct {
	logger *logger.Logger
	req    *entity.PostbackParamsRequest
}

func NewPostbackHandler(
	logger *logger.Logger,
	req *entity.PostbackParamsRequest,
) *PostbackHandler {
	return &PostbackHandler{
		logger: logger,
		req:    req,
	}
}

func (h *PostbackHandler) Postback() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service)
	p.Send()
}

func (h *PostbackHandler) Billable() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service)
	p.Billable()
}
