-- Schema setup for Neon database
CREATE TABLE IF NOT EXISTS team_stats (
    id SERIAL PRIMARY KEY,
    matches_played INT,
    win_rate NUMERIC(5,2),
    kd NUMERIC(5,2),
    kast NUMERIC(5,2),
    kpr NUMERIC(5,2),
    first_kills NUMERIC(5,2),
    adr NUMERIC(6,2),
    dmg_delta NUMERIC(6,2),
    trade_kill NUMERIC(5,2),
    clutch_wr NUMERIC(5,2)
);

CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    map_name TEXT NOT NULL,
    opponent TEXT NOT NULL,
    match_date TIMESTAMPTZ,
    status TEXT CHECK (status IN ('upcoming','victory','defeat')),
    score_team INT,
    score_opponent INT,
    duration INTERVAL
);

CREATE TABLE IF NOT EXISTS side_stats (
    id SERIAL PRIMARY KEY,
    side TEXT CHECK (side IN ('attack','defense')),
    wr NUMERIC(5,2),
    kd NUMERIC(5,2),
    first_kills NUMERIC(5,2),
    plant_wr NUMERIC(5,2),
    hold_wr NUMERIC(5,2),
    retake_wr NUMERIC(5,2)
);

CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    rating INT,
    team_id INT,
    FOREIGN KEY (team_id) REFERENCES team_stats(id)
);

CREATE TABLE IF NOT EXISTS player_stats (
    id SERIAL PRIMARY KEY,
    player_id INT REFERENCES players(id),
    kills INT,
    deaths INT,
    assists INT,
    kd NUMERIC(5,2),
    kda NUMERIC(5,2),
    kast NUMERIC(5,2),
    adr NUMERIC(6,2),
    dmg_delta NUMERIC(6,2),
    mk INT,
    clutch INT,
    wr NUMERIC(5,2)
);

-- Sample seed data
INSERT INTO team_stats (matches_played, win_rate, kd, kast, kpr, first_kills, adr, dmg_delta, trade_kill, clutch_wr)
VALUES (25, 64.00, 1.12, 72.50, 0.92, 52.10, 145.30, 15.40, 48.20, 36.00)
ON CONFLICT DO NOTHING;

INSERT INTO matches (map_name, opponent, match_date, status, score_team, score_opponent, duration)
VALUES
    ('Ascent', 'Team Ember', NOW() - INTERVAL '3 days', 'victory', 13, 9, INTERVAL '38 minutes'),
    ('Bind', 'Nova Six', NOW() - INTERVAL '10 days', 'defeat', 8, 13, INTERVAL '34 minutes')
ON CONFLICT DO NOTHING;

INSERT INTO side_stats (side, wr, kd, first_kills, plant_wr, hold_wr, retake_wr)
VALUES
    ('attack', 58.20, 1.05, 49.80, 62.50, 54.30, 0),
    ('defense', 68.40, 1.19, 55.10, 0, 66.20, 48.70)
ON CONFLICT DO NOTHING;

INSERT INTO players (name, role, rating)
VALUES
    ('Aspas', 'Duelist', 92),
    ('Less', 'Sentinel', 88)
ON CONFLICT DO NOTHING;

INSERT INTO player_stats (player_id, kills, deaths, assists, kd, kda, kast, adr, dmg_delta, mk, clutch, wr)
SELECT id, 245, 210, 98, 1.17, 1.45, 73.20, 151.60, 22.40, 12, 3, 62.00
FROM players
WHERE name = 'Aspas'
ON CONFLICT DO NOTHING;

INSERT INTO player_stats (player_id, kills, deaths, assists, kd, kda, kast, adr, dmg_delta, mk, clutch, wr)
SELECT id, 198, 185, 110, 1.07, 1.41, 76.10, 139.20, 14.90, 9, 4, 58.00
FROM players
WHERE name = 'Less'
ON CONFLICT DO NOTHING;
