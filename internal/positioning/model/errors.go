package model

import (
	"fmt"

	"github.com/cristiantebaldi/goal-detection-api/internal/util"
)

var _ util.BadRequestError = ErrInsufficientAnchors{}
var _ util.BadRequestError = ErrMissingAnchor{}
var _ util.BadRequestError = ErrColinearAnchors{}

type ErrInsufficientAnchors struct{}

func (e ErrInsufficientAnchors) Error() string {
	return "são necessárias pelo menos 3 distâncias de âncoras"
}

func (e ErrInsufficientAnchors) IsBadRequest() {}

type ErrMissingAnchor struct {
	ID string
}

func (e ErrMissingAnchor) Error() string {
	return fmt.Sprintf("distância da âncora %s não informada", e.ID)
}

func (e ErrMissingAnchor) IsBadRequest() {}

type ErrColinearAnchors struct{}

func (e ErrColinearAnchors) Error() string {
	return "âncoras colineares, não é possível calcular a posição"
}

func (e ErrColinearAnchors) IsBadRequest() {}
