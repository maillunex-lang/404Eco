INSERT INTO world_state (id, tick, day, month, year, oil_price, grain_index, labor_wage)
VALUES (1, 0, 1, 3, 2026, 78.5, 100, 48000)
ON CONFLICT (id) DO NOTHING;

-- Cities of the Republic of Bashkortostan (map coords ~0..1000)
INSERT INTO cities (id, name_ru, name_ba, x, y, population, wealth, industry_bias) VALUES
('ufa', 'Уфа', 'Өфө', 480, 320, 1140000, 1.35,
 '{"retail":1.4,"oil":1.2,"services":1.5,"food":1.2,"construction":1.3}'),
('sterlitamak', 'Стерлитамак', 'Стәрлетамаҡ', 520, 520, 275000, 1.1,
 '{"oil":1.5,"chemical":1.6,"retail":1.0,"manufacturing":1.3}'),
('salavat', 'Салават', 'Салауат', 500, 560, 148000, 1.05,
 '{"oil":1.8,"chemical":1.7,"manufacturing":1.2}'),
('neftekamsk', 'Нефтекамск', 'Нефтекама', 280, 180, 140000, 1.0,
 '{"oil":1.4,"manufacturing":1.5,"retail":0.95}'),
('oktyabrsky', 'Октябрьский', 'Октябрьский', 200, 380, 112000, 0.95,
 '{"oil":1.6,"services":1.0,"retail":0.9}'),
('beloretsk', 'Белорецк', 'Белорет', 720, 480, 64000, 0.85,
 '{"mining":1.6,"manufacturing":1.3,"tourism":1.2}'),
('tuymazy', 'Туймазы', 'Туймазы', 220, 340, 68000, 0.9,
 '{"oil":1.5,"agriculture":1.2,"food":1.1}'),
('ishimbay', 'Ишимбай', 'Ишембай', 490, 540, 64000, 0.9,
 '{"oil":1.7,"manufacturing":1.1}'),
('sibay', 'Сибай', 'Сибай', 820, 620, 60000, 0.8,
 '{"mining":1.8,"manufacturing":1.1}'),
('kumertau', 'Кумертау', 'Күмертау', 460, 680, 59000, 0.85,
 '{"manufacturing":1.2,"agriculture":1.3,"aviation":1.4}'),
('birsk', 'Бирск', 'Бөрө', 420, 200, 45000, 0.9,
 '{"agriculture":1.4,"food":1.3,"retail":0.95}'),
('belebey', 'Белебей', 'Бәләбәй', 260, 420, 58000, 0.88,
 '{"agriculture":1.3,"food":1.2,"manufacturing":1.0}')
ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, name_ru, category, base_price, unit) VALUES
('crude_oil', 'Нефть сырая', 'oil', 5200, 'т'),
('gasoline', 'Бензин', 'oil', 58, 'л'),
('plastics', 'Пластмасса', 'chemical', 95, 'кг'),
('wheat', 'Пшеница', 'agriculture', 18, 'кг'),
('honey', 'Башкирский мёд', 'agriculture', 650, 'кг'),
('milk', 'Молоко', 'agriculture', 55, 'л'),
('bread', 'Хлеб', 'food', 45, 'шт'),
('dairy', 'Молочные продукты', 'food', 120, 'кг'),
('copper_ore', 'Медная руда', 'mining', 280, 'т'),
('steel', 'Сталь', 'manufacturing', 42000, 'т'),
('auto_parts', 'Автозапчасти', 'manufacturing', 8500, 'комплект'),
('furniture', 'Мебель', 'manufacturing', 25000, 'шт'),
('retail_goods', 'Товары народного потребления', 'retail', 500, 'лот'),
('construction', 'Стройматериалы', 'construction', 15000, 'лот'),
('tourism_pkg', 'Туристический пакет', 'tourism', 12000, 'тур')
ON CONFLICT (id) DO NOTHING;

-- Available properties (unowned at start)
INSERT INTO properties (city_id, kind, name, size, base_rent, buy_price, condition)
SELECT c.id, k.kind, k.name || ' · ' || c.name_ru, k.size, k.rent * c.wealth, k.buy * c.wealth, 75 + random()*20
FROM cities c
CROSS JOIN (VALUES
  ('retail', 'Киоск у остановки', 1, 18000, 650000),
  ('retail', 'Торговая площадь', 1, 45000, 1800000),
  ('retail', 'ТЦ «Ярмолка»', 2, 90000, 4500000),
  ('factory', 'Мастерская', 1, 55000, 2200000),
  ('factory', 'Промплощадка', 2, 120000, 6500000),
  ('farm', 'Пасека / мини-ферма', 1, 15000, 900000),
  ('farm', 'Агроферма', 3, 35000, 2800000),
  ('oil', 'Нефтеучасток', 2, 200000, 12000000),
  ('warehouse', 'Склад', 2, 40000, 2100000),
  ('office', 'Офис', 1, 55000, 1600000),
  ('mine', 'Карьер', 3, 150000, 9500000)
) AS k(kind, name, size, rent, buy)
WHERE
  (k.kind = 'oil' AND c.id IN ('ufa','sterlitamak','salavat','neftekamsk','oktyabrsky','tuymazy','ishimbay'))
  OR (k.kind = 'mine' AND c.id IN ('beloretsk','sibay'))
  OR (k.kind IN ('farm') AND c.id IN ('birsk','belebey','tuymazy','kumertau','ufa'))
  OR (k.kind = 'retail' AND k.name = 'ТЦ «Ярмолка»' AND c.id IN ('ufa','sterlitamak','neftekamsk'))
  OR (k.kind NOT IN ('oil','mine','farm') AND NOT (k.kind = 'retail' AND k.name = 'ТЦ «Ярмолка»'));
