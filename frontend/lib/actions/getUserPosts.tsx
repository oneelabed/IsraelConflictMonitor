const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function getUserPosts(apiKey: string) {
  const res = await fetch(`${API_URL}/v1/posts`, {
    headers: {
      Authorization: `ApiKey ${apiKey}`,
    },
  });

  if (!res.ok) throw new Error("Failed to fetch posts");

  return res.json();
}