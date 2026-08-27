-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS kitchen;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION kitchen.set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE kitchen.restaurants (
    id              UUID            PRIMARY KEY     DEFAULT gen_random_uuid(),
    name            TEXT            NOT NULL,
    description     TEXT,
    address         JSONB           NOT NULL,
    status          TEXT            NOT NULL        DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'blocked')),
    service_url     TEXT            NOT NULL        DEFAULT '',
    created_at      TIMESTAMPTZ     NOT NULL        DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL        DEFAULT now(),

    CONSTRAINT valid_address CHECK (
        jsonb_typeof(address) = 'object'
        AND address ? 'city'
        AND address ? 'street'
        AND address ? 'house'
        AND address ? 'apartment'
        AND jsonb_typeof(address->'city') = 'string'
        AND jsonb_typeof(address->'street') = 'string'
        AND jsonb_typeof(address->'house') = 'string'
        AND jsonb_typeof(address->'apartment') = 'string'
    )
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER restaurants_set_updated_at
    BEFORE UPDATE ON kitchen.restaurants
    FOR EACH ROW EXECUTE FUNCTION kitchen.set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE kitchen.categories (
    id              UUID            PRIMARY KEY     DEFAULT gen_random_uuid(),
    name            TEXT            NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL        DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL        DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER categories_set_updated_at
    BEFORE UPDATE ON kitchen.categories
    FOR EACH ROW EXECUTE FUNCTION kitchen.set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE kitchen.menu_items (
    id              UUID            PRIMARY KEY     DEFAULT gen_random_uuid(),
    restaurant_id   UUID            NOT NULL        REFERENCES kitchen.restaurants (id) ON DELETE CASCADE,
    category_id     UUID                            REFERENCES kitchen.categories (id) ON DELETE SET NULL,
    name            TEXT            NOT NULL,
    description     TEXT,
    price           NUMERIC(12, 2)  NOT NULL        CHECK (price >= 0),
    available       BOOLEAN         NOT NULL        DEFAULT TRUE,
    created_at      TIMESTAMPTZ     NOT NULL        DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL        DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_menu_items_restaurant_id ON kitchen.menu_items (restaurant_id);
CREATE INDEX idx_menu_items_category_id ON kitchen.menu_items (category_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER menu_items_set_updated_at
    BEFORE UPDATE ON kitchen.menu_items
    FOR EACH ROW EXECUTE FUNCTION kitchen.set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE kitchen.orders (
    id                UUID              PRIMARY KEY     DEFAULT gen_random_uuid(),
    restaurant_id     UUID              NOT NULL        REFERENCES kitchen.restaurants (id),
    user_id           TEXT,
    delivery_address  JSONB             NOT NULL,
    status            TEXT              NOT NULL        DEFAULT 'pending' CHECK (status IN (
        'pending', 'sent_to_restaurant', 'accepted', 'preparing',
        'ready', 'in_delivery', 'delivered', 'rejected_by_restaurant'
    )),
    total_amount      NUMERIC(12, 2)    NOT NULL        CHECK (total_amount >= 0),
    comment           TEXT,
    created_at        TIMESTAMPTZ       NOT NULL        DEFAULT now(),
    updated_at        TIMESTAMPTZ       NOT NULL        DEFAULT now(),
    
    CONSTRAINT valid_delivery_address CHECK (
        jsonb_typeof(delivery_address) = 'object'
        AND delivery_address ? 'city'
        AND delivery_address ? 'street'
        AND delivery_address ? 'house'
        AND delivery_address ? 'apartment'
        AND jsonb_typeof(delivery_address->'city') = 'string'
        AND jsonb_typeof(delivery_address->'street') = 'string'
        AND jsonb_typeof(delivery_address->'house') = 'string'
        AND jsonb_typeof(delivery_address->'apartment') = 'string'
    )
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_orders_restaurant_id_status ON kitchen.orders (restaurant_id, status);
CREATE INDEX idx_orders_created_at ON kitchen.orders (created_at);
CREATE INDEX idx_orders_user_id ON kitchen.orders (user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER orders_set_updated_at
    BEFORE UPDATE ON kitchen.orders
    FOR EACH ROW EXECUTE FUNCTION kitchen.set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE kitchen.order_items (
    id                  UUID            PRIMARY KEY     DEFAULT gen_random_uuid(),
    order_id            UUID            NOT NULL        REFERENCES kitchen.orders (id) ON DELETE CASCADE,
    menu_item_id        UUID            NOT NULL        REFERENCES kitchen.menu_items (id),
    name                TEXT            NOT NULL,
    price               NUMERIC(12, 2)  NOT NULL        CHECK (price >= 0),
    quantity            INTEGER         NOT NULL        CHECK (quantity > 0)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_order_items_order_id ON kitchen.order_items (order_id);
CREATE INDEX idx_order_items_menu_item_id ON kitchen.order_items (menu_item_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kitchen.order_items;
DROP TABLE IF EXISTS kitchen.orders;
DROP TABLE IF EXISTS kitchen.menu_items;
DROP TABLE IF EXISTS kitchen.categories;
DROP TABLE IF EXISTS kitchen.restaurants;
DROP FUNCTION IF EXISTS kitchen.set_updated_at();
DROP SCHEMA IF EXISTS kitchen;
-- +goose StatementEnd
