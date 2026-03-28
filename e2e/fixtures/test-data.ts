import { API_BASE, TEST_TRACKS } from './test-constants';

export interface Listen {
  artist_name: string;
  track_name: string;
  release_name?: string;
  listened_at?: number;
}

export interface ApiKeyResponse {
  id: number;
  key: string;
  label: string;
  user_id: number;
  created_at: string;
}

export async function getApiKey(sessionCookie: string, label: string): Promise<string> {
  const response = await fetch(`${API_BASE}/apis/web/v1/user/apikeys`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      'Cookie': sessionCookie,
    },
    body: new URLSearchParams({ label }),
  });

  if (!response.ok) {
    throw new Error(`Failed to create API key: ${response.status} ${response.statusText}`);
  }

  const data: ApiKeyResponse = await response.json();
  return data.key;
}

export async function createTestListens(apiKey: string, listens: Listen[]): Promise<void> {
  const payload = {
    listen_type: 'import',
    payload: listens.map(listen => ({
      listened_at: listen.listened_at || Math.floor(Date.now() / 1000),
      track_metadata: {
        artist_name: listen.artist_name,
        track_name: listen.track_name,
        release_name: listen.release_name,
      },
    })),
  };

  const response = await fetch(`${API_BASE}/apis/listenbrainz/1/submit-listens`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Token ${apiKey}`,
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Error(`Failed to submit listens: ${response.status} ${response.statusText}`);
  }
}

export function generateRandomListens(count: number, daysBack: number = 7): Listen[] {
  const listens: Listen[] = [];
  const now = Date.now();
  const msPerDay = 24 * 60 * 60 * 1000;

  for (let i = 0; i < count; i++) {
    const track = TEST_TRACKS[Math.floor(Math.random() * TEST_TRACKS.length)];
    const randomDay = Math.floor(Math.random() * daysBack);
    const randomHour = Math.floor(Math.random() * 24);
    const randomMinute = Math.floor(Math.random() * 60);

    const timestamp = now - (randomDay * msPerDay) - (randomHour * 60 * 60 * 1000) - (randomMinute * 60 * 1000);

    listens.push({
      artist_name: track.artist,
      track_name: track.name,
      release_name: track.album,
      listened_at: Math.floor(timestamp / 1000),
    });
  }

  return listens.sort((a, b) => (a.listened_at || 0) - (b.listened_at || 0));
}

export async function seedDefaultData(sessionCookie: string, count: number = 20): Promise<string> {
  const apiKey = await getApiKey(sessionCookie, 'e2e-test-key');
  const listens = generateRandomListens(count);
  await createTestListens(apiKey, listens);
  return apiKey;
}
