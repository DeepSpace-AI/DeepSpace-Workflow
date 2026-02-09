package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"deepspace/internal/repo"
	"deepspace/internal/service/risk"

	"github.com/gin-gonic/gin"
)

type AdminRiskHandler struct {
	svc *risk.Service
}

func NewAdminRiskHandler(svc *risk.Service) *AdminRiskHandler {
	return &AdminRiskHandler{svc: svc}
}

type riskPolicyCreateRequest struct {
	Name      string `json:"name"`
	Scope     string `json:"scope"`
	UserID    *int64 `json:"user_id"`
	ProjectID *int64 `json:"project_id"`
	Status    string `json:"status"`
	Priority  int    `json:"priority"`
}

type riskPolicyUpdateRequest struct {
	Name      *string `json:"name"`
	Scope     *string `json:"scope"`
	UserID    *int64  `json:"user_id"`
	ProjectID *int64  `json:"project_id"`
	Status    *string `json:"status"`
	Priority  *int    `json:"priority"`
}

type rateLimitCreateRequest struct {
	PolicyID      int64  `json:"policy_id"`
	WindowSeconds int    `json:"window_seconds"`
	MaxRequests   int    `json:"max_requests"`
	MaxTokens     int    `json:"max_tokens"`
	Status        string `json:"status"`
}

type rateLimitUpdateRequest struct {
	WindowSeconds *int    `json:"window_seconds"`
	MaxRequests   *int    `json:"max_requests"`
	MaxTokens     *int    `json:"max_tokens"`
	Status        *string `json:"status"`
}

type ipRuleCreateRequest struct {
	PolicyID int64   `json:"policy_id"`
	Type     string  `json:"type"`
	IP       *string `json:"ip"`
	CIDR     *string `json:"cidr"`
	Status   string  `json:"status"`
}

type ipRuleUpdateRequest struct {
	Type   *string `json:"type"`
	IP     *string `json:"ip"`
	CIDR   *string `json:"cidr"`
	Status *string `json:"status"`
}

type budgetCapCreateRequest struct {
	PolicyID int64   `json:"policy_id"`
	Cycle    string  `json:"cycle"`
	MaxCost  float64 `json:"max_cost"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type budgetCapUpdateRequest struct {
	Cycle    *string  `json:"cycle"`
	MaxCost  *float64 `json:"max_cost"`
	Currency *string  `json:"currency"`
	Status   *string  `json:"status"`
}

// ListPolicies godoc
// @Summary 管理员：风控策略列表
// @Description 获取风控策略列表
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param scope query string false "作用域（global/user/project）"
// @Param user_id query int false "用户ID"
// @Param project_id query int false "项目ID"
// @Param status query string false "状态（active/disabled）"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/policies [get]
func (h *AdminRiskHandler) ListPolicies(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}

	userID, err := parseOptionalInt64(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不正确"})
		return
	}
	projectID, err := parseOptionalInt64(c.Query("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目ID不正确"})
		return
	}
	page := parseIntQueryAdmin(c, "page", 1)
	pageSize := parseIntQueryAdmin(c, "page_size", 20)

	items, total, err := h.svc.ListPolicies(c.Request.Context(), repo.RiskPolicyFilter{
		Scope:     strings.TrimSpace(c.Query("scope")),
		UserID:    userID,
		ProjectID: projectID,
		Status:    strings.TrimSpace(c.Query("status")),
		Limit:     pageSize,
		Offset:    (page - 1) * pageSize,
	})
	if err != nil {
		respondInternal(c, "获取风控策略失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// CreatePolicy godoc
// @Summary 管理员：创建风控策略
// @Description 创建风控策略
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body riskPolicyCreateRequest true "策略数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/policies [post]
func (h *AdminRiskHandler) CreatePolicy(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}

	var req riskPolicyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	item, err := h.svc.CreatePolicy(c.Request.Context(), risk.PolicyCreateInput{
		Name:      strings.TrimSpace(req.Name),
		Scope:     strings.TrimSpace(req.Scope),
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
		Status:    strings.TrimSpace(req.Status),
		Priority:  req.Priority,
	})
	if err != nil {
		handleRiskError(c, err, "创建风控策略失败")
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdatePolicy godoc
// @Summary 管理员：更新风控策略
// @Description 更新风控策略
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "策略ID"
// @Param data body riskPolicyUpdateRequest true "策略更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 404 {object} map[string]interface{} "策略不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/policies/{id} [patch]
func (h *AdminRiskHandler) UpdatePolicy(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "策略ID不正确"})
		return
	}

	var req riskPolicyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	item, err := h.svc.UpdatePolicy(c.Request.Context(), id, risk.PolicyUpdateInput{
		Name:      req.Name,
		Scope:     req.Scope,
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
		Status:    req.Status,
		Priority:  req.Priority,
	})
	if err != nil {
		handleRiskError(c, err, "更新风控策略失败")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "策略不存在"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeletePolicy godoc
// @Summary 管理员：删除风控策略
// @Description 删除风控策略
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "策略ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/policies/{id} [delete]
func (h *AdminRiskHandler) DeletePolicy(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "策略ID不正确"})
		return
	}
	_, err = h.svc.DeletePolicy(c.Request.Context(), id)
	if err != nil {
		handleRiskError(c, err, "删除风控策略失败")
		return
	}
	c.Status(http.StatusNoContent)
}

// ListRateLimits godoc
// @Summary 管理员：速率限制列表
// @Description 获取速率限制列表
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param policy_id query int false "策略ID"
// @Param status query string false "状态（active/disabled）"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/rate-limits [get]
func (h *AdminRiskHandler) ListRateLimits(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	policyID, err := parseOptionalInt64(c.Query("policy_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "策略ID不正确"})
		return
	}
	page := parseIntQueryAdmin(c, "page", 1)
	pageSize := parseIntQueryAdmin(c, "page_size", 20)
	items, total, err := h.svc.ListRateLimits(c.Request.Context(), repo.RateLimitFilter{
		PolicyID: policyID,
		Status:   strings.TrimSpace(c.Query("status")),
		Limit:    pageSize,
		Offset:   (page - 1) * pageSize,
	})
	if err != nil {
		respondInternal(c, "获取速率限制失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// CreateRateLimit godoc
// @Summary 管理员：创建速率限制
// @Description 创建速率限制
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body rateLimitCreateRequest true "速率限制数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/rate-limits [post]
func (h *AdminRiskHandler) CreateRateLimit(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	var req rateLimitCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	item, err := h.svc.CreateRateLimit(c.Request.Context(), risk.RateLimitInput{
		PolicyID:      req.PolicyID,
		WindowSeconds: req.WindowSeconds,
		MaxRequests:   req.MaxRequests,
		MaxTokens:     req.MaxTokens,
		Status:        strings.TrimSpace(req.Status),
	})
	if err != nil {
		handleRiskError(c, err, "创建速率限制失败")
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateRateLimit godoc
// @Summary 管理员：更新速率限制
// @Description 更新速率限制
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "速率限制ID"
// @Param data body rateLimitUpdateRequest true "速率限制更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 404 {object} map[string]interface{} "速率限制不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/rate-limits/{id} [patch]
func (h *AdminRiskHandler) UpdateRateLimit(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "速率限制ID不正确"})
		return
	}
	var req rateLimitUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	item, err := h.svc.UpdateRateLimit(c.Request.Context(), id, risk.RateLimitUpdateInput{
		WindowSeconds: req.WindowSeconds,
		MaxRequests:   req.MaxRequests,
		MaxTokens:     req.MaxTokens,
		Status:        req.Status,
	})
	if err != nil {
		handleRiskError(c, err, "更新速率限制失败")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "速率限制不存在"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteRateLimit godoc
// @Summary 管理员：删除速率限制
// @Description 删除速率限制
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "速率限制ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/rate-limits/{id} [delete]
func (h *AdminRiskHandler) DeleteRateLimit(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "速率限制ID不正确"})
		return
	}
	_, err = h.svc.DeleteRateLimit(c.Request.Context(), id)
	if err != nil {
		handleRiskError(c, err, "删除速率限制失败")
		return
	}
	c.Status(http.StatusNoContent)
}

// ListIPRules godoc
// @Summary 管理员：IP 规则列表
// @Description 获取 IP 黑白名单
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param policy_id query int false "策略ID"
// @Param type query string false "类型（allow/deny）"
// @Param status query string false "状态（active/disabled）"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/ip-rules [get]
func (h *AdminRiskHandler) ListIPRules(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	policyID, err := parseOptionalInt64(c.Query("policy_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "策略ID不正确"})
		return
	}
	page := parseIntQueryAdmin(c, "page", 1)
	pageSize := parseIntQueryAdmin(c, "page_size", 20)
	items, total, err := h.svc.ListIPRules(c.Request.Context(), repo.IPRuleFilter{
		PolicyID: policyID,
		Type:     strings.TrimSpace(c.Query("type")),
		Status:   strings.TrimSpace(c.Query("status")),
		Limit:    pageSize,
		Offset:   (page - 1) * pageSize,
	})
	if err != nil {
		respondInternal(c, "获取 IP 规则失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// CreateIPRule godoc
// @Summary 管理员：创建 IP 规则
// @Description 创建 IP 黑白名单规则
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body ipRuleCreateRequest true "IP 规则数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/ip-rules [post]
func (h *AdminRiskHandler) CreateIPRule(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	var req ipRuleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	item, err := h.svc.CreateIPRule(c.Request.Context(), risk.IPRuleInput{
		PolicyID: req.PolicyID,
		Type:     strings.TrimSpace(req.Type),
		IP:       req.IP,
		CIDR:     req.CIDR,
		Status:   strings.TrimSpace(req.Status),
	})
	if err != nil {
		handleRiskError(c, err, "创建 IP 规则失败")
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateIPRule godoc
// @Summary 管理员：更新 IP 规则
// @Description 更新 IP 黑白名单规则
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "IP 规则ID"
// @Param data body ipRuleUpdateRequest true "IP 规则更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 404 {object} map[string]interface{} "IP 规则不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/ip-rules/{id} [patch]
func (h *AdminRiskHandler) UpdateIPRule(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP 规则ID不正确"})
		return
	}
	var req ipRuleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	item, err := h.svc.UpdateIPRule(c.Request.Context(), id, risk.IPRuleUpdateInput{
		Type:   req.Type,
		IP:     req.IP,
		CIDR:   req.CIDR,
		Status: req.Status,
	})
	if err != nil {
		handleRiskError(c, err, "更新 IP 规则失败")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "IP 规则不存在"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteIPRule godoc
// @Summary 管理员：删除 IP 规则
// @Description 删除 IP 黑白名单规则
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "IP 规则ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/ip-rules/{id} [delete]
func (h *AdminRiskHandler) DeleteIPRule(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP 规则ID不正确"})
		return
	}
	_, err = h.svc.DeleteIPRule(c.Request.Context(), id)
	if err != nil {
		handleRiskError(c, err, "删除 IP 规则失败")
		return
	}
	c.Status(http.StatusNoContent)
}

// ListBudgetCaps godoc
// @Summary 管理员：预算上限列表
// @Description 获取预算上限列表
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param policy_id query int false "策略ID"
// @Param status query string false "状态（active/disabled）"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/budget-caps [get]
func (h *AdminRiskHandler) ListBudgetCaps(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	policyID, err := parseOptionalInt64(c.Query("policy_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "策略ID不正确"})
		return
	}
	page := parseIntQueryAdmin(c, "page", 1)
	pageSize := parseIntQueryAdmin(c, "page_size", 20)
	items, total, err := h.svc.ListBudgetCaps(c.Request.Context(), repo.BudgetCapFilter{
		PolicyID: policyID,
		Status:   strings.TrimSpace(c.Query("status")),
		Limit:    pageSize,
		Offset:   (page - 1) * pageSize,
	})
	if err != nil {
		respondInternal(c, "获取预算上限失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// CreateBudgetCap godoc
// @Summary 管理员：创建预算上限
// @Description 创建预算上限
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body budgetCapCreateRequest true "预算上限数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/budget-caps [post]
func (h *AdminRiskHandler) CreateBudgetCap(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	var req budgetCapCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	item, err := h.svc.CreateBudgetCap(c.Request.Context(), risk.BudgetCapInput{
		PolicyID: req.PolicyID,
		Cycle:    strings.TrimSpace(req.Cycle),
		MaxCost:  req.MaxCost,
		Currency: strings.TrimSpace(req.Currency),
		Status:   strings.TrimSpace(req.Status),
	})
	if err != nil {
		handleRiskError(c, err, "创建预算上限失败")
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateBudgetCap godoc
// @Summary 管理员：更新预算上限
// @Description 更新预算上限
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "预算上限ID"
// @Param data body budgetCapUpdateRequest true "预算上限更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 404 {object} map[string]interface{} "预算上限不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/budget-caps/{id} [patch]
func (h *AdminRiskHandler) UpdateBudgetCap(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预算上限ID不正确"})
		return
	}
	var req budgetCapUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}
	item, err := h.svc.UpdateBudgetCap(c.Request.Context(), id, risk.BudgetCapUpdateInput{
		Cycle:    req.Cycle,
		MaxCost:  req.MaxCost,
		Currency: req.Currency,
		Status:   req.Status,
	})
	if err != nil {
		handleRiskError(c, err, "更新预算上限失败")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预算上限不存在"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteBudgetCap godoc
// @Summary 管理员：删除预算上限
// @Description 删除预算上限
// @Tags 管理-风控
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "预算上限ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/risk/budget-caps/{id} [delete]
func (h *AdminRiskHandler) DeleteBudgetCap(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "风控服务未配置")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预算上限ID不正确"})
		return
	}
	_, err = h.svc.DeleteBudgetCap(c.Request.Context(), id)
	if err != nil {
		handleRiskError(c, err, "删除预算上限失败")
		return
	}
	c.Status(http.StatusNoContent)
}

func handleRiskError(c *gin.Context, err error, fallback string) {
	switch err {
	case risk.ErrInvalidScope:
		c.JSON(http.StatusBadRequest, gin.H{"error": "作用域不正确"})
	case risk.ErrInvalidStatus:
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态不正确"})
	case risk.ErrInvalidName:
		c.JSON(http.StatusBadRequest, gin.H{"error": "名称不正确"})
	case risk.ErrInvalidRule:
		c.JSON(http.StatusBadRequest, gin.H{"error": "规则不正确"})
	case risk.ErrInvalidCycle:
		c.JSON(http.StatusBadRequest, gin.H{"error": "周期不正确"})
	case risk.ErrInvalidPolicy:
		c.JSON(http.StatusNotFound, gin.H{"error": "策略不存在"})
	case risk.ErrInvalidIPRule:
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP 规则不正确"})
	default:
		respondInternal(c, fallback)
	}
}
