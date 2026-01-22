---
# credfolio2-hesu
title: Add CORS middleware to backend
status: completed
type: bug
priority: normal
created_at: 2026-01-22T10:16:36Z
updated_at: 2026-01-22T10:17:11Z
parent: credfolio2-k38n
---

File uploads fail with 'Network error during upload' because the backend lacks CORS headers. Browser blocks cross-origin requests from frontend (port 3000) to backend (port 8080).