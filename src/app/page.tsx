'use client';

import { useState } from 'react';

type KeyRecord = {
  id: string;
  word: string;
};

export default function Home() {
  const [record, setRecord] = useState<KeyRecord | null>(null);
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

      const data = (await response.json()) as KeyRecord;
      setRecord(data);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      setError(message);
      setRecord(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-6 bg-slate-950 p-8 text-slate-50">
      <h1 className="text-3xl font-semibold">Neon Key Fetch Demo</h1>
      <button
        type="button"
        onClick={handleQuery}
        disabled={loading}
        className="rounded bg-emerald-500 px-5 py-2 font-medium text-white transition hover:bg-emerald-400 disabled:cursor-not-allowed disabled:bg-emerald-700/60"
      >
        {loading ? 'Querying…' : 'Query data'}
      </button>

      {error ? (
        <p className="text-red-400">{error}</p>
      ) : record ? (
        <div className="rounded border border-emerald-500/50 bg-emerald-500/10 p-4 text-center">
          <p className="text-sm uppercase tracking-wider text-emerald-300">Current key</p>
          <p className="text-lg font-semibold">{record.id}</p>
          <p className="text-base">{record.word}</p>
        </div>
      ) : (
        <p className="text-slate-400">Press the button to fetch the key.</p>
      )}
    </main>
  );
}
