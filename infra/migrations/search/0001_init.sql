CREATE TABLE IF NOT EXISTS search_documents (
  id TEXT PRIMARY KEY,
  product_id TEXT NOT NULL,
  vendor_id TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  tier INTEGER NOT NULL DEFAULT 0
);
