"use client"

import { useEffect, useState } from "react"
import { getDiversePosts } from "@/lib/actions/getDiversePosts"
import { Post } from "@/types/Post"
import PostCard from "@/components/PostCard"

export default function HomePage() {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getDiversePosts()
      .then((data) => {
        setPosts(data)
        setLoading(false)
      })
      .catch((err) => {
        console.error("Failed to fetch posts:", err)
        setLoading(false)
      })
  }, [])

  return (
    <div className="min-h-screen bg-background p-6 lg:p-10">
      <header className="max-w-7xl mx-auto mb-10">
        <h1 className="text-3xl font-extrabold tracking-tight lg:text-4xl mb-2 text-center">
          Latest Conflict Updates
        </h1>
      </header>

      {loading ? (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : (
        <div className="max-w-7xl mx-auto">
          {/* Responsive Grid: 1 column on mobile, 2 on tablet, 3 on desktop */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {posts.map((post, index) => (
              <PostCard 
                key={post.id || index} 
                post={post} 
                index={index} 
              />
            ))}
          </div>

          {posts.length === 0 && (
            <div className="text-center py-20">
              <p className="text-muted-foreground text-lg">No posts found at the moment.</p>
            </div>
          )}
        </div>
      )}
    </div>
  )
}