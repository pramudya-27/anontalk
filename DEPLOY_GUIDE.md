# Deployment Guide: Backend + Redis + MySQL (Separated)

This guide covers deploying the **Go Backend**, **Redis**, and **MySQL** as separate services.

## Architecture

1.  **MySQL**: Managed Cloud DB (Aiven).
2.  **Redis**: Separate Docker Service (deployed from `redis` branch).
3.  **Backend**: Go App Service (deployed from `backend` branch).

---

## 1. Prerequisites (Database)

Ensure your Aiven MySQL is ready and initialized:

- **Host**: `anontalk-basisdata-1.d.aivencloud.com`
- **Port**: `17892`
- **User**: `avnadmin`
- **Password**: `(See Aiven Dashboard)`
- **DB Name**: `anontalk`

> [!IMPORTANT]
> **Initialize Schema**: If not done yet, run:
>
> ```bash
> mysql --user avnadmin --password=YOUR_PASSWORD --host anontalk-basisdata-1.d.aivencloud.com --port 17892 anontalk < backend/migrations/001_init.sql
> ```

---

## 2. Deploy Redis (Service 1)

1.  **Switch to Redis Branch**:
    ```bash
    git checkout redis
    git push origin redis
    ```
2.  **Railway**:
    - "New Project" -> "Deploy from GitHub".
    - Select Repo.
    - **Variables**: None needed (Defaults to port 6379).
    - **Settings**:
      - **Source Directory/Context**: `/` or default.
      - **Branch**: Select `redis` (Important!).

3.  **Get Internal URL**:
    - Once deployed, go to the Redis service settings/variables.
    - Find the **Internal Host** (e.g., `redis-production.up.railway.app`) and **Port** (e.g., `6379`).
    - Note this as `REDIS_ADDR`.
    - _Note: On Railway, services in the same project can talk via internal DNS names. If you use a private network, use the private domain._

---

## 3. Deploy Backend (Service 2)

1.  **Switch to Backend Branch**:

    ```bash
    git checkout backend
    git push -f origin backend
    ```

    _(Use -f if you rewrote history to remove secrets previously)_.

2.  **Railway**:
    - Add a **New Service** to the _same_ project (or create new).
    - "Deploy from GitHub".
    - Select Repo -> Branch `backend`.

3.  **Variables**: Add these variables to the Backend Service:
    - `DATABASE_URL`: `avnadmin:YOUR_PASSWORD@tcp(anontalk-basisdata-1.d.aivencloud.com:17892)/anontalk?parseTime=true&tls=skip-verify`
    - `REDIS_ADDR`: `redis-host:6379` (Replace with the actual Redis address from Step 2).
    - `JWT_SECRET`: (Random Secret).
    - `PORT`: `8080`.

4.  **Build Settings**:
    - **Dockerfile Path**: `/Dockerfile` (Root).
    - **Context**: `/` (Root).

5.  **Deploy**:
    - Railway will build the Go app and connect to Redis and MySQL.

---

## 4. Frontend (Vercel)

No changes needed if already deployed. Just ensure `NEXT_PUBLIC_API_URL` points to the new **Backend Service Public URL**.
