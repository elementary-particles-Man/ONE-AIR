CREATE TABLE IF NOT EXISTS officers (
    phone TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    pin_hash TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp TEXT NOT NULL,
    officer_phone TEXT NOT NULL,
    mrz TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    action TEXT NOT NULL,
    note TEXT,
    device_id TEXT NOT NULL,
    ip TEXT NOT NULL,
    user_agent TEXT,
    FOREIGN KEY (officer_phone) REFERENCES officers(phone)
);

CREATE TABLE IF NOT EXISTS secondary (
    entry_id INTEGER PRIMARY KEY,
    detail TEXT,
    photo BLOB,
    FOREIGN KEY (entry_id) REFERENCES entries(id)
);
