package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

type Match struct {
	ID              int        `json:"id"`
	MapName         string     `json:"mapName"`
	Opponent        string     `json:"opponent"`
	MatchDate       *time.Time `json:"matchDate,omitempty"`
	Status          *string    `json:"status,omitempty"`
	ScoreTeam       *int       `json:"scoreTeam,omitempty"`
	ScoreOpponent   *int       `json:"scoreOpponent,omitempty"`
	DurationSeconds *float64   `json:"durationSeconds,omitempty"`
}

type TeamStat struct {
	ID            int      `json:"id"`
	MatchesPlayed int      `json:"matchesPlayed"`
	WinRate       *float64 `json:"winRate,omitempty"`
	KD            *float64 `json:"kd,omitempty"`
	KAST          *float64 `json:"kast,omitempty"`
	KPR           *float64 `json:"kpr,omitempty"`
	FirstKills    *float64 `json:"firstKills,omitempty"`
	ADR           *float64 `json:"adr,omitempty"`
	DamageDelta   *float64 `json:"damageDelta,omitempty"`
	TradeKill     *float64 `json:"tradeKill,omitempty"`
	ClutchWR      *float64 `json:"clutchWR,omitempty"`
}

type SideStat struct {
	ID         int      `json:"id"`
	Side       string   `json:"side"`
	WR         *float64 `json:"wr,omitempty"`
	KD         *float64 `json:"kd,omitempty"`
	FirstKills *float64 `json:"firstKills,omitempty"`
	PlantWR    *float64 `json:"plantWR,omitempty"`
	HoldWR     *float64 `json:"holdWR,omitempty"`
	RetakeWR   *float64 `json:"retakeWR,omitempty"`
}

type Player struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Rating *int   `json:"rating,omitempty"`
	TeamID *int   `json:"teamId,omitempty"`
}

type PlayerStat struct {
	ID          int      `json:"id"`
	PlayerID    int      `json:"playerId"`
	Kills       int      `json:"kills"`
	Deaths      int      `json:"deaths"`
	Assists     int      `json:"assists"`
	KD          *float64 `json:"kd,omitempty"`
	KDA         *float64 `json:"kda,omitempty"`
	KAST        *float64 `json:"kast,omitempty"`
	ADR         *float64 `json:"adr,omitempty"`
	DamageDelta *float64 `json:"damageDelta,omitempty"`
	MultiKills  *int     `json:"multiKills,omitempty"`
	Clutch      *int     `json:"clutch,omitempty"`
	WinRate     *float64 `json:"winRate,omitempty"`
}

type fullResponse struct {
	Matches     []Match      `json:"matches"`
	TeamStats   []TeamStat   `json:"teamStats"`
	SideStats   []SideStat   `json:"sideStats"`
	Players     []Player     `json:"players"`
	PlayerStats []PlayerStat `json:"playerStats"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pool, err := getPool(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ensureSchema(ctx, pool); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := fetchAll(ctx, pool)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPool(ctx context.Context) (*pgxpool.Pool, error) {
	poolOnce.Do(func() {
		databaseURL := os.Getenv("DATABASE_URL")
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
		cfg.MaxConns = 4
		cfg.MaxConnLifetime = time.Minute * 5
		cfg.HealthCheckPeriod = time.Second * 30

		pool, poolErr = pgxpool.NewWithConfig(ctx, cfg)
	})

	return pool, poolErr
}

func ensureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS team_stats (
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
        )`,
		`CREATE TABLE IF NOT EXISTS matches (
            id SERIAL PRIMARY KEY,
            map_name TEXT NOT NULL,
            opponent TEXT NOT NULL,
            match_date TIMESTAMPTZ,
            status TEXT CHECK (status IN ('upcoming','victory','defeat')),
            score_team INT,
            score_opponent INT,
            duration INTERVAL
        )`,
		`CREATE TABLE IF NOT EXISTS side_stats (
            id SERIAL PRIMARY KEY,
            side TEXT CHECK (side IN ('attack','defense')),
            wr NUMERIC(5,2),
            kd NUMERIC(5,2),
            first_kills NUMERIC(5,2),
            plant_wr NUMERIC(5,2),
            hold_wr NUMERIC(5,2),
            retake_wr NUMERIC(5,2)
        )`,
		`CREATE TABLE IF NOT EXISTS players (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            role TEXT NOT NULL,
            rating INT,
            team_id INT,
            FOREIGN KEY (team_id) REFERENCES team_stats(id)
        )`,
		`CREATE TABLE IF NOT EXISTS player_stats (
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
        )`,
	}

	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

func fetchAll(ctx context.Context, pool *pgxpool.Pool) (fullResponse, error) {
	var result fullResponse

	matchesRows, err := pool.Query(ctx, `SELECT id, map_name, opponent, match_date, status, score_team, score_opponent, EXTRACT(EPOCH FROM duration) FROM matches ORDER BY match_date DESC NULLS LAST, id`)
	if err != nil {
		return result, err
	}
	defer matchesRows.Close()

	for matchesRows.Next() {
		var m Match
		var matchDate sql.NullTime
		var status sql.NullString
		var scoreTeam sql.NullInt32
		var scoreOpponent sql.NullInt32
		var duration sql.NullFloat64

		if err := matchesRows.Scan(&m.ID, &m.MapName, &m.Opponent, &matchDate, &status, &scoreTeam, &scoreOpponent, &duration); err != nil {
			return result, err
		}

		if matchDate.Valid {
			t := matchDate.Time
			m.MatchDate = &t
		}
		if status.Valid {
			s := status.String
			m.Status = &s
		}
		if scoreTeam.Valid {
			v := int(scoreTeam.Int32)
			m.ScoreTeam = &v
		}
		if scoreOpponent.Valid {
			v := int(scoreOpponent.Int32)
			m.ScoreOpponent = &v
		}
		if duration.Valid {
			v := duration.Float64
			m.DurationSeconds = &v
		}

		result.Matches = append(result.Matches, m)
	}
	if err := matchesRows.Err(); err != nil {
		return result, err
	}

	teamRows, err := pool.Query(ctx, `SELECT id, matches_played, win_rate, kd, kast, kpr, first_kills, adr, dmg_delta, trade_kill, clutch_wr FROM team_stats ORDER BY id`)
	if err != nil {
		return result, err
	}
	defer teamRows.Close()

	for teamRows.Next() {
		var t TeamStat
		var winRate, kd, kast, kpr, firstKills, adr, dmgDelta, tradeKill, clutchWR sql.NullFloat64

		if err := teamRows.Scan(&t.ID, &t.MatchesPlayed, &winRate, &kd, &kast, &kpr, &firstKills, &adr, &dmgDelta, &tradeKill, &clutchWR); err != nil {
			return result, err
		}

		t.WinRate = nullableFloat(winRate)
		t.KD = nullableFloat(kd)
		t.KAST = nullableFloat(kast)
		t.KPR = nullableFloat(kpr)
		t.FirstKills = nullableFloat(firstKills)
		t.ADR = nullableFloat(adr)
		t.DamageDelta = nullableFloat(dmgDelta)
		t.TradeKill = nullableFloat(tradeKill)
		t.ClutchWR = nullableFloat(clutchWR)

		result.TeamStats = append(result.TeamStats, t)
	}
	if err := teamRows.Err(); err != nil {
		return result, err
	}

	sideRows, err := pool.Query(ctx, `SELECT id, side, wr, kd, first_kills, plant_wr, hold_wr, retake_wr FROM side_stats ORDER BY id`)
	if err != nil {
		return result, err
	}
	defer sideRows.Close()

	for sideRows.Next() {
		var s SideStat
		var wr, kd, firstKills, plantWR, holdWR, retakeWR sql.NullFloat64

		if err := sideRows.Scan(&s.ID, &s.Side, &wr, &kd, &firstKills, &plantWR, &holdWR, &retakeWR); err != nil {
			return result, err
		}

		s.WR = nullableFloat(wr)
		s.KD = nullableFloat(kd)
		s.FirstKills = nullableFloat(firstKills)
		s.PlantWR = nullableFloat(plantWR)
		s.HoldWR = nullableFloat(holdWR)
		s.RetakeWR = nullableFloat(retakeWR)

		result.SideStats = append(result.SideStats, s)
	}
	if err := sideRows.Err(); err != nil {
		return result, err
	}

	playerRows, err := pool.Query(ctx, `SELECT id, name, role, rating, team_id FROM players ORDER BY id`)
	if err != nil {
		return result, err
	}
	defer playerRows.Close()

	for playerRows.Next() {
		var p Player
		var rating sql.NullInt32
		var teamID sql.NullInt32

		if err := playerRows.Scan(&p.ID, &p.Name, &p.Role, &rating, &teamID); err != nil {
			return result, err
		}

		if rating.Valid {
			v := int(rating.Int32)
			p.Rating = &v
		}
		if teamID.Valid {
			v := int(teamID.Int32)
			p.TeamID = &v
		}

		result.Players = append(result.Players, p)
	}
	if err := playerRows.Err(); err != nil {
		return result, err
	}

	playerStatRows, err := pool.Query(ctx, `SELECT id, player_id, kills, deaths, assists, kd, kda, kast, adr, dmg_delta, mk, clutch, wr FROM player_stats ORDER BY id`)
	if err != nil {
		return result, err
	}
	defer playerStatRows.Close()

	for playerStatRows.Next() {
		var ps PlayerStat
		var kd, kda, kast, adr, dmgDelta, wr sql.NullFloat64
		var mk, clutch sql.NullInt32

		if err := playerStatRows.Scan(&ps.ID, &ps.PlayerID, &ps.Kills, &ps.Deaths, &ps.Assists, &kd, &kda, &kast, &adr, &dmgDelta, &mk, &clutch, &wr); err != nil {
			return result, err
		}

		ps.KD = nullableFloat(kd)
		ps.KDA = nullableFloat(kda)
		ps.KAST = nullableFloat(kast)
		ps.ADR = nullableFloat(adr)
		ps.DamageDelta = nullableFloat(dmgDelta)
		ps.WinRate = nullableFloat(wr)
		if mk.Valid {
			v := int(mk.Int32)
			ps.MultiKills = &v
		}
		if clutch.Valid {
			v := int(clutch.Int32)
			ps.Clutch = &v
		}

		result.PlayerStats = append(result.PlayerStats, ps)
	}
	if err := playerStatRows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func nullableFloat(value sql.NullFloat64) *float64 {
	if value.Valid {
		v := value.Float64
		return &v
	}
	return nil
}
