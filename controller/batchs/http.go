package batchs

import (
	"marketplace-backend/business/batchs"
)

type Controller struct {
	batchUC batchs.UseCase
}

func NewBatchController(batchUC batchs.UseCase) *Controller {
	return &Controller{
		batchUC: batchUC,
	}
}

/*
Create
*/

// func (bc *Controller)

/*
Read
*/

/*
Update
*/

/*
Delete
*/
