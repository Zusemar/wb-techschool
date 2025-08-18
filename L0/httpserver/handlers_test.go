package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wb-techschool/L0/internal/domain"

	"github.com/go-chi/chi/v5"
)

type fakeRepo struct{ m map[string]domain.Order }

func (f *fakeRepo) CreateOrder(ctx context.Context, o *domain.Order) error {
	f.m[o.Order_uid] = *o
	return nil
}
func (f *fakeRepo) GetOrderById(ctx context.Context, id string) (*domain.Order, error) {
	o := f.m[id]
	return &o, nil
}
func (f *fakeRepo) UpdateOrder(ctx context.Context, o *domain.Order) error {
	f.m[o.Order_uid] = *o
	return nil
}
func (f *fakeRepo) DeleteOrder(ctx context.Context, id string) error { delete(f.m, id); return nil }
func (f *fakeRepo) ListAllOrderUIDs(ctx context.Context) ([]string, error) {
	var ids []string
	for k := range f.m {
		ids = append(ids, k)
	}
	return ids, nil
}

func TestHandlers_GetAndCreate(t *testing.T) {
	repo := &fakeRepo{m: map[string]domain.Order{}}
	srv := &Server{Repo: repo}
	r := chi.NewRouter()
	r.Get("/order/{order_uid}", srv.GetOrder)
	r.Post("/order", srv.CreateOrder)

	// create
	body := domain.Order{Order_uid: "u1"}
	bb, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status: %d", w.Code)
	}

	// get
	req2 := httptest.NewRequest(http.MethodGet, "/order/u1", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("get status: %d", w2.Code)
	}
}

