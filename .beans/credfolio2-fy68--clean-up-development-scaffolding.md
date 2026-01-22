---
# credfolio2-fy68
title: Clean up development scaffolding
status: todo
type: task
created_at: 2026-01-22T10:24:24Z
updated_at: 2026-01-22T10:24:24Z
---

Remove temporary development artifacts created during the file upload pipeline implementation.

## Checklist

- [ ] Remove demo user seed migration (20260122021922_seed_demo_user)
- [ ] Remove /upload demo page (src/frontend/src/app/upload/page.tsx)
- [ ] Remove hardcoded DEMO_USER_ID from upload page
- [ ] Review CORS allowed origins for production (currently localhost only)
