'use client';

import { useState } from 'react';

type Match = {
  id: number;
  mapName: string;
  opponent: string;
  matchDate?: string;
  status?: string;
  scoreTeam?: number;
  scoreOpponent?: number;
  durationSeconds?: number;
};

type TeamStat = {
  id: number;
  matchesPlayed: number;
  winRate?: number;
  kd?: number;
  kast?: number;
  kpr?: number;
  firstKills?: number;
  adr?: number;
  damageDelta?: number;
  tradeKill?: number;
  clutchWR?: number;
};

type SideStat = {
  id: number;
  side: string;
  wr?: number;
  kd?: number;
  firstKills?: number;
  plantWR?: number;
  holdWR?: number;
  retakeWR?: number;
};

type Player = {
  id: number;
  name: string;
  role: string;
  rating?: number;
  teamId?: number;
};

type PlayerStat = {
  id: number;
  playerId: number;
  kills: number;
  deaths: number;
  assists: number;
  kd?: number;
  kda?: number;
  kast?: number;
  adr?: number;
  damageDelta?: number;
  multiKills?: number;
  clutch?: number;
  winRate?: number;
};

type ApiResponse = {
  matches: Match[];
  teamStats: TeamStat[];
  sideStats: SideStat[];
  players: Player[];
  playerStats: PlayerStat[];
};

const formatNumber = (value?: number) =>
  value === undefined || Number.isNaN(value) ? '–' : value.toFixed(2).replace(/\.00$/, '');

const formatDuration = (seconds?: number) => {
  if (seconds === undefined || Number.isNaN(seconds)) return '–';
  const total = Math.round(seconds);
  const mins = Math.floor(total / 60);
  const secs = total % 60;
  return `${mins}m ${secs}s`;
};

const formatDate = (iso?: string) => (iso ? new Date(iso).toLocaleString() : '–');

export default function Home() {
  const [data, setData] = useState<ApiResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleQuery = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch('/api/query-data');
      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
      }

      const payload = (await response.json()) as ApiResponse;
      setData(payload);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      setError(message);
      setData(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center gap-10 bg-slate-950 px-6 py-12 text-slate-50">
      <header className="flex flex-col items-center gap-4 text-center">
        <h1 className="text-3xl font-semibold">Neon Match Dashboard</h1>
        <p className="max-w-xl text-slate-300">
          Load the latest matches, team metrics, and player statistics directly from your Neon database.
        </p>
        <button
          type="button"
          onClick={handleQuery}
          disabled={loading}
          className="rounded bg-emerald-500 px-5 py-2 font-medium text-white transition hover:bg-emerald-400 disabled:cursor-not-allowed disabled:bg-emerald-700/60"
        >
          {loading ? 'Fetching…' : data ? 'Refresh data' : 'Load data'}
        </button>
      </header>

      {error && <p className="text-red-400">{error}</p>}

      {!data && !error && !loading && (
        <p className="text-slate-400">Press the button to fetch match and team data.</p>
      )}

      {data && (
        <div className="flex w-full max-w-5xl flex-col gap-8">
          <section className="space-y-3">
            <h2 className="text-2xl font-semibold text-emerald-300">Matches</h2>
            {data.matches.length === 0 ? (
              <p className="text-slate-400">No matches found.</p>
            ) : (
              <ul className="space-y-3">
                {data.matches.map((match) => (
                  <li key={match.id} className="rounded border border-emerald-500/40 bg-emerald-500/10 p-4">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <p className="text-lg font-semibold text-white">
                        {match.mapName} vs {match.opponent}
                      </p>
                      <span className="rounded-full bg-slate-900 px-3 py-1 text-sm uppercase tracking-wide text-emerald-200">
                        {match.status ?? 'unknown'}
                      </span>
                    </div>
                    <div className="mt-2 grid gap-1 text-sm text-slate-200 sm:grid-cols-2">
                      <span>When: {formatDate(match.matchDate)}</span>
                      <span>
                        Score: {match.scoreTeam ?? '–'} - {match.scoreOpponent ?? '–'}
                      </span>
                      <span>Duration: {formatDuration(match.durationSeconds)}</span>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </section>

          <section className="space-y-3">
            <h2 className="text-2xl font-semibold text-emerald-300">Team Stats</h2>
            {data.teamStats.length === 0 ? (
              <p className="text-slate-400">No team stats available.</p>
            ) : (
              <div className="overflow-hidden rounded border border-slate-800">
                <table className="w-full table-auto border-collapse text-sm">
                  <thead className="bg-slate-900 text-slate-200">
                    <tr>
                      <th className="px-3 py-2 text-left">Matches</th>
                      <th className="px-3 py-2 text-left">Win %</th>
                      <th className="px-3 py-2 text-left">KD</th>
                      <th className="px-3 py-2 text-left">KAST</th>
                      <th className="px-3 py-2 text-left">KPR</th>
                      <th className="px-3 py-2 text-left">ADR</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.teamStats.map((stat) => (
                      <tr key={stat.id} className="odd:bg-slate-900/40">
                        <td className="px-3 py-2">{stat.matchesPlayed}</td>
                        <td className="px-3 py-2">{formatNumber(stat.winRate)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.kd)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.kast)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.kpr)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.adr)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>

          <section className="space-y-3">
            <h2 className="text-2xl font-semibold text-emerald-300">Side Stats</h2>
            {data.sideStats.length === 0 ? (
              <p className="text-slate-400">No side stats available.</p>
            ) : (
              <div className="grid gap-3 sm:grid-cols-2">
                {data.sideStats.map((stat) => (
                  <div key={stat.id} className="rounded border border-slate-800 bg-slate-900/60 p-4">
                    <p className="text-lg font-semibold capitalize text-white">{stat.side}</p>
                    <dl className="mt-2 space-y-1 text-sm text-slate-300">
                      <div className="flex justify-between">
                        <dt>Win %</dt>
                        <dd>{formatNumber(stat.wr)}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt>KD</dt>
                        <dd>{formatNumber(stat.kd)}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt>First Kills %</dt>
                        <dd>{formatNumber(stat.firstKills)}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt>Plant WR</dt>
                        <dd>{formatNumber(stat.plantWR)}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt>Hold WR</dt>
                        <dd>{formatNumber(stat.holdWR)}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt>Retake WR</dt>
                        <dd>{formatNumber(stat.retakeWR)}</dd>
                      </div>
                    </dl>
                  </div>
                ))}
              </div>
            )}
          </section>

          <section className="space-y-3">
            <h2 className="text-2xl font-semibold text-emerald-300">Players</h2>
            {data.players.length === 0 ? (
              <p className="text-slate-400">No players found.</p>
            ) : (
              <ul className="grid gap-3 sm:grid-cols-2">
                {data.players.map((player) => (
                  <li key={player.id} className="rounded border border-slate-800 bg-slate-900/60 p-4">
                    <p className="text-lg font-semibold text-white">{player.name}</p>
                    <p className="text-sm uppercase tracking-wide text-emerald-200">{player.role}</p>
                    <p className="mt-2 text-sm text-slate-300">Rating: {player.rating ?? '–'}</p>
                    <p className="text-sm text-slate-300">Team ID: {player.teamId ?? '–'}</p>
                  </li>
                ))}
              </ul>
            )}
          </section>

          <section className="space-y-3">
            <h2 className="text-2xl font-semibold text-emerald-300">Player Stats</h2>
            {data.playerStats.length === 0 ? (
              <p className="text-slate-400">No player stats recorded.</p>
            ) : (
              <div className="overflow-x-auto rounded border border-slate-800">
                <table className="min-w-[720px] table-auto border-collapse text-sm">
                  <thead className="bg-slate-900 text-slate-200">
                    <tr>
                      <th className="px-3 py-2 text-left">Player ID</th>
                      <th className="px-3 py-2 text-left">Kills</th>
                      <th className="px-3 py-2 text-left">Deaths</th>
                      <th className="px-3 py-2 text-left">Assists</th>
                      <th className="px-3 py-2 text-left">K/D</th>
                      <th className="px-3 py-2 text-left">KDA</th>
                      <th className="px-3 py-2 text-left">KAST</th>
                      <th className="px-3 py-2 text-left">ADR</th>
                      <th className="px-3 py-2 text-left">Win %</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.playerStats.map((stat) => (
                      <tr key={stat.id} className="odd:bg-slate-900/40">
                        <td className="px-3 py-2">{stat.playerId}</td>
                        <td className="px-3 py-2">{stat.kills}</td>
                        <td className="px-3 py-2">{stat.deaths}</td>
                        <td className="px-3 py-2">{stat.assists}</td>
                        <td className="px-3 py-2">{formatNumber(stat.kd)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.kda)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.kast)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.adr)}</td>
                        <td className="px-3 py-2">{formatNumber(stat.winRate)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>
        </div>
      )}
    </main>
  );
}
