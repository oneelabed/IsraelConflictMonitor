import { UUID } from "crypto";
import { Timestamp } from "next/dist/server/lib/cache-handlers/types";

export interface Feed {
  id: UUID;
  created_at: Timestamp;
  updated_at: Timestamp;
  name: string;
  url: string;
  icon_url: string;
  user_id: UUID;
}