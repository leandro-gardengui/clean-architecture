package web

import (
	"net/http"

	"clean-architecture/internal/usecase"

	"github.com/gin-gonic/gin"
)

type WebOrderHandler struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewWebOrderHandler(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
) *WebOrderHandler {
	return &WebOrderHandler{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (h *WebOrderHandler) CreateOrder(c *gin.Context) {
	var input usecase.CreateOrderInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.CreateOrderUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (h *WebOrderHandler) ListOrders(c *gin.Context) {
	output, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
