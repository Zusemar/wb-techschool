package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"wb-techschool/L0/internal/domain"

	_ "github.com/lib/pq"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	cfg := Config{
		Host:     getenv("DB_HOST", "localhost"),
		Port:     getenvInt("DB_PORT", 5432),
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASSWORD", "postgres"),
		DBName:   getenv("DB_NAME", "postgres"),
		SSLMode:  getenv("DB_SSLMODE", "disable"),
	}
	db, err := NewDB(cfg)
	if err != nil {
		t.Fatalf("db open: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("db ping: %v", err)
	}
	return db
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func getenvInt(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		var n int
		fmt.Sscanf(v, "%d", &n)
		if n > 0 {
			return n
		}
	}
	return d
}

func TestOrderRepo_CRUD(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	repo := NewOrderRepo(db)
	ctx := context.Background()

	o := &domain.Order{
		Order_uid:        "it-uid-1",
		Track_number:     "TN",
		Entry:            "WBIL",
		Delivery:         domain.Delivery{Name: "n", Phone: "p", Zip: "z", City: "c", Address: "a", Region: "r", Email: "e"},
		Payment:          domain.Payment{Transaction: "t", Currency: "USD", Provider: "wb", Amount: 1, Payment_dt: 2, Bank: "b", Delivery_cost: 3, Goods_total: 4, Custom_fee: 5},
		Items:            []domain.Item{{Chrt_id: 1, Track_number: "TN", Price: 2, Rid: "rid", Name: "nm", Sale: 3, Size: "s", Total_price: 4, Nm_id: 5, Brand: "br", Status: 6}},
		Locale:           "en",
		Customer_id:      "c1",
		Delivery_service: "svc",
		Shardkey:         "1",
		Sm_id:            9,
		Date_created:     time.Now().UTC(),
		Oof_shard:        "1",
	}

	if err := repo.CreateOrder(ctx, o); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.GetOrderById(ctx, o.Order_uid)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Order_uid != o.Order_uid || len(got.Items) != 1 {
		t.Fatalf("mismatch: %+v", got)
	}

	// Update some fields and re-save
	o.Locale = "ru"
	if err := repo.UpdateOrder(ctx, o); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err = repo.GetOrderById(ctx, o.Order_uid)
	if err != nil || got.Locale != "ru" {
		t.Fatalf("update check: %+v %v", got, err)
	}

	// List IDs
	ids, err := repo.ListAllOrderUIDs(ctx)
	if err != nil || len(ids) == 0 {
		t.Fatalf("list ids: %v %v", ids, err)
	}

	// Delete
	if err := repo.DeleteOrder(ctx, o.Order_uid); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
