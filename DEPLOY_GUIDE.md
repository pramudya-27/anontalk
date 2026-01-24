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
- Your Database Credentials (Aiven):
  - **Host**: `anontalk-basisdata-1.d.aivencloud.com`
  - **Port**: `17892`
  - **User**: `avnadmin`
  - **Password**: `(See your Aiven Dashboard)`
  - **DB Name**: `anontalk`

> [!IMPORTANT]
> **Initialize Database**: Your new Aiven database is likely empty. You MUST run the `migrations/001_init.sql` script to create the necessary tables before deploying.
>
> Run this command in your terminal (if you have mysql installed):
>
> ```bash
> mysql --user avnadmin --password=YOUR_PASSWORD --host anontalk-basisdata-1.d.aivencloud.com --port 17892 anontalk < backend/migrations/001_init.sql
> ```
>
> Or use a database tool (like DBeaver) to connect and run the SQL file content manually.

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
    - `DATABASE_URL`: `avnadmin:YOUR_PASSWORD@tcp(anontalk-basisdata-1.d.aivencloud.com:17892)/anontalk?parseTime=true&tls=skip-verify`
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

## 5. Deploying Frontend to Vercel

1.  **Push Changes**:
    Ensure you have committed and pushed the latest changes (fixing the WebSocket connection):

    ```bash
    git push
    ```

2.  **Go to Vercel**:
    - Log in to [vercel.com](https://vercel.com).
    - Click **"Add New..."** -> **"Project"**.
    - Import your `anontalk` repository.

3.  **Configure Project**:
    - **Framework Preset**: Next.js (should be auto-detected).
    - **Root Directory**: `./` (default).

4.  **Environment Variables**:
    Expand the "Environment Variables" section and add:
    - **Key**: `NEXT_PUBLIC_API_URL`
    - **Value**: `https://anontalk-production.up.railway.app` (Your active Railway URL)

5.  **Deploy**:
    Click **"Deploy"**. Vercel will build your frontend and publish it.

6.  **Verify**:
    Open your new Vercel URL. Try to:
    - Signup/Login.
    - Start a chat.
    - Confirm messages send and receive.

---

## Summary of Architecture

- **Frontend**: Hosted on **Vercel** (Next.js).
- **Backend (API + Redis)**: Hosted on **Railway/Render** (Docker Container).
- **Database**: Hosted on **Aiven** (MySQL).
