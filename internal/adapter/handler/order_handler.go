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
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type OrderHandlerInterface interface {
	CreateOrder(c *gin.Context)
	GetOrderByID(c *gin.Context)
}
type OrderHandler struct {
	orderService service.OrderServiceInterface
	validator    *v.Validator
}

// GetOrder implements OrderHandlerInterface.
func (o *OrderHandler) GetOrderByID(c *gin.Context) {

	var (
		ctx = c.Request.Context()
		res = response.GetOrderByIDResponse{}
	)

	orderID, err := uuid.Parse(c.Param("orderID"))
	if err != nil {
		log.Error().Err(err).Msg("[OrderHandler-4] GetOrderByID")
		c.JSON(http.StatusBadRequest, response.ResponseError(http.StatusBadRequest, "order ID must be a valid UUID"))
		return
	}

	order, err := o.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Error().Err(err).Msg("[OrderHandler-5] GetOrderByID")
		if errors.Is(err, errs.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, response.ResponseError(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	res.OrderID = order.ID
	res.ProductID = order.ProductID
	res.BuyerID = order.BuyerID
	res.Quantity = order.Quantity
	res.TotalCents = order.TotalCents
	res.Status = order.Status

	if order.Product != nil {
		res.Product = &response.ProductDetails{
			ID:         order.Product.ID,
			Name:       order.Product.Name,
			PriceCents: order.Product.PriceCents,
		}
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(http.StatusOK, "success", res))

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
