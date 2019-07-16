package postgres

import (
	"github.com/jmoiron/sqlx"
)

const schema = `CREATE TABLE IF NOT EXISTS "exchanges"
(
    id                INT PRIMARY KEY,
    name              TEXT UNIQUE NOT NULL,
    trading_fee_maker FLOAT,
    trading_fee_taker FLOAT,
    disable           BOOLEAN     NOT NULL DEFAULT TRUE
        -- only allow to enable exchange if trading_fee_maker and trading_fee_taker are both set
        CONSTRAINT disable_check CHECK (disable OR
                                        ((trading_fee_maker IS NOT NULL) AND (trading_fee_taker IS NOT NULL)))
);

CREATE TABLE IF NOT EXISTS "addresses"
(
    id      SERIAL PRIMARY KEY,
    address TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "assets"
(
    id                            SERIAL PRIMARY KEY,
    symbol                        TEXT      NOT NULL UNIQUE,
    name                          TEXT      NOT NULL,
    address_id                    INT       NULL REFERENCES addresses (id)
        CONSTRAINT address_id_check CHECK ( address_id IS NOT NULL OR NOT transferable),
    decimals                      INT       NOT NULL,
    -- transferable must be set to true if it is possible to withdraw/deposit 
    -- to reserve (ETH or ERC20 tokens). If transferable is true, the address and 
    -- deposit address of related asset_exchange records are required.
    transferable                  BOOLEAN   NOT NULL,
    set_rate                      TEXT      NOT NULL,
    rebalance                     BOOLEAN   NOT NULL,
    is_quote                      BOOLEAN   NOT NULL,

    pwi_ask_a                     FLOAT     NULL,
    pwi_ask_b                     FLOAT     NULL,
    pwi_ask_c                     FLOAT     NULL,
    pwi_ask_min_min_spread        FLOAT     NULL,
    pwi_ask_price_multiply_factor FLOAT     NULL,
    pwi_bid_a                     FLOAT     NULL,
    pwi_bid_b                     FLOAT     NULL,
    pwi_bid_c                     FLOAT     NULL,
    pwi_bid_min_min_spread        FLOAT     NULL,
    pwi_bid_price_multiply_factor FLOAT     NULL,

    rebalance_quadratic_a         FLOAT     NULL,
    rebalance_quadratic_b         FLOAT     NULL,
    rebalance_quadratic_c         FLOAT     NULL,

    target_total                  FLOAT     NULL,
    target_reserve                FLOAT     NULL,
    target_rebalance_threshold    FLOAT     NULL,
    target_transfer_threshold     FLOAT     NULL,

    created                       TIMESTAMP NOT NULL,
    updated                       TIMESTAMP NOT NULL,
    -- if set_rate strategy is defined, pwi columns are required
    CONSTRAINT pwi_check CHECK (
            set_rate = 'not_set'
            OR (pwi_ask_a IS NOT NULL AND
                pwi_ask_b IS NOT NULL AND
                pwi_ask_c IS NOT NULL AND
                pwi_ask_min_min_spread IS NOT NULL AND
                pwi_ask_price_multiply_factor IS NOT NULL AND
                pwi_bid_a IS NOT NULL AND
                pwi_bid_b IS NOT NULL AND
                pwi_bid_c IS NOT NULL AND
                pwi_bid_min_min_spread IS NOT NULL AND
                pwi_bid_price_multiply_factor IS NOT NULL
                )),
    -- if rebalance is true, rebalance quadratic is required
    CONSTRAINT rebalance_quadratic_check CHECK (
            NOT rebalance OR
            (rebalance_quadratic_a IS NOT NULL AND
             rebalance_quadratic_b IS NOT NULL AND
             rebalance_quadratic_c IS NOT NULL)),
    -- if rebalance is true, target configuration is required
    CONSTRAINT target_check CHECK (
            NOT rebalance OR
            (target_total IS NOT NULL AND
             target_reserve IS NOT NULL AND
             target_rebalance_threshold IS NOT NULL AND
             target_transfer_threshold IS NOT NULL))
);

CREATE TABLE IF NOT EXISTS "asset_old_addresses"
(
    id         SERIAL PRIMARY KEY,
    address_id INT NOT NULL REFERENCES addresses (id),
    asset_id   INT NOT NULL REFERENCES assets (id)
    -- TODO add a constraint to ensure that asset_id is not linked to any asset in address field already 
);

CREATE TABLE IF NOT EXISTS "asset_exchanges"
(
    id                 SERIAL PRIMARY KEY,
    exchange_id        INT REFERENCES exchanges (id) NOT NULL,
    asset_id           INT REFERENCES assets (id)    NOT NULL,
    symbol             TEXT                          NOT NULL,
    deposit_address    TEXT                          NULL,
    min_deposit        FLOAT                         NOT NULL,
    withdraw_fee       FLOAT                         NOT NULL,
    price_precision    INT                           NOT NULL,
    amount_precision   INT                           NOT NULL,
    amount_limit_min   FLOAT                         NOT NULL,
    amount_limit_max   FLOAT                         NOT NULL,
    price_limit_min    FLOAT                         NOT NULL,
    price_limit_max    FLOAT                         NOT NULL,
    target_recommended FLOAT                         NOT NULL,
    target_ratio       FLOAT                         NOT NULL
);

CREATE TABLE IF NOT EXISTS trading_pairs
(
    id           SERIAL PRIMARY KEY,
    exchange_id  INT REFERENCES exchanges (id) NOT NULL,
    base_id      INT REFERENCES assets (id)    NOT NULL,
    quote_id     INT REFERENCES assets (id)    NOT NULL,
    min_notional FLOAT                         NOT NULL,
    UNIQUE (exchange_id, base_id, quote_id),
    CONSTRAINT trading_pair_check CHECK ( base_id != quote_id)
);

CREATE OR REPLACE FUNCTION new_asset(_symbol assets.symbol%TYPE,
                                     _name assets.symbol%TYPE,
                                     _address addresses.address%TYPE,
                                     _decimals assets.decimals%TYPE,
                                     _transferable assets.transferable%TYPE,
                                     _set_rate assets.set_rate%TYPE,
                                     _rebalance assets.rebalance%TYPE,
                                     _is_quote assets.is_quote%TYPE,
                                     _pwi_ask_a assets.pwi_ask_a%TYPE,
                                     _pwi_ask_b assets.pwi_ask_b%TYPE,
                                     _pwi_ask_c assets.pwi_ask_c%TYPE,
                                     _pwi_ask_min_min_spread assets.pwi_ask_min_min_spread%TYPE,
                                     _pwi_ask_price_multiply_factor assets.pwi_ask_price_multiply_factor%TYPE,
                                     _pwi_bid_a assets.pwi_bid_a%TYPE,
                                     _pwi_bid_b assets.pwi_bid_b%TYPE,
                                     _pwi_bid_c assets.pwi_bid_c%TYPE,
                                     _pwi_bid_min_min_spread assets.pwi_bid_min_min_spread%TYPE,
                                     _pwi_bid_price_multiply_factor assets.pwi_bid_price_multiply_factor%TYPE,
                                     _rebalance_quadratic_a assets.rebalance_quadratic_a%TYPE,
                                     _rebalance_quadratic_b assets.rebalance_quadratic_b%TYPE,
                                     _rebalance_quadratic_c assets.rebalance_quadratic_c%TYPE,
                                     _target_total assets.target_total%TYPE,
                                     _target_reserve assets.target_reserve%TYPE,
                                     _target_rebalance_threshold assets.target_rebalance_threshold%TYPE,
                                     _target_transfer_threshold assets.target_total%TYPE)
    RETURNS int AS
$$
DECLARE
    _address_id addresses.id%TYPE;
    _id         assets.id%TYPE;
BEGIN
    IF _address IS NOT NULL THEN
        INSERT INTO "addresses" (address) VALUES (_address) RETURNING id INTO _address_id;
    END IF;

    INSERT
    INTO assets(symbol,
                name,
                address_id,
                decimals,
                transferable,
                set_rate,
                rebalance,
                is_quote,
                pwi_ask_a,
                pwi_ask_b,
                pwi_ask_c,
                pwi_ask_min_min_spread,
                pwi_ask_price_multiply_factor,
                pwi_bid_a,
                pwi_bid_b,
                pwi_bid_c,
                pwi_bid_min_min_spread,
                pwi_bid_price_multiply_factor,
                rebalance_quadratic_a,
                rebalance_quadratic_b,
                rebalance_quadratic_c,
                target_total,
                target_reserve,
                target_rebalance_threshold,
                target_transfer_threshold,
                created,
                updated)
    VALUES (_symbol,
            _name,
            _address_id,
            _decimals,
            _transferable,
            _set_rate,
            _rebalance,
            _is_quote,
            _pwi_ask_a,
            _pwi_ask_b,
            _pwi_ask_c,
            _pwi_ask_min_min_spread,
            _pwi_ask_price_multiply_factor,
            _pwi_bid_a,
            _pwi_bid_b,
            _pwi_bid_c,
            _pwi_bid_min_min_spread,
            _pwi_bid_price_multiply_factor,
            _rebalance_quadratic_a,
            _rebalance_quadratic_b,
            _rebalance_quadratic_c,
            _target_total,
            _target_reserve,
            _target_rebalance_threshold,
            _target_transfer_threshold,
            now(),
            now()) RETURNING id INTO _id;

    RETURN _id;
END
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION change_asset_address(_id assets.id%TYPE, _address addresses.address%TYPE) RETURNS VOID AS
$$
DECLARE
    _new_address_id addresses.id%TYPE;
BEGIN
    PERFORM id FROM assets WHERE id = _id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'asset with id % does not exists', _id USING ERRCODE = 'assert_failure';
    END IF;

    INSERT INTO "asset_old_addresses" (address_id, asset_id)
    SELECT addresses.id, assets.id
    FROM assets
             LEFT JOIN addresses ON assets.address_id = addresses.id
    WHERE assets.id = _id;

    INSERT INTO "addresses" (address) VALUES (_address) RETURNING id INTO _new_address_id;

    UPDATE "assets"
    SET address_id = _new_address_id,
        updated    = now()
    WHERE assets.id = _id;
    RETURN;
END
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE FUNCTION new_trading_pair(_exchange_id trading_pairs.exchange_id%TYPE,
                                            _base_id trading_pairs.base_id%TYPE,
                                            _quote_id trading_pairs.quote_id%TYPE,
                                            _min_notional trading_pairs.min_notional%TYPE)
    RETURNS INT AS
$$
DECLARE
    _id                   trading_pairs.id%TYPE;
    _quote_asset_is_quote assets.is_quote%TYPE;
BEGIN
    PERFORM id FROM asset_exchanges WHERE exchange_id = _exchange_id AND asset_id = _base_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'base asset is not configured for exchange base_id=% exchange_id=%',
            _base_id,_exchange_id USING ERRCODE = 'assert_failure';
    END IF;

    PERFORM id FROM asset_exchanges WHERE exchange_id = _exchange_id AND asset_id = _quote_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'quote asset is not configured for exchange quote_id=% exchange_id=%',
            _quote_id,_exchange_id USING ERRCODE = 'assert_failure';
    END IF;

    PERFORM id FROM assets WHERE id = _base_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'base asset is not found base_id=%', _base_id USING ERRCODE = 'assert_failure';
    END IF;

    SELECT is_quote FROM assets WHERE id = _quote_id INTO _quote_asset_is_quote;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'quote asset is not found quote_id=%', _quote_id USING ERRCODE = 'assert_failure';
    END IF;

    IF NOT _quote_asset_is_quote THEN
        RAISE EXCEPTION 'quote asset is not configured as quote id=%', _quote_id USING ERRCODE = 'assert_failure';
    END IF;

    INSERT INTO trading_pairs (exchange_id, base_id, quote_id, min_notional)
    VALUES (_exchange_id, _base_id, _quote_id, _min_notional) RETURNING id INTO _id;
    RETURN _id;
END
$$ LANGUAGE PLPGSQL;
`

type preparedStmts struct {
	getExchanges     *sqlx.Stmt
	getExchange      *sqlx.Stmt
	updateExchange   *sqlx.NamedStmt
	newAsset         *sqlx.NamedStmt
	newAssetExchange *sqlx.NamedStmt
	newTradingPair   *sqlx.Stmt

	getAsset             *sqlx.Stmt
	getAssetExchange     *sqlx.Stmt
	getTradingPair       *sqlx.Stmt
	updateAsset          *sqlx.NamedStmt
	changeAssetAddress   *sqlx.Stmt
	updateDepositAddress *sqlx.Stmt

	getTradingPairSymbols *sqlx.Stmt
	getMinNotional        *sqlx.Stmt
	getTransferableAssets *sqlx.Stmt
}

func newPreparedStmts(db *sqlx.DB) (*preparedStmts, error) {
	const getExchangesQuery = `SELECT * FROM "exchanges";`
	getExchanges, err := db.Preparex(getExchangesQuery)
	if err != nil {
		return nil, err
	}

	const getExchangeQuery = `SELECT * FROM "exchanges" WHERE id = $1`
	getExchange, err := db.Preparex(getExchangeQuery)
	if err != nil {
		return nil, err
	}

	const updateExchangeQuery = `UPDATE "exchanges"
SET trading_fee_maker = COALESCE(:trading_fee_maker, trading_fee_maker),
    trading_fee_taker = COALESCE(:trading_fee_taker, trading_fee_taker),
    disable           = COALESCE(:disable, disable)
WHERE id = :id
`
	updateExchange, err := db.PrepareNamed(updateExchangeQuery)
	if err != nil {
		return nil, err
	}

	const newAssetQuery = `SELECT new_asset
FROM new_asset(
             :symbol,
             :name,
             :address,
             :decimals,
             :transferable,
             :set_rate,
             :rebalance,
             :is_quote,
             :ask_a,
             :ask_b,
             :ask_c,
             :ask_min_min_spread,
             :ask_price_multiply_factor,
             :bid_a,
             :bid_b,
             :bid_c,
             :bid_min_min_spread,
             :bid_price_multiply_factor,
             :rebalance_quadratic_a,
             :rebalance_quadratic_b,
             :rebalance_quadratic_c,
             :target_total,
             :target_reserve,
             :target_rebalance_threshold,
             :target_transfer_threshold
         );`
	newAsset, err := db.PrepareNamed(newAssetQuery)
	if err != nil {
		return nil, err
	}

	const newAssetExchangeQuery string = `INSERT INTO asset_exchanges(exchange_id,
                            asset_id,
                            symbol,
                            deposit_address,
                            min_deposit,
                            withdraw_fee,
                            price_precision,
                            amount_precision,
                            amount_limit_min,
                            amount_limit_max,
                            price_limit_min,
                            price_limit_max,
                            target_recommended,
                            target_ratio)
VALUES (:exchange_id,
        :asset_id,
        :symbol,
        :deposit_address,
        :min_deposit,
        :withdraw_fee,
        :price_precision,
        :amount_precision,
        :amount_limit_min,
        :amount_limit_max,
        :price_limit_min,
        :price_limit_max,
        :target_recommended,
        :target_ratio) RETURNING id`
	newAssetExchange, err := db.PrepareNamed(newAssetExchangeQuery)
	if err != nil {
		return nil, err
	}

	const newTradingPairQuery = `SELECT new_trading_pair
FROM new_trading_pair($1, $2, $3, $4);`
	newOrderBook, err := db.Preparex(newTradingPairQuery)
	if err != nil {
		return nil, err
	}

	const getAssetQuery = `SELECT assets.id,
       assets.symbol,
       assets.name,
       a.address,
       array_agg(oa.address) FILTER ( WHERE oa.address IS NOT NULL ) AS old_addresses,
       assets.decimals,
       assets.transferable,
       assets.set_rate,
       assets.rebalance,
       assets.is_quote,
       assets.pwi_ask_a,
       assets.pwi_ask_b,
       assets.pwi_ask_c,
       assets.pwi_ask_min_min_spread,
       assets.pwi_ask_price_multiply_factor,
       assets.pwi_bid_a,
       assets.pwi_bid_b,
       assets.pwi_bid_c,
       assets.pwi_bid_min_min_spread,
       assets.pwi_bid_price_multiply_factor,
       assets.rebalance_quadratic_a,
       assets.rebalance_quadratic_b,
       assets.rebalance_quadratic_c,
       assets.target_total,
       assets.target_reserve,
       assets.target_rebalance_threshold,
       assets.target_transfer_threshold,
       assets.created,
       assets.updated
FROM assets
         LEFT JOIN addresses a on assets.address_id = a.id
         LEFT JOIN asset_old_addresses aoa on assets.id = aoa.asset_id
         LEFT JOIN addresses oa ON aoa.address_id = oa.id
WHERE assets.id = coalesce($1, assets.id)
GROUP BY assets.id,
         assets.symbol,
         assets.name,
         a.address,
         assets.decimals,
         assets.transferable,
         assets.set_rate,
         assets.rebalance,
         assets.is_quote,
         assets.pwi_ask_a,
         assets.pwi_ask_b,
         assets.pwi_ask_c,
         assets.pwi_ask_min_min_spread,
         assets.pwi_ask_price_multiply_factor,
         assets.pwi_bid_a,
         assets.pwi_bid_b,
         assets.pwi_bid_c,
         assets.pwi_bid_min_min_spread,
         assets.pwi_bid_price_multiply_factor,
         assets.rebalance_quadratic_a,
         assets.rebalance_quadratic_b,
         assets.rebalance_quadratic_c,
         assets.target_total,
         assets.target_reserve,
         assets.target_rebalance_threshold,
         assets.target_transfer_threshold,
         assets.created,
         assets.updated
ORDER BY assets.id`
	getAsset, err := db.Preparex(getAssetQuery)
	if err != nil {
		return nil, err
	}

	const getAssetExchangeQuery = `SELECT id,
       exchange_id,
       asset_id,
       symbol,
       deposit_address,
       min_deposit,
       withdraw_fee,
       price_precision,
       amount_precision,
       amount_limit_min,
       amount_limit_max,
       price_limit_min,
       price_limit_max,
       target_recommended,
       target_ratio
FROM asset_exchanges
WHERE asset_id = coalesce($1, asset_id)`
	getAssetExchange, err := db.Preparex(getAssetExchangeQuery)
	if err != nil {
		return nil, err
	}

	const getTradingPairQuery = `SELECT DISTINCT tp.id,
                tp.exchange_id,
                tp.base_id,
                tp.quote_id,
                tp.min_notional
FROM trading_pairs tp
         INNER JOIN asset_exchanges ae ON tp.exchange_id = ae.exchange_id AND ae.asset_id = coalesce($1, ae.asset_id);
`
	getOrderBook, err := db.Preparex(getTradingPairQuery)
	if err != nil {
		return nil, err
	}

	const updateAssetQuery = `WITH updated AS (
    UPDATE "addresses"
        SET address = COALESCE(:address, addresses.address)
        FROM "assets"
        WHERE assets.id = :id AND assets.address_id = addresses.id
)
UPDATE "assets"
SET symbol    = COALESCE(:symbol, symbol),
    name      = COALESCE(:name, name),
    decimals  = COALESCE(:decimals, decimals),
    set_rate  = COALESCE(:set_rate, set_rate),
    rebalance = COALESCE(:rebalance, rebalance),
    is_quote  = COALESCE(:is_quote, is_quote),
    updated   = now()
WHERE id = :id RETURNING id;
`
	updateAsset, err := db.PrepareNamed(updateAssetQuery)
	if err != nil {
		return nil, err
	}

	const changeAssetAddressQuery = `SELECT change_asset_address($1, $2);`
	changeAssetAddress, err := db.Preparex(changeAssetAddressQuery)
	if err != nil {
		return nil, err
	}

	const getMinNotionalQuery = `SELECT min_notional
FROM trading_pairs
WHERE exchange_id = $1
  AND base_id = $2
  AND quote_id = $3;
`
	getMinMotional, err := db.Preparex(getMinNotionalQuery)
	if err != nil {
		return nil, err
	}

	const getTradingPairSymbolsQuery = `SELECT DISTINCT bae.symbol AS base_symbol, qae.symbol AS quote_symbol
FROM trading_pairs AS tp
         INNER JOIN assets AS ba ON tp.base_id = ba.id
         INNER JOIN asset_exchanges AS bae ON ba.id = bae.asset_id
         INNER JOIN assets AS qa ON tp.quote_id = qa.id
         INNER JOIN asset_exchanges AS qae ON qa.id = qae.asset_id
WHERE bae.exchange_id = $1
  AND qae.exchange_id = $1; `
	getTradingPairSymbols, err := db.Preparex(getTradingPairSymbolsQuery)
	if err != nil {
		return nil, err
	}

	const updateDepositAddressQuery = `UPDATE asset_exchanges
SET deposit_address = $3
WHERE asset_id = $1
  AND exchange_id = $2 RETURNING id;`
	updateDepositAddress, err := db.Preparex(updateDepositAddressQuery)
	if err != nil {
		return nil, err
	}

	return &preparedStmts{
		getExchanges:     getExchanges,
		getExchange:      getExchange,
		updateExchange:   updateExchange,
		newAsset:         newAsset,
		newAssetExchange: newAssetExchange,
		newTradingPair:   newOrderBook,

		getAsset:             getAsset,
		getAssetExchange:     getAssetExchange,
		getTradingPair:       getOrderBook,
		updateAsset:          updateAsset,
		changeAssetAddress:   changeAssetAddress,
		updateDepositAddress: updateDepositAddress,

		getTradingPairSymbols: getTradingPairSymbols,
		getMinNotional:        getMinMotional,
		//getTransferableAssets:
	}, nil
}
