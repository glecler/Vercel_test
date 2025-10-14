package handler

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "os"
    "sync"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

const (
    keyID       = "id_237429"
    defaultWord = "Bingo"
)

var (
    pool     *pgxpool.Pool
    poolOnce sync.Once
    poolErr  error
)

type keyRecord struct {
    ID   string `json:"id"`
    Word string `json:"word"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.Header().Set("Allow", http.MethodGet)
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    pool, err := getPool(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    record, err := ensureKey(ctx, pool)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(record); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func getPool(ctx context.Context) (*pgxpool.Pool, error) {
    poolOnce.Do(func() {
        databaseURL := os.Getenv("NEON_DATABASE_URL")
        if databaseURL == "" {
            databaseURL = os.Getenv("DATABASE_URL")
        }
        if databaseURL == "" {
            poolErr = errors.New("missing NEON_DATABASE_URL or DATABASE_URL environment variable")
            return
        }

        cfg, err := pgxpool.ParseConfig(databaseURL)
        if err != nil {
            poolErr = err
            return
        }

        cfg.MinConns = 0
        cfg.MaxConns = 2
        cfg.MaxConnLifetime = time.Minute * 5
        cfg.HealthCheckPeriod = time.Second * 30

        pool, poolErr = pgxpool.NewWithConfig(ctx, cfg)
    })

    return pool, poolErr
}

func ensureKey(ctx context.Context, pool *pgxpool.Pool) (keyRecord, error) {
    if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS keys (
        id TEXT PRIMARY KEY,
        word TEXT NOT NULL
    )`); err != nil {
        return keyRecord{}, err
    }

    var record keyRecord
    err := pool.QueryRow(ctx, "SELECT id, word FROM keys WHERE id = $1", keyID).Scan(&record.ID, &record.Word)
    if err == nil {
        return record, nil
    }

    if !errors.Is(err, pgx.ErrNoRows) {
        return keyRecord{}, err
    }

    record = keyRecord{ID: keyID, Word: defaultWord}
    if _, err := pool.Exec(ctx, `INSERT INTO keys (id, word) VALUES ($1, $2)
        ON CONFLICT (id) DO UPDATE SET word = EXCLUDED.word`, record.ID, record.Word); err != nil {
        return keyRecord{}, err
    }

    return record, nil
}
