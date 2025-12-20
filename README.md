# Postastiq

**Postastiq** is a self-hosted, privacy-first blogging platform written in Go.

It’s designed for individuals and small teams who want a **clean, focused publishing experience**, **full ownership of their data**, and **no external service dependencies**.

Postastiq runs as a single binary, stores data locally using SQLite, and supports rich media, custom domains, RSS feeds, and optional access controls — all without requiring a separate database, background workers, or third-party services.

---

## Why Postastiq?

- **Self-hosted by default**  
  Your content and data stay under your control.

- **Minimal by design**  
  Fast, focused, and distraction-free — for both writers and readers.

- **Single binary architecture**  
  No MySQL, Redis, or external infrastructure required.

- **Privacy-aware**  
  Optional password protection for public blogs.

---

## License & Project Model

Postastiq is **source-available software** released under the  
**Business Source License (BSL) 1.1**.

The source code is publicly available for transparency and self-hosting.
You are free to run Postastiq for personal or internal use.

However:
- Postastiq may not be offered as a hosted or managed service
- The software may not be used to build or operate competing platforms
- The Postastiq name and branding may not be reused or rebranded

Postastiq is developed and maintained exclusively by the Postastiq team.
External pull requests are not accepted at this time, though issues and
feedback are welcome.

On the license Change Date, the project will automatically convert to
the MIT License. See `LICENSE.md` for full terms.

---

## Quick Start

### Docker (recommended)

```bash
docker run -dit -p 8080:8080 postastiq/postastiq
```

### Local binary

```bash
go build -o postastiq .
./postastiq
```

Once running:

- **Blog:** http://localhost:8080/
- **Admin:** http://localhost:8080/admin  
  **Default credentials:**  
  - Username: `admin`  
  - Password: `admin`

> You will be prompted to change the admin password on first login.

---

## Core Features

- **Public Blog Viewer (`/`)**  
  Clean, minimal feed with infinite scroll.

- **Admin Dashboard (`/admin`)**  
  Create, edit, and manage posts, media, and site settings.

- **Rich Media Support**  
  Upload photos, audio, and video with optional thumbnails.

- **Privacy Controls**  
  Optional password protection for public access.

- **Themes**  
  Built-in light and dark modes.

- **Custom Domains**  
  Bring your own domain with automatic HTTPS (via Caddy).

- **Backups & Restore**  
  One-click ZIP export and restore (database + media).

- **RSS Feed**  
  Automatically generated RSS 2.0 feed at `/rss`.

- **SEO-Friendly URLs**  
  Human-readable post URLs, for example:  
  `/posts/my-title-2025-12-15/`

- **Embedded SQLite**  
  No external database required.

---

## Routes Overview

### Public Routes

| Method | Path | Description |
|------|------|------------|
| GET | `/` | Main blog feed |
| GET | `/posts/:slug/` | Individual post |
| GET | `/api/entries` | JSON API |
| GET | `/rss` | RSS feed |
| GET | `/uploads/:filename` | Media files |

---

### Admin Routes

| Method | Path | Description |
|------|------|------------|
| GET | `/admin` | Admin dashboard |
| POST | `/admin/create` | Create post |
| POST | `/admin/update` | Update post |
| POST | `/admin/delete` | Delete post |
| GET | `/admin/settings` | Site settings |
| GET | `/admin/backup` | Download backup |
| POST | `/admin/restore` | Restore backup |

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|--------|---------|------------|
| `PORT` | `8080` | HTTP port |
| `DB_PATH` | `/app/data/blog.db` | SQLite database location |
| `UPLOADS_DIR` | `/app/data/uploads` | Media storage directory |
| `ADMIN_PASSWORD` | `admin` | Initial admin password |

> For security, change the admin password immediately after first login.

---

## Media Support

| Type | Extensions | Max Size |
|-----|-----------|---------|
| Photos | jpg, jpeg, png, gif, webp | 5 MB |
| Audio | mp3, m4a, wav, ogg, aac, webm | 20 MB |
| Video | mp4, webm, mov, avi | 50 MB |

---

## Health Check

`GET /health`

```json
{
  "status": "ready",
  "bootstrap": true
}
```

---

## License

Business Source License 1.1

See `LICENSE.md` for full terms.
