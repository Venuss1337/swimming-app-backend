package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"testProject/internal/models"
)

func (h *Handler) NewWorkout(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var workout map[string]interface{}

	if err := c.Bind(&workout); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.DB.SaveWorkout(user.Id, workout)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "well done maciek")
}
func (h *Handler) ChangeWorkout(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var workout map[string]interface{}
	if err := c.Bind(&workout); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.DB.UpdateWorkout(user.Id, workout);
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "well done maciek")
}
func (h *Handler) DeleteWorkout(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var workout map[string]interface{}
	if err := c.Bind(&workout); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if workout["id"] == nil || workout["id"] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "workout id is required")
	}
	err := h.DB.DeleteWorkout(user.Id, workout["id"].(string));
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "well done maciek")
}
func (h *Handler) GetWorkout(c echo.Context) error {
	user := c.Get("user").(*models.User)

	workouts, err := h.DB.GetAllWorkouts(user.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, workouts)
}