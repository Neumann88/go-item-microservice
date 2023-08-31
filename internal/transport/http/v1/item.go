package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"item-service/internal/entity"
	"item-service/pkg/http/response"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func (h HttpHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var inItem entity.Item
	err = json.Unmarshal(b, &inItem)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.app.Command.CreateItem.Handle(ctx, inItem)
	if err != nil {
		logrus.Errorf("CreateItem: ERROR(%v), DATA(%v)", err, inItem)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h HttpHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	inItemID := chi.URLParam(r, "id")
	if inItemID == "" {
		response.Error(w, "undefined item id", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(inItemID)
	if err != nil {
		response.Error(w, fmt.Sprintf("invalid item id %s", err.Error()), http.StatusBadRequest)
		return
	}

	item, err := h.app.Query.GetItemByID.Handle(ctx, itemID)
	if err != nil {
		logrus.Errorf("GetItemByID: ERROR(%v), DATA(%v)", err, itemID)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = response.Success(w, item, http.StatusOK)
	if err != nil {
		logrus.Errorf("GetItemByID-Response: ERROR(%v), DATA(%v)", err, itemID)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
