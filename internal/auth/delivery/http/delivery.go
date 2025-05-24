package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t1d333/smartlectures/internal/auth"
	authmodels "github.com/t1d333/smartlectures/internal/auth/models"
	"github.com/t1d333/smartlectures/internal/auth/service"
	"github.com/t1d333/smartlectures/internal/errors"

	// "github.com/t1d333/smartlectures/internal/errors"
	// "github.com/t1d333/smartlectures/internal/models"
	"github.com/t1d333/smartlectures/pkg/logger"
)

type Delivery struct {
	logger  logger.Logger
	mux     *gin.Engine
	service auth.Service
}

func NewDelivery(logger logger.Logger, mux *gin.Engine, service auth.Service) *Delivery {
	return &Delivery{
		logger:  logger,
		mux:     mux,
		service: service,
	}
}

func (d *Delivery) RegisterHandler() {
	d.mux.POST("/api/v1/auth/login", d.Login)
	d.mux.DELETE("/api/v1/auth/logout", d.Logout)
	d.mux.POST("/api/v1/auth/refresh", d.Refresh)
	d.mux.PUT("/api/v1/auth/register", d.Register)
	d.mux.GET("/api/v1/auth/me", d.GetMe)
}

func (d *Delivery) Login(c *gin.Context) {
	data := authmodels.LoginRequest{}

	if err := c.BindJSON(&data); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	}

	token, err := d.service.Login(context.WithValue(c.Request.Context(), "client_ip", c.ClientIP()), data)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.SetCookie(
		"session",
		token,
		int(service.Expire.Seconds()),
		"/",
		"localhost",
		false,
		true,
	)
	c.Status(http.StatusNoContent)
}

func (d *Delivery) Logout(c *gin.Context) {
	token, err := c.Cookie("session")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err = d.service.Logout(c.Request.Context(), token); err != nil {
		_ = c.Error(err)
		return
	}

	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	c.Status(http.StatusNoContent)
}

func (d *Delivery) GetMe(c *gin.Context) {
	session, err := c.Cookie("session")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := d.service.GetMe(c.Request.Context(), session)
	
	if err != nil {
		_ = c.Error(err)
		return
	}
	
	c.JSON(http.StatusOK, user)

}

func (d *Delivery) Register(c *gin.Context) {
	data := authmodels.RegisterRequest{}

	if err := c.BindJSON(&data); err != nil {
		_ = c.Error(errors.BadRequestError)
		return
	}

	user, err := d.service.Register(c.Request.Context(), data)
	if err != nil {
		_ = c.Error(err)
		return
	}

	token, err := d.service.Login(context.WithValue(c.Request.Context(), "client_ip", c.ClientIP()), authmodels.LoginRequest{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		_ =  c.Error(err)
		return
	}

	c.SetCookie(
		"session",
		token,
		int(service.Expire.Seconds()),
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, user)
}

func (d *Delivery) Refresh(c *gin.Context) {
}
