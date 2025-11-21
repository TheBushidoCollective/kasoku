# Kasoku Web Dashboard

Next.js web dashboard for Kasoku cache server.

## Features

- User authentication (sign up, login, logout)
- Analytics dashboard with cache hit/miss rates and time saved
- Cache management (view and delete cache entries)
- API token management for CLI/CI-CD integration
- Responsive design with Tailwind CSS

## Setup

1. **Install dependencies**:

```bash
npm install
```

2. **Configure environment**:

```bash
cp .env.example .env
```

Edit `.env` and set `NEXT_PUBLIC_API_URL` to your kasoku-server URL:

```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

3. **Run development server**:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Building for Production

```bash
npm run build
npm start
```

## Deployment

### Vercel (Recommended)

1. Push code to GitHub
2. Import project in Vercel
3. Set environment variable: `NEXT_PUBLIC_API_URL=https://api.kasoku.dev`
4. Deploy

### Docker

```bash
# Build
docker build -t kasoku-web .

# Run
docker run -p 3000:3000 -e NEXT_PUBLIC_API_URL=http://your-server:8080 kasoku-web
```

### Self-Hosted

```bash
# Build static export
npm run build

# Serve with any static server
npx serve@latest out
```

## Features

### Landing Page
- Product overview
- Feature highlights
- Pricing tiers
- Call-to-action buttons

### Authentication
- Sign up with email/password
- Login
- Secure token storage in localStorage

### Dashboard
- **Analytics**: View cache hit rates, time saved, storage usage
- **Cache**: Browse and manage cache entries
- **Tokens**: Create and revoke API tokens for CLI

## Tech Stack

- **Next.js 14** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Radix UI** - Accessible components
- **Recharts** - Analytics charts (ready for future enhancements)

## API Integration

The dashboard connects to kasoku-server via REST API:

- `POST /auth/register` - Create account
- `POST /auth/login` - Authenticate
- `GET /auth/me` - Get current user
- `POST /auth/tokens` - Create API token
- `GET /auth/tokens` - List tokens
- `DELETE /auth/tokens/:id` - Revoke token
- `GET /cache` - List cache entries
- `DELETE /cache/:hash` - Delete cache entry
- `GET /analytics` - Get analytics summary

## Future Enhancements

- [ ] Team management
- [ ] Billing integration (Stripe)
- [ ] Advanced analytics with charts
- [ ] Cache usage breakdown by command
- [ ] Real-time updates with WebSockets
- [ ] Dark mode toggle
- [ ] Email verification
- [ ] Password reset flow
