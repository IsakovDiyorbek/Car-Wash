package handler

import (
	"strconv"
	"strings"

	"github.com/exam-5/Car-Wash-Api-Gateway/genproto/carwash"
	"github.com/exam-5/Car-Wash-Api-Gateway/genproto/user"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new provider
// @Description Create a new provider
// @Tags Providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider body carwash.CreateProviderRequest true "provider"
// @Success 200 {object} carwash.CreateProviderResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /providers [post]
func (h *Handler) CreateProvider(c *gin.Context) {
	req := carwash.CreateProviderRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check USER_ID
	// res, err := h.Client.User.GetProfile(c, &user.GetProfileRequest{Id: req.UserId})
	// if err != nil{
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }

	// if res == nil{
	// 	c.JSON(400, gin.H{"error": "User not found"})
	// 	return
	// }

	// check service_id
	count := 1
	for _, serviceId := range req.ServiceId {
		_, err := h.Client.Service.GetService(c, &carwash.GetServiceRequest{Id: serviceId})
		if err != nil {
			c.JSON(400, gin.H{"Service not found id": count})
			return
		}
		count++
	}

	_, err = h.Client.Provider.CreateProvider(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Provider created"})
}

// @Summary Get all providers
// @Description Get all providers
// @Tags Providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param averageRating query number false "averageRating"
// @Param companyName query string false "companyName"
// @Param description query string false "description"
// @Param userId query string false "userId"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} carwash.ListProvidersResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /providers [get]
func (h *Handler) GetProviders(c *gin.Context) {
	req := carwash.ListProvidersRequest{}
	avarageRatingStr := c.Query("avarageRating")
	if avarageRatingStr != "" {
		avarageRatingFloat, err := strconv.Atoi(avarageRatingStr)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.AverageRating = float32(avarageRatingFloat)
	}

	limitStr := c.Query("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Limit = int32(limitInt)
	}

	OffsetStr := c.Query("offset")
	if OffsetStr != "" {
		OffsetInt, err := strconv.Atoi(OffsetStr)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Offset = int32(OffsetInt)
	}

	req.CompanyName = c.Query("companyName")

	req.Description = c.Query("description")

	req.UserId = c.Query("userId")

	res, err := h.Client.Provider.ListProviders(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return

	}
	c.JSON(200, res)
}

// @Summary Update a provider
// @Description Update a provider
// @Tags Providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Param averageRating query number false "averageRating"
// @Param companyName query string false "companyName"
// @Param description query string false "description"
// @Param userId query string false "userId"
// @Param serviceId query string false "serviceId"
// @Param availability query string false "availability"
// @Success 200 {object} carwash.UpdateProviderResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /providers/{id} [put]
func (h *Handler) UpdateProvider(c *gin.Context) {
	req := carwash.UpdateProviderRequest{}

	avarageRatingStr := c.Query("averageRating")
	if avarageRatingStr != "" {
		avarageRatingFloat, err := strconv.ParseFloat(avarageRatingStr, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid averageRating value"})
			return
		}
		req.AverageRating = float32(avarageRatingFloat)
	}

	req.CompanyName = c.Query("companyName")
	req.Description = c.Query("description")
	req.UserId = c.Query("userId")

	req.ServiceId = c.QueryArray("serviceId")

	availabilityStrs := c.QueryArray("availability")
	for _, availabilityStr := range availabilityStrs {
		parts := strings.Split(availabilityStr, ",")
		if len(parts) == 2 {
			req.Availability = append(req.Availability, &carwash.Availability{
				StartTime: parts[0],
				EndTime:   parts[1],
			})
		}
	}

	req.Id = c.Query("id")
	if req.Id == "" {
		c.JSON(400, gin.H{"error": "Missing provider ID"})
		return
	}

	_, err := h.Client.User.GetProfile(c, &user.GetProfileRequest{Id: req.UserId})
	if err != nil {
		c.JSON(400, gin.H{"error": "Provider not found"})
		return
	}
	count := 1
	for _, serviceId := range req.ServiceId {
		_, err := h.Client.Service.GetService(c, &carwash.GetServiceRequest{Id: serviceId})
		if err != nil {
			c.JSON(400, gin.H{"Service not found id": count})
			return
		}
		count++
	}

	_, err = h.Client.Provider.UpdateProvider(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Provider updated successfully"})
}

// @Summary Delete a provider
// @Description Delete a provider
// @Tags Providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.DeleteProviderResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /providers/{id} [delete]
func (h *Handler) DeleteProvider(c *gin.Context) {
	req := carwash.DeleteProviderRequest{}
	req.Id = c.Query("id")


	_, err := h.Client.Provider.DeleteProvider(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

// @Summary Get a provider
// @Description Get a provider
// @Tags Providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.GetProviderResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /providers/{id} [get]
func (h *Handler) GetProvider(c *gin.Context) {
	req := carwash.GetProviderRequest{}
	req.Id = c.Query("id")

	_, err := h.Client.Provider.GetProvider(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "Delete Provider")
}


// @Summary Search providers
// @Description Search providers by company name or description
// @Tags Providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CompanyName query string true "Company Name" 
// @Param Description query string true "Description" 
// @Success 200 {object} carwash.SearchProvidersResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /providers/search [get] // Changed from /providers/search/{companyName}/{description} to /providers/search
func (h *Handler) SearchProviders(c *gin.Context) {
	req := carwash.SearchProvidersRequest{
		CompanyName: c.Query("companyName"),
		Description: c.Query("description"),
	}

	res, err := h.Client.Provider.SearchProviders(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
