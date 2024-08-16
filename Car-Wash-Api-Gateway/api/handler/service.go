package handler

import (
	"fmt"
	"strconv"

	"github.com/exam-5/Car-Wash-Api-Gateway/genproto/carwash"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new service
// @Description Create a new service
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param service body carwash.CreateServiceRequest true "service"
// @Success 200 {object} carwash.CreateServiceResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /services [post]
func (h *Handler) CreateService(c *gin.Context) {
	req := carwash.CreateServiceRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err = h.Client.Service.CreateService(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Service created"})
}

// @Summary Get all services
// @Description Get all services
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param description query string false "description"
// @Param name query string false "name"
// @Param price query number false "price"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param duration query int false "duration"
// @Success 200 {object} carwash.ListServicesResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /services [get]
func (h *Handler) GetServices(c *gin.Context) {
	req := carwash.ListServicesRequest{}
	req.Description = c.Query("description")
	req.Name = c.Query("name")
	priceStr := c.Query("price")
	if priceStr != "" {
		priceFloat, err := strconv.ParseFloat(priceStr, 32)
		req.Price = float32(priceFloat)
		if err != nil {
			fmt.Println("err while parsing float", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	limitStr := c.Query("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			fmt.Println("error while converting limit", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Limit = int32(limitInt)
	}

	OffsetStr := c.Query("offset")
	if OffsetStr != "" {
		OffsetInt, err := strconv.Atoi(OffsetStr)
		if err != nil {
			fmt.Println("error while converting offset", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Offset = int32(OffsetInt)
	}

	durationStr := c.Query("duration")
	if durationStr != "" {
		durationInt, err := strconv.Atoi(durationStr)
		if err != nil {
			fmt.Println("error while converting duration", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Duration = int32(durationInt)
	}

	res, err := h.Client.Service.ListServices(c, &req)
	if err != nil {
		fmt.Println("error while getting services", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// @Summary Update a service
// @Description Update a service
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Param name query string false "name"
// @Param description query string false "description"
// @Param price query number false "price"
// @Param duration query int false "duration"
// @Success 200 {object} carwash.UpdateServiceResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /services/{id} [put]
func (h *Handler) UpdateService(c *gin.Context) {
	req := carwash.UpdateServiceRequest{}
	req.Id = c.Query("id")
	req.Description = c.Query("description")
	req.Name = c.Query("name")
	priceStr := c.Query("price")
	if priceStr != "" {
		priceFloat, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Price = float32(priceFloat)
	}

	durationStrs := c.Query("duration")
	if durationStrs != "" {
		durationInt, err := strconv.Atoi(durationStrs)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Duration = int32(durationInt)

	}

	_, err := h.Client.Service.UpdateService(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Service updated"})
}

// @Summary Delete a service
// @Description Delete a service
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.DeleteServiceResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /services/{id} [delete]
func (h *Handler) DeleteService(c *gin.Context) {
	id := c.Query("id")
	_, err := h.Client.Service.DeleteService(c, &carwash.DeleteServiceRequest{Id: id})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Service deleted"})
}

// @Summary Get a service
// @Description Get a service
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.GetServiceResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /services/{id} [get]
func (h *Handler) GetService(c *gin.Context) {
	id := c.Query("id")
	res, err := h.Client.Service.GetService(c, &carwash.GetServiceRequest{Id: id})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// write swagger

// @Summary Search services
// @Description Search services by name or description
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string true "name"
// @Param description query string true "description"
// @Success 200 {object} carwash.SearchServicesResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /services/search [get]
func (h *Handler) SearchServices(c *gin.Context) {
	req := carwash.SearchServicesRequest{}
	req.Name = c.Query("name")
	req.Description = c.Query("description")
	res, err := h.Client.Service.SearchServices(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
