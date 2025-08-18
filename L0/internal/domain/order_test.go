package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOrderJSONRoundtrip(t *testing.T) {
	src := `{
      "order_uid":"uid1",
      "track_number":"TN",
      "entry":"WBIL",
      "delivery":{"name":"n","phone":"p","zip":"z","city":"c","address":"a","region":"r","email":"e"},
      "payment":{"transaction":"t","request_id":"","currency":"USD","provider":"wb","amount":1,"payment_dt":2,"bank":"b","delivery_cost":3,"goods_total":4,"custom_fee":5},
      "items":[{"chrt_id":1,"track_number":"TN","price":2,"rid":"rid","name":"nm","sale":3,"size":"s","total_price":4,"nm_id":5,"brand":"br","status":6}],
      "locale":"en","internal_signature":"","customer_id":"c1","delivery_service":"svc","shardkey":"1","sm_id":9,
      "date_created":"2021-11-26T06:22:19Z","oof_shard":"1"
    }`
	var o Order
	if err := json.Unmarshal([]byte(src), &o); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if o.Order_uid != "uid1" || o.Track_number != "TN" || o.Entry != "WBIL" {
		t.Fatalf("unexpected fields: %+v", o)
	}
	// date parsed
	if o.Date_created.IsZero() {
		t.Fatalf("date_created not parsed")
	}
	// marshal back
	out, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// ensure date format compatible (RFC3339)
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal back: %v", err)
	}
	if _, err := time.Parse(time.RFC3339, m["date_created"].(string)); err != nil {
		t.Fatalf("date_created format: %v", err)
	}
}

