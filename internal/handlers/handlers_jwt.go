// Package handlers : file contains operation with requests
package handlers

import (
	"awesomeProject/internal/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Registration godoc
// @Summary Registration
// @Tags    auth
// @Param   person body model.Person true "create user"
// @Produce json
// @Success 200 string
// @Failure 500 string
// @Router  /sign-up [post]
func (h *Handler) Registration(c echo.Context) error {
	person := model.Person{}

	err := json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		log.Errorf("failed parse json, %e", err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%e", err))
	}
	newID, err := h.s.Registration(c.Request().Context(), &person)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%e", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("You register with "+`{"ID":%v}`, newID))
}

// Authentication godoc
// @Summary Authentication
// @Tags    auth
// @Param   id    path string       true "Account ID"
// @Param   login body model.Person true "user password & id"
// @Produce json
// @Accept  json
// @Success 200 string
// @Failure 500 string
// @Router  /login/{id} [post]
func (h *Handler) Authentication(c echo.Context) error {
	auth := model.Authentication{}
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "id cant be empty")
	}
	err = json.NewDecoder(c.Request().Body).Decode(&auth)
	if err != nil {
		log.Errorf("failed parse json, %e", err)
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error with authentication: %e", err))
	}
	accessToken, refreshToken, err := h.s.Authentication(c.Request().Context(), id, auth.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error with authentication: %e", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("You_entry_with "+`{"refreshToken":%v,"accessToken" : %v}`, refreshToken, accessToken))
}

// RefreshToken refresh tokens
func (h *Handler) RefreshToken(c echo.Context) error {
	refreshToken := model.RefreshTokens{}
	err := json.NewDecoder(c.Request().Body).Decode(&refreshToken)
	if err != nil {
		log.Errorf("failed parse json, %e", err)
		return err
	}
	newAccessTokenString, newRefreshTokenString, err := h.s.RefreshToken(c.Request().Context(), refreshToken.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error while creating tokens, %e", err))
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf("Tokens_refresh: "+`{
			"accessToken" : %v,
			"refreshToken" : %v}`, newAccessTokenString, newRefreshTokenString),
		),
	)
}

// Logout godoc
// @Summary  Logout
// @Tags     auth
// @Param    id path string true "Account ID"
// @Accept   string
// @Security ApiKeyAuth
// @Router   /logout/{id} [post]
func (h *Handler) Logout(c echo.Context) error {
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "id cant be empty")
	}
	err = h.s.DeleteFromCache(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed delete user from cache, %e", err))
	}
	err = h.s.UpdateUserAuth(c.Request().Context(), id, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.String(http.StatusOK, "logout")
}
