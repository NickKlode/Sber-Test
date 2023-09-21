package handler

import (
	"net/http"
	"sber-test"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Create List
// @Tags list
// @Description create list
// @ID create-list
// @Accept json
// @Produce json
// @Param input body sber.TodoList true "list creation"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /lists [post]
func (h *Handler) createList(c *gin.Context) {
	var input sber.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	if input.Date == "" {
		layout := "2006-01-02"
		t := time.Now()
		dateString := t.Format(layout)
		input.Date = dateString
	}

	id, err := h.service.TodoList.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get All Lists
// @Tags list
// @Description get all lists
// @ID get-all-lists
// @Accept json
// @Produce json
// @Success 200 {object} []sber.TodoList
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /lists [get]
func (h *Handler) getAll(c *gin.Context) {
	list, err := h.service.TodoList.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

// @Summary Delete List
// @Tags list
// @Description delete list
// @ID delete-list
// @Accept json
// @Produce json
// @Param id path int true "list id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id")
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"OK"})
}

// @Summary Update List
// @Tags list
// @Description update list
// @ID update-list
// @Accept json
// @Produce json
// @Param input body sber.UpdateInput true "update list"
// @Param id path int true "list id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id")
		return
	}
	var input sber.UpdateInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.service.TodoList.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"OK"})
}

// @Summary Get List
// @Tags list
// @Description get list
// @ID get-list
// @Accept json
// @Produce json
// @Param input body sber.FindInput true "get list"
// @Success 200 {object} []sber.TodoList
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /lists/find [post]
func (h *Handler) getByDate(c *gin.Context) {
	var input sber.FindInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	list, err := h.service.TodoList.GetByDate(input.Date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}
