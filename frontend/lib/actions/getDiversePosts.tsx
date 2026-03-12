const API_URL = process.env.NEXT_PUBLIC_API_URL;
import { Post } from "@/types/Post";

export async function getDiversePosts(): Promise<Post[]> {
  const res = await fetch(`${API_URL}/v1/posts/diverse`, {
    cache: 'no-store', // Use 'no-store' for real-time conflict news
  });
  
  if (!res.ok) throw new Error('Failed to fetch news');
  return res.json();
}