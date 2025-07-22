CREATE TABLE IF NOT EXISTS metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT,
    pod TEXT,
    labels TEXT,
    value REAL NOT NULL,
    timestamp INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_metrics_name ON metrics (name);
CREATE INDEX IF NOT EXISTS idx_metrics_pod ON metrics (pod);
CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics (timestamp);
