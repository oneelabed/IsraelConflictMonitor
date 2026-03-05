# RSS Aggregator
A high-performance, concurrent backend service built in Go that allows users to register, add RSS feeds, and automatically scrape posts into a searchable database. 

## 🚀 Features
- **Concurrent Scraper**: A background worker pool that uses goroutines and WaitGroups to fetch multiple RSS feeds simultaneously without blocking the main API.
- **Type-Safe SQL**: Utilizes SQLC to generate compile-time safe Go code from raw PostgreSQL queries.
- **Database Migrations**: Managed schema evolution using Goose to handle versioned SQL migrations.
- **Custom Authentication**: Secure API access via API Key-based authentication implemented through custom middleware.
- **Automated Scheduling**: Intelligent fetching logic that prioritizes feeds based on their last fetched timestamp.

## 🛠️ Tech Stack
- **Language**: Go (Golang)
- **Database**: PostgreSQL
- **Router**: Chi Router (with sub-routing and middleware)
- **SQL Toolkit**: SQLC (Type-safe query generation)
- **Migrations**: Goose
