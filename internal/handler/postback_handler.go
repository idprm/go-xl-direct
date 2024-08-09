package handler

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/logger"
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
}

func (h *PostbackHandler) Billable() {
}
