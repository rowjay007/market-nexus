CREATE TABLE IF NOT EXISTS shipments (
  id TEXT PRIMARY KEY,
  order_id TEXT NOT NULL UNIQUE,
  vendor_id TEXT NOT NULL,
  status TEXT NOT NULL,
  address TEXT
);
