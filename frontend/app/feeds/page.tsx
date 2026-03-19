"use client"

import { useEffect, useState } from "react"
import { Check, Plus, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import Image from "next/image"
import { Feed } from "@/types/Feed"
import { useRouter } from "next/navigation";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export default function DiscoverFeeds() {
  const [feeds, setFeeds] = useState<Feed[]>([])
  const [loading, setLoading] = useState(true)
  const router = useRouter();

  useEffect(() => {
    const apiKey = localStorage.getItem("api_key");
    
    if (!apiKey) {
      router.push("/login");
    }
  }, [router]);

  useEffect(() => {
    // 1. Fetch all feeds from your Go Backend
    // You should send your API Key in the headers for the backend to check follows
    const fetchFeeds = async () => {
      try {
        const res = await fetch(`${API_URL}/v1/feeds`, {
            headers: { 'Authorization': `ApiKey ${localStorage.getItem("api_key")}` }
        })
        const data = await res.json()
        setFeeds(data)
      } catch (err) {
        console.error("Failed to load feeds", err)
      } finally {
        setLoading(false)
      }
    }
    fetchFeeds()
  }, [])

  const toggleFollow = async (feedId: string, currentlyFollowing: boolean) => {
    const method = currentlyFollowing ? 'DELETE' : 'POST'
    const endpoint = `/v1/feed_follows`
    const apiKey = localStorage.getItem("api_key")

    try {
      // Optimistic Update: Change UI immediately
      setFeeds(current => current.map(f => 
        f.id === feedId ? { ...f, is_following: !currentlyFollowing } : f
      ))

      await fetch(`${API_URL}${endpoint}`, {
        method,
        headers: { 
            'Content-Type': 'application/json',
            'Authorization': `ApiKey ${apiKey}` 
        },
        body: JSON.stringify({ feed_id: feedId })
      })
    } catch (err) {
      // Revert if it fails
      console.error("Failed to update follow", err)
    }
  }

  if (loading) return <div className="flex justify-center p-20"><Loader2 className="animate-spin" /></div>

  return (
    <div className="max-w-4xl mx-auto p-6 lg:p-10">
      <header className="mb-10 text-center">
        <h1 className="text-4xl font-extrabold tracking-tight mb-3">Discover Sources</h1>
        <p className="text-muted-foreground text-lg">
          Follow the news outlets you trust to customize your Conflict Monitor feed.
        </p>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {feeds.map((feed) => (
          <Card key={feed.id} className="overflow-hidden border-2 transition-all hover:border-primary/20">
            <CardHeader className="flex flex-row items-center gap-4 space-y-0">
              <div className="relative w-12 h-12 flex-shrink-0">
                <div className="w-full h-full rounded-full overflow-hidden border bg-white flex items-center justify-center shadow-sm">
                  <Image 
                    src={feed.icon_url} 
                    alt={feed.name}
                    width={48} 
                    height={48}
                    className="object-contain p-1" 
                  />
                </div>
              </div>
              <div className="flex-1">
                <CardTitle className="text-xl">{feed.name}</CardTitle>
                <CardDescription className="truncate">{feed.url}</CardDescription>
              </div>
            </CardHeader>
            <CardContent className="flex justify-end pt-0 pb-6 pr-6">
              <Button 
                variant={feed.is_following ? "outline" : "default"}
                onClick={() => toggleFollow(feed.id, feed.is_following)}
                className="w-32 rounded-full font-semibold transition-all"
              >
                {feed.is_following ? (
                  <>
                    <Check className="mr-2 h-4 w-4" /> Following
                  </>
                ) : (
                  <>
                    <Plus className="mr-2 h-4 w-4" /> Follow
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}