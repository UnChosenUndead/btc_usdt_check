CREATE TABLE IF NOT EXISTS fiat_date
(
    id SERIAL PRIMARY KEY,
    fiat_date_value BIGINT
);

CREATE TABLE IF NOT EXISTS btc_usd_date
(
    id SERIAL PRIMARY KEY,
    btc_date_value BIGINT
);

CREATE TABLE IF NOT EXISTS btc_usd
(
    id SERIAL PRIMARY KEY,
    btc_date_id int,
    average_price DOUBLE PRECISION,
    FOREIGN KEY (btc_date_id)
    REFERENCES btc_usd_date (id)
);

CREATE TABLE IF NOT EXISTS fiat
(
    id SERIAL PRIMARY KEY,
    fiat_date_id int,
    fiat_char_code TEXT,
    fiat_nominal int,
    fiat_name TEXT,
    fiat_value DOUBLE PRECISION,
        FOREIGN KEY (fiat_date_id)
        REFERENCES fiat_date (id)
);

CREATE TABLE IF NOT EXISTS fiat_btc
(
    id SERIAL PRIMARY KEY,
    fiat_date_id int,
    btc_usd_date_id int,
    fiat_char_code TEXT,
    fiat_btc_sum_value DOUBLE PRECISION,
    FOREIGN KEY (fiat_date_id)
        REFERENCES fiat_date (id),
    FOREIGN KEY (btc_usd_date_id)
        REFERENCES btc_usd_date (id)
);