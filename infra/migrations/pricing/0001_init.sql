CREATE TABLE IF NOT EXISTS price_rules (
  id TEXT PRIMARY KEY,
  vendor_id TEXT NOT NULL,
  discount_bps BIGINT NOT NULL
);
