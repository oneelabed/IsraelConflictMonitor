import { UUID } from "crypto";
import { Timestamp } from "next/dist/server/lib/cache-handlers/types";

export interface User {
  id: UUID;
  created_at: Timestamp;
  updated_at: Timestamp;
  name: string;
  api_key: string;
}