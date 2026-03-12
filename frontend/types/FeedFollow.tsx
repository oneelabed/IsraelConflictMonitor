import { UUID } from "crypto";
import { Timestamp } from "next/dist/server/lib/cache-handlers/types";

export interface FeedFollow {
  id: UUID;
  created_at: Timestamp;
  updated_at: Timestamp;
  user_id: UUID;
  feed_id: UUID;
}