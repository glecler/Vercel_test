import { neon } from '@neondatabase/serverless';

const sql = neon(process.env.DATABASE_URL);

function choice(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

function rand(min, max, dec = 2) {
  return parseFloat((Math.random() * (max - min) + min).toFixed(dec));
}

async function seed() {
  // Clear existing data
  await sql`TRUNCATE player_stats, players, side_stats, team_stats, matches RESTART IDENTITY CASCADE;`;

  // Matches
  const maps = ['Haven', 'Bind', 'Ascent', 'Sunset', 'Lotus'];
  const opponents = ['FNATIC', 'G2 ESPORTS', 'LOUD', 'SENTINELS', 'TEAM LIQUID'];
  const statuses = ['upcoming', 'victory', 'defeat'];

  for (let i = 0; i < 10; i++) {
    await sql`
      INSERT INTO matches (map_name, opponent, match_date, status, score_team, score_opponent, duration)
      VALUES (
        ${choice(maps)},
        ${choice(opponents)},
        ${new Date(Date.now() + i * 86400000).toISOString()},
        ${choice(statuses)},
        ${10 + Math.floor(Math.random() * 4)},
        ${7 + Math.floor(Math.random() * 5)},
        ${`${35 + Math.floor(Math.random() * 10)} minutes`}
      );
    `;
  }

  // Team stats
  for (let i = 0; i < 2; i++) {
    await sql`
      INSERT INTO team_stats (
        matches_played, win_rate, kd, kast, kpr, first_kills, adr,
        dmg_delta, trade_kill, clutch_wr
      )
      VALUES (
        ${100 + Math.floor(Math.random() * 50)},
        ${rand(60, 75)},
        ${rand(1.6, 2.0)},
        ${rand(74, 78)},
        ${rand(3.5, 4.0)},
        ${rand(60, 65)},
        ${rand(150, 180)},
        ${rand(10, 30)},
        ${rand(20, 25)},
        ${rand(15, 20)}
      );
    `;
  }

  // Side stats
  const sides = ['attack', 'defense'];
  for (const side of sides) {
    await sql`
      INSERT INTO side_stats (side, wr, kd, first_kills, plant_wr, hold_wr, retake_wr)
      VALUES (
        ${side},
        ${rand(65, 72)},
        ${rand(1.6, 2.0)},
        ${rand(55, 70)},
        ${rand(50, 80)},
        ${rand(50, 75)},
        ${rand(30, 35)}
      );
    `;
  }

  // Players
  const players = [
    ['Aspas', 'Duelist'],
    ['Demon1', 'Duelist'],
    ['Chronicle', 'Initiator'],
    ['Saadhak', 'Controller'],
    ['Less', 'Sentinel'],
    ['Leo', 'Initiator'],
    ['Derke', 'Duelist'],
    ['Boaster', 'Controller'],
    ['Alfajer', 'Sentinel'],
    ['Nats', 'Sentinel']
  ];

  for (let i = 0; i < players.length; i++) {
    await sql`
      INSERT INTO players (name, role, rating, team_id)
      VALUES (
        ${players[i][0]},
        ${players[i][1]},
        ${75 + Math.floor(Math.random() * 20)},
        ${1 + (i % 2)}
      );
    `;
  }

  // Player stats
  for (let i = 0; i < 10; i++) {
    await sql`
      INSERT INTO player_stats (
        player_id, kills, deaths, assists, kd, kda, kast, adr,
        dmg_delta, mk, clutch, wr
      )
      VALUES (
        ${i + 1},
        ${1900 + Math.floor(Math.random() * 500)},
        ${1000 + Math.floor(Math.random() * 300)},
        ${600 + Math.floor(Math.random() * 200)},
        ${rand(1.6, 2.1)},
        ${rand(2.2, 2.8)},
        ${rand(73, 79)},
        ${rand(150, 175)},
        ${rand(0, 35)},
        ${150 + Math.floor(Math.random() * 100)},
        ${120 + Math.floor(Math.random() * 70)},
        ${rand(66, 72)}
      );
    `;
  }

  console.log('✅ Seed complete');
}

seed().catch(err => {
  console.error('Seed error:', err);
});

