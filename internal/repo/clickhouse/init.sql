CREATE TABLE IF NOT EXISTS events (
    id           String,
    source       LowCardinality(String),
    type         LowCardinality(String),
    external_id  String,
    title        String,
    payload      String,
    collected_at DateTime,
    occurred_at  DateTime
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(occurred_at)
ORDER BY (source, type, occurred_at)
SETTINGS index_granularity = 8192;
