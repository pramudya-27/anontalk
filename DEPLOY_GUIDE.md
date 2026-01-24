# Deployment Guide: Go Backend + Redis (Combined Image)

> [!WARNING]
> **Vercel Incompatibility**: You specifically requested to "make a combined Docker image" for the backend. **Vercel does NOT support deploying custom Docker images.** Vercel is for frontend and serverless functions only.
>
> To deploy this **Docker image**, you must use a container hosting platform like **Railway**, **Render**, **Fly.io**, or a **VPS** (DigitalOcean, etc.).
>
> This guide provides instructions for **Railway** and **Render** as they are the easiest alternatives to Vercel for Docker hosting.

## 1. Prerequisites

- [Docker](https://www.docker.com/) installed locally (for testing).
- A GitHub repository containing this project.
- Your Database Credentials:
  - **Host**: `sql101.byethost3.com`
  - **Port**: `3306`
  - **User**: `b3_40872396`
  - **Password**: `daimyo266713`
  - **DB Name**: `b3_40872396_anontalk`

---

## 2. Prepare the Repository

Ensure you have pushed the new files to your GitHub repository:

- `backend/Dockerfile.combined`
- `backend/entrypoint.sh`

```bash
git add backend/Dockerfile.combined backend/entrypoint.sh
git commit -m "Add combined Dockerfile and entrypoint"
git push origin main
```

---

## 3. Deployment Option A: Railway (Recommended)

Railway matches your "Zero Config" desire similar to Vercel but supports Docker natively.

1.  **Sign Up/Login**: Go to [railway.app](https://railway.app/) and login with GitHub.
2.  **New Project**: Click "New Project" -> "Deploy from GitHub repo".
3.  **Select Repo**: Choose your `anontalk` repository.
4.  **Configuration**:
    - Railway usually auto-detects Dockerfiles.
    - Go to **Settings** -> **Build**.
    - **Dockerfile Path**: Set this to `/Dockerfile.combined` (It's now in the root!).
    - **Context**: Set to `/` (Root directory).
5.  **Variables**: Go to the **Variables** tab and add your Database details:
    - `DATABASE_URL`: `b3_40872396:daimyo266713@tcp(sql101.byethost3.com:3306)/b3_40872396_anontalk`
    - `JWT_SECRET`: (Create a secure random string)
    - `SERVER_PORT`: `8080`
6.  **Deploy**: Railway will build and deploy.
7.  **Public URL**: Go to **Settings** -> **Networking** and "Generate Domain" to get your public API URL (e.g., `https://anontalk-production.up.railway.app`).

---

## 4. Deployment Option B: Render

1.  **Sign Up/Login**: Go to [render.com](https://render.com/).
2.  **New Web Service**: Click "New +" -> "Web Service".
3.  **Connect Repo**: Connect your GitHub account and select `anontalk`.
4.  **Settings**:
    - **Runtime**: Docker
    - **Context Directory**: `backend` (This is crucial, everything is relative to this).
    - **Dockerfile Path**: `Dockerfile.combined` (It's inside the context directory).
5.  **Environment Variables**:
    - Add `DATABASE_URL`, `JWT_SECRET`, `SERVER_PORT` (Value: `8080`).
6.  **Deploy**: Click "Create Web Service".

---

## 5. Connecting Frontend to Backend

Once deployed, copy the **Public URL** from Railway/Render.

1.  Open your local frontend code.
2.  Update your `lib/api-client.ts` (or equivalent config).
3.  Replace `http://localhost:8080` with your new production URL (e.g., `https://anontalk-production.up.railway.app`).
4.  Commit and push the frontend changes to deploy the frontend (which _can_ stay on Vercel).

---

## Summary of Architecture

- **Frontend**: Hosted on **Vercel** (Next.js).
- **Backend (API + Redis)**: Hosted on **Railway/Render** (Docker Container).
- **Database**: Hosted on **ByetHost** (MySQL).
