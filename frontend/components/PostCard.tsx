import { motion } from "framer-motion";
import { ExternalLink, Clock, Newspaper, Share2 } from "lucide-react";
import { Post } from "@/types/Post";

interface PostCardProps {
  post: Post;
  index: number;
}

const PostCard = ({ post, index }: PostCardProps) => {
  const displayDate = post.published_at;
  const timeAgo = getTimeAgo(displayDate.toString());

  return (
    <motion.a
      href={post.url}
      target="_blank"
      rel="noopener noreferrer"
      className="group relative flex flex-col h-full overflow-hidden rounded-[2rem] bg-white border border-transparent shadow-sm hover:shadow-md transition-all duration-500 hover:-translate-y-1 p-8"
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ duration: 0.4, delay: Math.min(index * 0.08, 0.4) }}
    >

      {/* Top Row: Source Icon and Time */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          {/* Main Feed Icon - Circular and soft colored like the image */}
          <div className="h-10 w-10 rounded-full bg-orange-50 flex items-center justify-center border border-orange-100 group-hover:bg-orange-500 group-hover:text-white transition-colors duration-300">
            {post.feed_icon ? (
              <img 
                src={post.feed_icon} 
                alt={"icon"} 
                className="h-5 w-5 rounded-full object-contain"
              />
             
            ) : (
              <Newspaper className="h-5 w-5 text-orange-600 group-hover:text-white" />
            )}
          </div>
          <div>
            <span className="block text-sm font-bold text-[#2d241e] group-hover:text-orange-600 transition-colors capitalize">
              {post.feed_name?.replace(/_/g, ' ') || "News Feed"}
            </span>
            <span className="flex items-center gap-1 text-[11px] text-[#6b5d52]/70 font-medium">
              <Clock className="h-3.5 w-3.5" />
              {timeAgo}
            </span>
          </div>
        </div>
        
        <Share2 className="h-4 w-4 text-gray-300 hover:text-orange-500 transition-colors cursor-pointer" />
      </div>

      {/* Title - Increased size and font-weight to match the "Stories" look */}
      <h3 className="text-xl font-bold leading-tight text-[#2d241e] mb-4 line-clamp-2 group-hover:text-orange-600 transition-colors">
        {post.title}
      </h3>

      {/* Description - Softer gray and better leading (line height) */}
      {post.description && (
        <p className="text-[0.95rem] leading-relaxed text-[#6b5d52] line-clamp-3 mb-6">
          {post.description}
        </p>
      )}

      {/* Footer Link - Minimalist with an underline transition */}
      <div className="mt-auto flex items-center">
        <div className="inline-flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-orange-600/80 group-hover:text-orange-600">
          Read Article
          <ExternalLink className="h-3.5 w-3.5 transition-transform group-hover:translate-x-1 group-hover:-translate-y-1" />
        </div>
      </div>
    </motion.a>
  );
};

function getTimeAgo(dateString: string): string {
  try {
    const now = new Date();
    const date = new Date(dateString);
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);
    if (isNaN(seconds)) return "Recently";
    if (seconds < 60) return "just now";
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    if (days < 7) return `${days}d ago`;
    return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
  } catch (e) { return "Recently"; }
}

export default PostCard;