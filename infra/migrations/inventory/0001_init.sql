CREATE TABLE IF NOT EXISTS stock_items (
  id TEXT PRIMARY KEY,
  vendor_id TEXT NOT NULL,
  sku TEXT NOT NULL,
  available INTEGER NOT NULL,
  reserved INTEGER NOT NULL DEFAULT 0,
  version INTEGER NOT NULL DEFAULT 0,
  UNIQUE (vendor_id, sku)
);
