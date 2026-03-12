import { UUID } from "crypto";
import { Timestamp } from "next/dist/server/lib/cache-handlers/types";

export interface Post {
  id: UUID;
  created_at: Timestamp;
  updated_at: Timestamp;
  title: string;
  description: string;
  published_at: Timestamp;
  url: string;
  feed_icon: string;
  feed_name: string;
  feed_id: UUID;
}

