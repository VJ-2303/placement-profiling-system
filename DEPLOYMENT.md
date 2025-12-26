# 🚀 Complete Deployment Guide

## KCT Placement Profiling System

Deploy your application with:
- **Frontend**: Netlify (Free)
- **Backend**: Railway ($5 free credit/month)
- **Database**: Neon PostgreSQL (Free tier)
- **Auth**: Microsoft Azure AD

---

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Step 1: Neon PostgreSQL Setup](#step-1-neon-postgresql-setup)
3. [Step 2: Microsoft Azure AD Setup](#step-2-microsoft-azure-ad-setup)
4. [Step 3: Railway Backend Deployment](#step-3-railway-backend-deployment)
5. [Step 4: Netlify Frontend Deployment](#step-4-netlify-frontend-deployment)
6. [Step 5: Final Configuration](#step-5-final-configuration)
7. [Step 6: Add Admin Users](#step-6-add-admin-users)
8. [Testing & Troubleshooting](#testing--troubleshooting)

---

## Prerequisites

Before starting, ensure you have:
- [ ] GitHub account with this repository pushed
- [ ] Email account for creating service accounts
- [ ] Access to your college's Microsoft Azure AD (for OAuth)

---

## Step 1: Neon PostgreSQL Setup

### 1.1 Create Neon Account

1. Go to [https://neon.tech](https://neon.tech)
2. Click **Sign Up** → Sign in with GitHub
3. Verify your email if required

### 1.2 Create Database Project

1. Click **Create Project**
2. Configure:
   - **Project name**: `kct-placement`
   - **Region**: Choose closest to your users (e.g., `Asia Pacific (Singapore)`)
   - **PostgreSQL version**: `16` (latest)
3. Click **Create Project**

### 1.3 Get Connection String

1. After project creation, you'll see the **Connection Details**
2. Click **Copy** on the connection string
3. It looks like:
   ```
   postgresql://username:password@ep-cool-name-123456.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
   ```
4. **Save this securely** - you'll need it for Railway

### 1.4 Run Database Migrations

1. In Neon Dashboard, click **SQL Editor** (left sidebar)
2. Copy the entire contents of `backend/migrations/001_initial_schema.sql`
3. Paste into the SQL Editor
4. Click **Run** (or press Ctrl+Enter)
5. You should see "Query executed successfully"

### 1.5 Verify Tables Created

Run this query in SQL Editor:
```sql
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
```

You should see tables like: `students`, `admins`, `companies`, `placements`, etc.

---

## Step 2: Microsoft Azure AD Setup

### 2.1 Access Azure Portal

1. Go to [https://portal.azure.com](https://portal.azure.com)
2. Sign in with your **college Microsoft account** (e.g., `admin@kct.ac.in`)
3. If you don't have admin access, contact your IT department

### 2.2 Navigate to App Registrations

1. In the search bar, type **"App registrations"**
2. Click on **App registrations** under Services
3. Click **+ New registration**

### 2.3 Register the Application

Fill in the form:

| Field | Value |
|-------|-------|
| **Name** | `KCT Placement Portal` |
| **Supported account types** | `Accounts in this organizational directory only (Single tenant)` |
| **Redirect URI - Platform** | `Web` |
| **Redirect URI - URL** | `http://localhost:4000/auth/callback` (we'll update this later) |

Click **Register**

### 2.4 Note Important IDs

After registration, you'll see the **Overview** page. Copy these values:

| Value | Where to Find | Example |
|-------|---------------|---------|
| **Application (client) ID** | Overview page | `12345678-abcd-1234-efgh-123456789abc` |
| **Directory (tenant) ID** | Overview page | `87654321-dcba-4321-hgfe-987654321cba` |

### 2.5 Create Client Secret

1. In left sidebar, click **Certificates & secrets**
2. Click **+ New client secret**
3. Configure:
   - **Description**: `Production Secret`
   - **Expires**: `24 months` (recommended)
4. Click **Add**
5. **IMPORTANT**: Copy the **Value** immediately (it won't be shown again!)
   - Example: `abc123~XYZ789secretKeyHere`

### 2.6 Configure API Permissions

1. In left sidebar, click **API permissions**
2. Click **+ Add a permission**
3. Select **Microsoft Graph**
4. Select **Delegated permissions**
5. Search and check these permissions:
   - [x] `email`
   - [x] `openid`
   - [x] `profile`
   - [x] `User.Read`
6. Click **Add permissions**
7. Click **Grant admin consent for [Your Organization]**
8. Confirm by clicking **Yes**

### 2.7 Configure Token Settings (Optional but Recommended)

1. In left sidebar, click **Token configuration**
2. Click **+ Add optional claim**
3. Select **ID** token type
4. Check:
   - [x] `email`
   - [x] `preferred_username`
5. Click **Add**

### 2.8 Summary of Azure AD Values

Save these values securely:

```
MICROSOFT_CLIENT_ID=<Application (client) ID>
MICROSOFT_CLIENT_SECRET=<Client secret Value>
MICROSOFT_TENANT_ID=<Directory (tenant) ID>
```

---

## Step 3: Railway Backend Deployment

### 3.1 Create Railway Account

1. Go to [https://railway.app](https://railway.app)
2. Click **Login** → **Login with GitHub**
3. Authorize Railway to access your GitHub

### 3.2 Create New Project

1. Click **+ New Project**
2. Select **Deploy from GitHub repo**
3. If prompted, click **Configure GitHub App** and give access to your repository
4. Select `placement-profiling-system` repository

### 3.3 Configure Service

1. Railway will detect the Dockerfile automatically
2. Click on the service card
3. Go to **Settings** tab
4. Set:
   - **Root Directory**: `backend`
   - **Watch Paths**: `/backend/**`

### 3.4 Configure Environment Variables

1. Click on **Variables** tab
2. Click **+ New Variable** for each:

| Variable | Value |
|----------|-------|
| `DATABASE_URL` | Your Neon connection string |
| `MICROSOFT_CLIENT_ID` | From Azure AD |
| `MICROSOFT_CLIENT_SECRET` | From Azure AD |
| `MICROSOFT_TENANT_ID` | From Azure AD |
| `MICROSOFT_REDIRECT_URL` | `https://YOUR-APP.railway.app/auth/callback` (update after getting domain) |
| `JWT_SECRET` | Generate with: `openssl rand -hex 32` |
| `FRONTEND_URL` | `https://YOUR-SITE.netlify.app` (update after Netlify deploy) |
| `ENV` | `production` |
| `PORT` | `4000` |

### 3.5 Generate Railway Domain

1. Go to **Settings** tab
2. Scroll to **Networking**
3. Click **Generate Domain**
4. You'll get a URL like: `https://placement-api-production.up.railway.app`
5. **Copy this URL**

### 3.6 Update Microsoft Redirect URL

1. Go back to **Variables** tab in Railway
2. Update `MICROSOFT_REDIRECT_URL` to: `https://YOUR-RAILWAY-DOMAIN/auth/callback`
3. Also update in **Azure Portal**:
   - Go to App registrations → Your app → Authentication
   - Update Redirect URI to match

### 3.7 Deploy

1. Railway auto-deploys when you push to GitHub
2. Or click **Deploy** button manually
3. Wait for build to complete (2-3 minutes)
4. Check **Deployments** tab for logs

### 3.8 Verify Backend

Visit: `https://YOUR-RAILWAY-DOMAIN/health`

You should see:
```json
{
  "status": "healthy",
  "service": "placement-api"
}
```

---

## Step 4: Netlify Frontend Deployment

### 4.1 Update Frontend Configuration

Before deploying, update the API URL:

**Edit `frontend/js/config.js`:**
```javascript
const API_BASE_URL = 'https://YOUR-RAILWAY-DOMAIN';  // e.g., https://placement-api-production.up.railway.app
```

Commit and push this change to GitHub.

### 4.2 Create Netlify Account

1. Go to [https://netlify.com](https://netlify.com)
2. Click **Sign up** → **Sign up with GitHub**
3. Authorize Netlify

### 4.3 Deploy from GitHub

1. Click **Add new site** → **Import an existing project**
2. Select **GitHub**
3. Choose your `placement-profiling-system` repository
4. Configure build settings:

| Setting | Value |
|---------|-------|
| **Base directory** | `frontend` |
| **Build command** | (leave empty) |
| **Publish directory** | `frontend` |

5. Click **Deploy site**

### 4.4 Get Netlify URL

1. After deployment, you'll get a URL like: `https://random-name-123.netlify.app`
2. **Copy this URL**

### 4.5 (Optional) Custom Domain

1. Go to **Domain settings**
2. Click **Add custom domain**
3. Follow instructions to add your domain

### 4.6 Update Railway with Netlify URL

Go back to Railway and update these environment variables:

| Variable | Value |
|----------|-------|
| `FRONTEND_URL` | `https://YOUR-SITE.netlify.app` |

---

## Step 5: Final Configuration

### 5.1 Update Azure AD Redirect URI

1. Go to [Azure Portal](https://portal.azure.com)
2. Navigate to **App registrations** → **KCT Placement Portal**
3. Click **Authentication** in left sidebar
4. Under **Web** → **Redirect URIs**, ensure you have:
   - `https://YOUR-RAILWAY-DOMAIN/auth/callback`
5. Click **Save**

### 5.2 Verify All URLs Match

Create a checklist:

| Location | URL Should Be |
|----------|---------------|
| `frontend/js/config.js` → `API_BASE_URL` | Railway backend URL |
| Railway → `FRONTEND_URL` | Netlify frontend URL |
| Railway → `MICROSOFT_REDIRECT_URL` | `{Railway URL}/auth/callback` |
| Azure AD → Redirect URI | `{Railway URL}/auth/callback` |

### 5.3 Redeploy if Needed

- **Frontend changes**: Push to GitHub, Netlify auto-deploys
- **Backend changes**: Push to GitHub, Railway auto-deploys
- **Environment variable changes**: Railway auto-redeploys

---

## Step 6: Add Admin Users

### 6.1 Add Admins to Database

1. Go to [Neon Dashboard](https://console.neon.tech)
2. Open **SQL Editor**
3. Run this query (replace with actual admin emails):

```sql
INSERT INTO admins (email, name, department, is_active) VALUES 
('placement.officer@kct.ac.in', 'Placement Officer', 'Training & Placement', true),
('hod.cse@kct.ac.in', 'Dr. HOD Name', 'Computer Science', true),
('admin@kct.ac.in', 'Admin User', 'Administration', true);
```

### 6.2 Verify Admins Added

```sql
SELECT * FROM admins;
```

---

## Testing & Troubleshooting

### Test the Complete Flow

1. **Health Check**
   - Visit: `https://YOUR-RAILWAY-DOMAIN/health`
   - Should return: `{"status":"healthy"}`

2. **Frontend Load**
   - Visit: `https://YOUR-NETLIFY-SITE.netlify.app`
   - Should see login page

3. **Student Login**
   - Click "Login with Microsoft"
   - Sign in with a `@kct.ac.in` student email
   - Should redirect to profile page

4. **Admin Login**
   - Click "Login as Admin"
   - Sign in with an admin email (added to database)
   - Should redirect to admin dashboard

### Common Issues & Fixes

#### CORS Errors
```
Access to fetch at 'https://api...' from origin 'https://frontend...' has been blocked by CORS
```
**Fix**: Ensure `FRONTEND_URL` in Railway exactly matches your Netlify URL (including `https://`)

#### OAuth Redirect Mismatch
```
AADSTS50011: The redirect URI specified in the request does not match
```
**Fix**: 
1. Check Railway `MICROSOFT_REDIRECT_URL` matches Azure AD redirect URI exactly
2. No trailing slashes
3. Must be `https://`

#### Database Connection Failed
```
failed to connect to host=... connection refused
```
**Fix**:
1. Check `DATABASE_URL` includes `?sslmode=require`
2. Verify Neon project is active (not paused)
3. Check connection string is correct

#### 401 Unauthorized
```
{"error":"unauthorized"}
```
**Fix**:
1. Token might be expired - try logging in again
2. Check `JWT_SECRET` is set in Railway
3. Verify token is being sent in Authorization header

#### Domain Not Allowed (OAuth)
```
You're not allowed to sign in from this domain
```
**Fix**: Ensure the user's email domain is `@kct.ac.in` (or your configured domain)

### View Logs

**Railway Logs:**
1. Go to your Railway project
2. Click on the service
3. Click **Deployments** → Click latest deployment → **View Logs**

**Netlify Logs:**
1. Go to your Netlify site
2. Click **Deploys** → Click latest deploy → **Deploy log**

---

## 🎉 Deployment Complete!

Your application is now live:

| Service | URL |
|---------|-----|
| **Frontend** | `https://YOUR-SITE.netlify.app` |
| **Backend API** | `https://YOUR-APP.railway.app` |
| **Health Check** | `https://YOUR-APP.railway.app/health` |
| **Database** | Neon Dashboard |

### Quick Links

- [Neon Dashboard](https://console.neon.tech)
- [Railway Dashboard](https://railway.app/dashboard)
- [Netlify Dashboard](https://app.netlify.com)
- [Azure Portal](https://portal.azure.com)

---

## 📞 Support

If you encounter issues:
1. Check the troubleshooting section above
2. Review logs in Railway/Netlify
3. Verify all environment variables are set correctly
4. Ensure database migrations ran successfully
