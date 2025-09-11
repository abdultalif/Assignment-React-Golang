package handler

import (
	"backend-service/internal/adapter/handler/request"
	"backend-service/internal/adapter/handler/response"
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"backend-service/internal/core/service"

	v "backend-service/pkg/validator"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type OrderHandlerInterface interface {
	CreateOrder(c *gin.Context)
}
type OrderHandler struct {
	orderService service.OrderServiceInterface
	validator    *v.Validator
}

// CreateOrder implements OrderHandlerInterface.
func (o *OrderHandler) CreateOrder(c *gin.Context) {

	var (
		req = request.CreateOrderRequest{}
		ctx = c.Request.Context()
		res = response.CreateOrderResponse{}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("[OrderHandler-1] Create")
		c.JSON(http.StatusBadRequest, response.ResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := o.validator.Validate(req); err != nil {
		log.Error().Err(err).Msg("[OrderHandler-2] Create")

		if ve, ok := err.(v.ValidationError); ok {
			c.JSON(http.StatusUnprocessableEntity, response.ResponseError(http.StatusUnprocessableEntity, ve.Errors))
			return
		}

		c.JSON(http.StatusUnprocessableEntity, response.ResponseError(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	request := entity.OrderEntity{
		ProductID: req.ProductID,
		BuyerID:   req.BuyerID,
		Quantity:  req.Quantity,
	}

	order, err := o.orderService.CreateOrder(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("[OrderHandler-3] Create")
		if errors.Is(err, errs.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, response.ResponseError(http.StatusNotFound, err.Error()))
			return
		} else if errors.Is(err, errs.ErrOutOfStock) {
			c.JSON(http.StatusUnprocessableEntity, response.ResponseError(http.StatusUnprocessableEntity, err.Error()))
			return
		} else {
			c.JSON(http.StatusInternalServerError, response.ResponseError(http.StatusInternalServerError, err.Error()))
			return
		}
	}

	res.OrderID = order.ID
	res.Status = order.Status

	c.JSON(http.StatusCreated, response.ResponseSuccess(http.StatusCreated, "success", res))

}

func NewOrderHandler(orderService service.OrderServiceInterface, validator *v.Validator) OrderHandlerInterface {
	return &OrderHandler{orderService: orderService, validator: validator}
}
