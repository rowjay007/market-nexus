CREATE TABLE IF NOT EXISTS behavior_events (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  product_id TEXT NOT NULL,
  vendor_id TEXT NOT NULL,
  behavior TEXT NOT NULL,
  occurred_at TIMESTAMP NOT NULL
);
