-- +goose Up
INSERT INTO kitchen.restaurants (id, name, description, address, status, service_url) VALUES
(
    '123e4567-e89b-12d3-a456-426614174000',
    'Вкусно и точка',
    'Ресторан русской кухни с домашними блюдами',
    '{
        "city": "Москва",
        "street": "Тверская",
        "house": "12",
        "apartment": "1",
        "comment": "Вход со двора"
    }'::jsonb,
    'active',
    'http://restaurant-simulator:8081/api/v1'
),
(
    '123e4567-e89b-12d3-a456-426614174001',
    'Итальянский дворик',
    'Уютный ресторан итальянской кухни',
    '{
        "city": "Санкт-Петербург",
        "street": "Невский проспект",
        "house": "45",
        "apartment": "2",
        "comment": "2 этаж"
    }'::jsonb,
    'active',
    'http://localhost:8082'
);

-- +goose Down
DELETE FROM kitchen.restaurants 
WHERE id IN (
    '123e4567-e89b-12d3-a456-426614174000',
    '123e4567-e89b-12d3-a456-426614174001'
);

