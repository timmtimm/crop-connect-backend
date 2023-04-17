package regions

import (
	"marketplace-backend/business/regions"
)

type Controller struct {
	regionUC regions.UseCase
}

func NewRegionController(regionUC regions.UseCase) *Controller {
	return &Controller{
		regionUC: regionUC,
	}
}
