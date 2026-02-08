package handlers

import (
	"net/http"
	"strings"

	"deepspace/internal/integrations/newapi"
	modelservice "deepspace/internal/service/model"

	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	svc    *modelservice.Service
	newapi *newapi.Client
}

func NewModelHandler(svc *modelservice.Service, newapiClient *newapi.Client) *ModelHandler {
	return &ModelHandler{svc: svc, newapi: newapiClient}
}

type modelSyncResponse struct {
	Items []newapi.UpstreamModel `json:"items"`
}

func (h *ModelHandler) Sync(c *gin.Context) {
	if h == nil || h.newapi == nil {
		respondInternal(c, "上游接口未配置")
		return
	}

	items, err := h.newapi.ListModels(c.Request.Context())
	if err != nil {
		respondInternal(c, "同步上游模型失败")
		return
	}

	c.JSON(http.StatusOK, modelSyncResponse{Items: items})
}

type modelCreateRequest struct {
	Name         string         `json:"name"`
	Provider     string         `json:"provider"`
	PriceInput   float64        `json:"price_input"`
	PriceOutput  float64        `json:"price_output"`
	Currency     string         `json:"currency"`
	Capabilities []string       `json:"capabilities"`
	Status       string         `json:"status"`
	Metadata     map[string]any `json:"metadata"`
}

func (h *ModelHandler) Create(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}

	var req modelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), modelservice.CreateInput{
		Name:         strings.TrimSpace(req.Name),
		Provider:     strings.TrimSpace(req.Provider),
		PriceInput:   req.PriceInput,
		PriceOutput:  req.PriceOutput,
		Currency:     strings.TrimSpace(req.Currency),
		Capabilities: req.Capabilities,
		Status:       strings.TrimSpace(req.Status),
		Metadata:     req.Metadata,
	})
	if err != nil {
		handleModelError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}

type modelUpdateRequest struct {
	Provider     *string         `json:"provider"`
	PriceInput   *float64        `json:"price_input"`
	PriceOutput  *float64        `json:"price_output"`
	Currency     *string         `json:"currency"`
	Capabilities *[]string       `json:"capabilities"`
	Status       *string         `json:"status"`
	ProviderIcon *string         `json:"provider_icon"`
	Metadata     *map[string]any `json:"metadata"`
}

type modelConfirmItem struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type modelConfirmRequest struct {
	Items []modelConfirmItem `json:"items"`
}

type modelPricingItem struct {
	ID           string    `json:"id"`
	PriceInput   *float64  `json:"price_input"`
	PriceOutput  *float64  `json:"price_output"`
	Currency     *string   `json:"currency"`
	Status       *string   `json:"status"`
	Capabilities *[]string `json:"capabilities"`
}

type modelPricingRequest struct {
	Items []modelPricingItem `json:"items"`
}

func (h *ModelHandler) Update(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型ID不正确"})
		return
	}

	var req modelUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	item, err := h.svc.Update(c.Request.Context(), id, modelservice.UpdateInput{
		Provider:     req.Provider,
		PriceInput:   req.PriceInput,
		PriceOutput:  req.PriceOutput,
		Currency:     req.Currency,
		Capabilities: req.Capabilities,
		Status:       req.Status,
		ProviderIcon: req.ProviderIcon,
		Metadata:     req.Metadata,
	})
	if err != nil {
		handleModelError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ModelHandler) List(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}
	provider := strings.TrimSpace(c.Query("provider"))
	var (
		items []modelservice.ModelItem
		err   error
	)
	if provider == "" {
		items, err = h.svc.ListActive(c.Request.Context())
	} else {
		items, err = h.svc.ListActiveByProvider(c.Request.Context(), provider)
	}
	if err != nil {
		respondInternal(c, "获取模型失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ModelHandler) ListProviders(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}

	items, err := h.svc.ListProviders(c.Request.Context(), true)
	if err != nil {
		respondInternal(c, "获取模型提供商失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ModelHandler) ConfirmBatch(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}

	var req modelConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型数据不能为空"})
		return
	}

	inputs := make([]modelservice.ConfirmInput, 0, len(req.Items))
	for _, item := range req.Items {
		inputs = append(inputs, modelservice.ConfirmInput{
			Name:     item.Name,
			Provider: item.Provider,
		})
	}

	result, err := h.svc.ConfirmBatch(c.Request.Context(), inputs)
	if err != nil {
		handleModelError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ModelHandler) ListAll(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}
	provider := strings.TrimSpace(c.Query("provider"))
	var (
		items []modelservice.ModelItem
		err   error
	)
	if provider == "" {
		items, err = h.svc.ListAll(c.Request.Context())
	} else {
		items, err = h.svc.ListAllByProvider(c.Request.Context(), provider)
	}
	if err != nil {
		respondInternal(c, "获取模型失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ModelHandler) BatchPricing(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}

	var req modelPricingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "定价数据不能为空"})
		return
	}

	inputs := make([]modelservice.BatchPricingItem, 0, len(req.Items))
	for _, item := range req.Items {
		inputs = append(inputs, modelservice.BatchPricingItem{
			ID:           strings.TrimSpace(item.ID),
			PriceInput:   item.PriceInput,
			PriceOutput:  item.PriceOutput,
			Currency:     item.Currency,
			Status:       item.Status,
			Capabilities: item.Capabilities,
		})
	}

	result, err := h.svc.BatchPricing(c.Request.Context(), inputs)
	if err != nil {
		handleModelError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ModelHandler) ListAllProviders(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "模型服务未配置")
		return
	}

	items, err := h.svc.ListProviders(c.Request.Context(), false)
	if err != nil {
		respondInternal(c, "获取模型提供商失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func handleModelError(c *gin.Context, err error) {
	switch err {
	case modelservice.ErrInvalidName:
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型名称不正确"})
	case modelservice.ErrInvalidProvider:
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型提供者不正确"})
	case modelservice.ErrInvalidStatus:
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型状态不正确"})
	case modelservice.ErrInvalidCapability, modelservice.ErrInvalidCapabilities:
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型能力不正确"})
	case modelservice.ErrInvalidPrice:
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型价格不正确"})
	case modelservice.ErrInvalidCurrency:
		c.JSON(http.StatusBadRequest, gin.H{"error": "货币类型不正确"})
	case modelservice.ErrDuplicateName:
		c.JSON(http.StatusConflict, gin.H{"error": "模型名称已存在"})
	case modelservice.ErrModelNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "模型不存在"})
	default:
		respondInternal(c, "模型操作失败")
	}
}
