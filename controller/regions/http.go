package regions

import (
	"crop_connect/business/regions"
)

type Controller struct {
	regionUC regions.UseCase
}

func NewController(regionUC regions.UseCase) *Controller {
	return &Controller{
		regionUC: regionUC,
	}
}
