const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function createUser(username: string, password: string) {
  const res = await fetch(`${API_URL}/v1/users`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      username: username,
      password_hash: password,
    }),
  });
  
  if (res.status === 409) {
    throw new Error('Username already taken');
  }

  if (!res.ok) throw new Error('Failed to create user');
  return res.json();
}