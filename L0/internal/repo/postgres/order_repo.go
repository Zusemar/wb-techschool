package postgres

import (
	"context"
	"database/sql"

	"wb-techschool/L0/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

const insertOrder = `
INSERT INTO orders (
  order_uid, track_number, entry, locale, internal_signature, customer_id,
  delivery_service, shardkey, sm_id, date_created, oof_shard
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

const insertDelivery = `
INSERT INTO delivery (
  order_uid, name, phone, zip, city, address, region, email
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

const insertPayment = `
INSERT INTO payment (
  order_uid, transaction, request_id, currency, provider, amount,
  payment_dt, bank, delivery_cost, goods_total, custom_fee
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

const insertItem = `
INSERT INTO order_items (
  order_uid, chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

const getOrder = `
SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
       delivery_service, shardkey, sm_id, date_created, oof_shard
FROM orders WHERE order_uid = $1`

const getDelivery = `
SELECT name, phone, zip, city, address, region, email
FROM delivery WHERE order_uid = $1`

const getPayment = `
SELECT transaction, request_id, currency, provider, amount, payment_dt, bank,
       delivery_cost, goods_total, custom_fee
FROM payment WHERE order_uid = $1`

const getItems = `
SELECT chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status
FROM order_items WHERE order_uid = $1 ORDER BY item_id`

const updateOrder = `
UPDATE orders SET
  track_number=$2, entry=$3, locale=$4, internal_signature=$5, customer_id=$6,
  delivery_service=$7, shardkey=$8, sm_id=$9, date_created=$10, oof_shard=$11
WHERE order_uid=$1`

const updateDelivery = `
UPDATE delivery SET
  name=$2, phone=$3, zip=$4, city=$5, address=$6, region=$7, email=$8
WHERE order_uid=$1`

const updatePayment = `
UPDATE payment SET
  transaction=$2, request_id=$3, currency=$4, provider=$5, amount=$6,
  payment_dt=$7, bank=$8, delivery_cost=$9, goods_total=$10, custom_fee=$11
WHERE order_uid=$1`

const deleteItemsByOrder = `DELETE FROM order_items WHERE order_uid=$1`
const deleteOrder = `DELETE FROM orders WHERE order_uid=$1`

func (r *OrderRepo) CreateOrder(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// orders
	if _, err = tx.ExecContext(ctx, insertOrder,
		order.Order_uid,
		order.Track_number,
		order.Entry,
		order.Locale,
		order.Internal_signature,
		order.Customer_id,
		order.Delivery_service,
		order.Shardkey,
		order.Sm_id,
		order.Date_created,
		order.Oof_shard,
	); err != nil {
		tx.Rollback()
		return err
	}

	// delivery
	if _, err = tx.ExecContext(ctx, insertDelivery,
		order.Order_uid,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	); err != nil {
		tx.Rollback()
		return err
	}

	// payment
	if _, err = tx.ExecContext(ctx, insertPayment,
		order.Order_uid,
		order.Payment.Transaction,
		order.Payment.Request_id,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.Payment_dt,
		order.Payment.Bank,
		order.Payment.Delivery_cost,
		order.Payment.Goods_total,
		order.Payment.Custom_fee,
	); err != nil {
		tx.Rollback()
		return err
	}

	// items
	for _, it := range order.Items {
		if _, err = tx.ExecContext(ctx, insertItem,
			order.Order_uid,
			it.Chrt_id,
			it.Price,
			it.Rid,
			it.Name,
			it.Sale,
			it.Size,
			it.Total_price,
			it.Nm_id,
			it.Brand,
			it.Status,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepo) GetOrderById(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order

	// order
	row := r.db.QueryRowContext(ctx, getOrder, id)
	if err := row.Scan(
		&o.Order_uid,
		&o.Track_number,
		&o.Entry,
		&o.Locale,
		&o.Internal_signature,
		&o.Customer_id,
		&o.Delivery_service,
		&o.Shardkey,
		&o.Sm_id,
		&o.Date_created,
		&o.Oof_shard,
	); err != nil {
		return nil, err
	}

	// delivery
	drow := r.db.QueryRowContext(ctx, getDelivery, id)
	if err := drow.Scan(
		&o.Delivery.Name,
		&o.Delivery.Phone,
		&o.Delivery.Zip,
		&o.Delivery.City,
		&o.Delivery.Address,
		&o.Delivery.Region,
		&o.Delivery.Email,
	); err != nil {
		return nil, err
	}

	// payment
	prow := r.db.QueryRowContext(ctx, getPayment, id)
	if err := prow.Scan(
		&o.Payment.Transaction,
		&o.Payment.Request_id,
		&o.Payment.Currency,
		&o.Payment.Provider,
		&o.Payment.Amount,
		&o.Payment.Payment_dt,
		&o.Payment.Bank,
		&o.Payment.Delivery_cost,
		&o.Payment.Goods_total,
		&o.Payment.Custom_fee,
	); err != nil {
		return nil, err
	}

	// items
	rows, err := r.db.QueryContext(ctx, getItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	o.Items = make([]domain.Item, 0)
	for rows.Next() {
		var it domain.Item
		if err := rows.Scan(
			&it.Chrt_id,
			&it.Price,
			&it.Rid,
			&it.Name,
			&it.Sale,
			&it.Size,
			&it.Total_price,
			&it.Nm_id,
			&it.Brand,
			&it.Status,
		); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, updateOrder,
		order.Order_uid,
		order.Track_number,
		order.Entry,
		order.Locale,
		order.Internal_signature,
		order.Customer_id,
		order.Delivery_service,
		order.Shardkey,
		order.Sm_id,
		order.Date_created,
		order.Oof_shard,
	); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.ExecContext(ctx, updateDelivery,
		order.Order_uid,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.ExecContext(ctx, updatePayment,
		order.Order_uid,
		order.Payment.Transaction,
		order.Payment.Request_id,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.Payment_dt,
		order.Payment.Bank,
		order.Payment.Delivery_cost,
		order.Payment.Goods_total,
		order.Payment.Custom_fee,
	); err != nil {
		tx.Rollback()
		return err
	}

	// Recreate items set
	if _, err = tx.ExecContext(ctx, deleteItemsByOrder, order.Order_uid); err != nil {
		tx.Rollback()
		return err
	}
	for _, it := range order.Items {
		if _, err = tx.ExecContext(ctx, insertItem,
			order.Order_uid,
			it.Chrt_id,
			it.Price,
			it.Rid,
			it.Name,
			it.Sale,
			it.Size,
			it.Total_price,
			it.Nm_id,
			it.Brand,
			it.Status,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepo) DeleteOrder(ctx context.Context, id string) error {
	if _, err := r.db.ExecContext(ctx, deleteOrder, id); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepo) ListAllOrderUIDs(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT order_uid FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}
