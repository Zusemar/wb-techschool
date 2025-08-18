//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"wb-techschool/L0/internal/domain"
	pg "wb-techschool/L0/internal/repo/postgres"

	_ "github.com/lib/pq"
)

func waitPostgres(dsn string, attempts int, delay time.Duration) error {
	for i := 0; i < attempts; i++ {
		db, err := sql.Open("postgres", dsn)
		if err == nil && db.Ping() == nil {
			db.Close()
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("postgres not ready")
}

func TestIntegration_PostgresRepo(t *testing.T) {
	// Ensure migrations are applied via docker compose profile
	cmd := exec.Command("docker", "compose", "-f", "L0/docker-compose.yaml", "--profile", "migrate", "up", "migrate", "--exit-code-from", "migrate")
	_ = cmd.Run() // best-effort

	cfg := pg.Config{Host: "localhost", Port: 5432, User: "postgres", Password: "postgres", DBName: "postgres", SSLMode: "disable"}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	if err := waitPostgres(dsn, 10, time.Second); err != nil {
		t.Skip("db not ready")
	}

	db, err := pg.NewDB(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := pg.NewOrderRepo(db)
	ctx := context.Background()
	o := &domain.Order{Order_uid: "it2", Track_number: "TN", Entry: "WBIL", Delivery: domain.Delivery{Name: "n", Phone: "p", Zip: "z", City: "c", Address: "a", Region: "r", Email: "e"}, Payment: domain.Payment{Transaction: "t", Currency: "USD", Provider: "wb", Amount: 1, Payment_dt: 2, Bank: "b", Delivery_cost: 3, Goods_total: 4, Custom_fee: 5}, Items: []domain.Item{{Chrt_id: 1, Track_number: "TN", Price: 2, Rid: "rid", Name: "nm", Sale: 3, Size: "s", Total_price: 4, Nm_id: 5, Brand: "br", Status: 6}}, Locale: "en", Customer_id: "c1", Delivery_service: "svc", Shardkey: "1", Sm_id: 1, Date_created: time.Now().UTC(), Oof_shard: "1"}

	if err := repo.CreateOrder(ctx, o); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetOrderById(ctx, o.Order_uid)
	if err != nil {
		t.Fatal(err)
	}
	if got.Order_uid != o.Order_uid {
		t.Fatalf("got %+v", got)
	}
}
