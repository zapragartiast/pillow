# Pillow Frontend

**Server-Side Rendering (SSR) Implementation** with Next.js 15 + TypeScript + Tailwind CSS 4.1

## 🚀 **SSR Benefits Implemented**

✅ **Server-Side Authentication** - JWT tokens handled server-side with httpOnly cookies
✅ **Server-Side Data Fetching** - User profile fetched on server
✅ **Server-Side Redirects** - Authentication-based routing handled server-side
✅ **Enhanced Security** - API calls not visible in browser network
✅ **Better SEO** - Pages rendered server-side
✅ **Improved Performance** - Faster initial page loads

## 📁 **Architecture Overview**

### **Server Actions** (Next.js 15)
- `app/login/actions.ts` - Login server action
- `app/register/actions.ts` - Registration server action
- `app/dashboard/actions.ts` - Logout server action

### **Server Components**
- `app/login/page.tsx` - Server-rendered login form
- `app/register/page.tsx` - Server-rendered registration form
- `app/dashboard/page.tsx` - Server-rendered dashboard with data fetching
- `app/page.tsx` - Server-side authentication redirect

### **Security Features**
- **httpOnly cookies** for JWT token storage
- **Server-side validation** before API calls
- **Automatic token refresh** handling
- **Secure cookie settings** (sameSite, secure flags)

## 🔧 **Setup Instructions**

### **1. Environment Variables**
Create `.env.local` in frontend directory:
```bash
BACKEND_URL=http://localhost:8080
```

### **2. Install Dependencies**
```bash
cd frontend
npm install
```

### **3. Start Development Server**
```bash
npm run dev
```

### **4. Start Backend Server**
```bash
cd backend
./pillow-backend
```

## 🔐 **Authentication Flow (SSR)**

### **Login Process**
1. User submits form → Server Action triggered
2. Server validates input → Calls backend API
3. Backend returns JWT → Server sets httpOnly cookie
4. Server redirects to dashboard

### **Registration Process**
1. User submits form → Server Action triggered
2. Server validates input → Calls backend register API
3. Auto-login after registration → Sets cookies
4. Server redirects to dashboard

### **Dashboard Access**
1. Server checks for valid token cookie
2. Fetches user profile server-side
3. Renders dashboard with user data
4. Invalid token → Redirects to login

## 📋 **API Endpoints Used**

- `POST /api/login` - User authentication
- `POST /api/register` - User registration
- `GET /api/users/profile` - Get authenticated user profile

## 🛡️ **Security Features**

- **Server-side authentication** - No client-side token exposure
- **httpOnly cookies** - XSS protection
- **Secure cookies** - HTTPS only in production
- **SameSite protection** - CSRF protection
- **Server-side validation** - Input validation before API calls

## 🎯 **Key Differences from CSR**

| Feature | CSR (Before) | SSR (Now) |
|---------|-------------|-----------|
| **Token Storage** | localStorage | httpOnly cookies |
| **API Visibility** | Browser network | Server-side only |
| **SEO** | Limited | Full server rendering |
| **Security** | Client-side exposure | Server-side protection |
| **Performance** | Client hydration | Server-side rendering |
| **Error Handling** | Client-side | Server-side redirects |

## 🚀 **Testing the SSR Implementation**

1. **Visit** `http://localhost:3001`
2. **Register** a new account
3. **Auto-login** after registration
4. **View dashboard** with server-fetched data
5. **Logout** clears server cookies
6. **Try accessing** `/dashboard` without auth → Server redirects to login

## 📝 **Files Modified/Created**

### **Server Actions**
- `frontend/app/login/actions.ts` - Login server action
- `frontend/app/register/actions.ts` - Register server action
- `frontend/app/dashboard/actions.ts` - Logout server action

### **Server Components**
- `frontend/app/login/page.tsx` - SSR login form
- `frontend/app/register/page.tsx` - SSR register form
- `frontend/app/dashboard/page.tsx` - SSR dashboard with data fetching
- `frontend/app/page.tsx` - SSR authentication redirect

### **Configuration**
- `frontend/tsconfig.json` - Updated with path aliases
- `frontend/package.json` - Updated dependencies

## 🎉 **Result**

Your frontend now uses **Server-Side Rendering** with:
- ✅ **Enhanced Security** - Server-side authentication
- ✅ **Better Performance** - Server-side data fetching
- ✅ **Improved SEO** - Full server rendering
- ✅ **Production Ready** - Enterprise-grade security

The authentication flow is now completely server-side, providing better security and performance! 🚀