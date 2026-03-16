"use client"

import { useEffect, useState } from "react"
import { Post } from "@/types/Post"
import PostCard from "@/components/PostCard"
import { Loader2, Newspaper, Search, X } from "lucide-react"
import Link from "next/link"

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export default function UserFeedPage() {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState("")
  // We use this to track if we are currently looking at search results
  const [isSearching, setIsSearching] = useState(false)

  const fetchFeed = async (query = "") => {
    setLoading(true)
    try {
      const apiKey = localStorage.getItem("api_key")
      
      // Determine which endpoint to hit
      const endpoint = query 
        ? `${API_URL}/v1/posts/search?q=${encodeURIComponent(query)}`
        : `${API_URL}/v1/posts`

      const res = await fetch(endpoint, {
        headers: { 'Authorization': `ApiKey ${apiKey}` }
      })

      if (!res.ok) throw new Error("Fetch failed")
      const data = await res.json()
      setPosts(data)
      setIsSearching(!!query)
    } catch (err) {
      console.error("Error:", err)
    } finally {
      setLoading(false)
    }
  }

  // Effect for Debounced Search
  useEffect(() => {
    const timer = setTimeout(() => {
      fetchFeed(searchQuery)
    }, 400) // Wait 400ms after user stops typing

    return () => clearTimeout(timer)
  }, [searchQuery])

  return (
    <div className="min-h-screen bg-slate-50/50 p-6 lg:p-10">
      <header className="max-w-7xl mx-auto mb-10 space-y-6">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-extrabold tracking-tight lg:text-4xl">
              {isSearching ? "Search Results" : "My Personal Monitor"}
            </h1>
            <p className="text-muted-foreground mt-1">
              {isSearching 
                ? `Showing results for "${searchQuery}"` 
                : "Updates from your followed sources."}
            </p>
          </div>

          {/* Search Bar UI */}
          <div className="relative w-full md:w-80">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
            <input
              type="text"
              placeholder="Search your feed..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-10 py-2 bg-white border border-slate-200 rounded-full focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all shadow-sm"
            />
            {searchQuery && (
              <button 
                onClick={() => setSearchQuery("")}
                className="absolute right-3 top-1/2 -translate-y-1/2 hover:text-red-500 transition-colors"
              >
                <X className="h-4 w-4" />
              </button>
            )}
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto">
        {loading ? (
          <div className="flex justify-center py-20">
            <Loader2 className="animate-spin h-10 w-10 text-primary" />
          </div>
        ) : posts.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {posts.map((post, index) => (
              <PostCard key={post.id || index} post={post} index={index}/>
            ))}
          </div>
        ) : (
          <div className="text-center py-24 bg-white rounded-3xl border-2 border-dashed border-slate-200">
            <div className="bg-slate-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <Search className="h-8 w-8 text-slate-400" />
            </div>
            <h3 className="text-xl font-semibold text-slate-900">
              {isSearching ? "No matches found" : "Your feed is empty"}
            </h3>
            <p className="text-muted-foreground mt-2 max-w-xs mx-auto">
              {isSearching 
                ? "Try searching for something else, like 'Hezbollah' or 'Gaza'." 
                : "Follow some news outlets in the Discover section to get started."}
            </p>
            {isSearching && (
              <button 
                onClick={() => setSearchQuery("")}
                className="mt-6 text-primary font-semibold hover:underline"
              >
                Clear Search
              </button>
            )}
          </div>
        )}
      </main>
    </div>
  )
}