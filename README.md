# Workout Tracker

A full-stack workout tracking application built with Go (Gin), Vue.js, PostgreSQL, and Biome.js.

## Tech Stack

### Backend
- **Go 1.23** with Gin framework
- **PostgreSQL 18** for database
- RESTful API architecture

### Frontend
- **Vue 3** with Composition API
- **Vite** for development and building
- **Biome.js** for linting and formatting
- **Vue Router** for navigation
- **Axios** for API calls

### Infrastructure
- **Docker & Docker Compose** for containerization
- PostgreSQL in Docker

## Project Structure

```
workout-tracker/
├── main.go                 # Go backend entry point
├── go.mod                  # Go dependencies
├── Dockerfile             # Backend container
├── docker-compose.yml     # Docker orchestration
├── init.sql              # Database initialization
└── frontend/
    ├── package.json      # Node dependencies
    ├── biome.json       # Biome configuration
    ├── vite.config.js   # Vite configuration
    ├── index.html       # HTML entry point
    └── src/
        ├── main.js      # Vue app entry
        ├── App.vue      # Root component
        └── views/       # Page components
```

## Getting Started

### Prerequisites
- Go 1.23 or later
- Node.js 18 or later
- Docker and Docker Compose
- Make (optional)

### Backend Setup

1. **Install Go dependencies:**
   ```bash
   go mod download
   ```

2. **Start PostgreSQL with Docker:**
   ```bash
   docker-compose up postgres -d
   ```

3. **Run the backend:**
   ```bash
   go run main.go
   ```

   The backend will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Run development server:**
   ```bash
   npm run dev
   ```

   The frontend will be available at `http://localhost:3000`

### Using Docker Compose (Recommended)

Start everything with one command:

```bash
docker-compose up --build
```

This will start:
- PostgreSQL on port 5432
- Backend API on port 8080

For the frontend, run separately:
```bash
cd frontend
npm install
npm run dev
```

## Development

### Backend Development

The Go backend uses Gin framework and connects to PostgreSQL. Main features:
- Health check endpoint: `GET /health`
- Workouts API: `GET /api/workouts`, `POST /api/workouts`
- CORS enabled for frontend development

### Frontend Development

The Vue frontend includes:
- Vue Router for navigation
- Axios for API calls
- Biome.js for code quality

**Available Scripts:**
```bash
npm run dev        # Start dev server
npm run build      # Build for production
npm run preview    # Preview production build
npm run lint       # Lint with Biome
npm run lint:fix   # Fix linting issues
npm run format     # Format code
```

### Database

PostgreSQL is configured with:
- Database: `workout_tracker`
- User: `postgres`
- Password: `postgres`
- Port: `5432`

The `init.sql` file creates initial tables:
- `workouts` - Store workout sessions
- `exercises` - Store individual exercises

## API Endpoints

### Health Check
```
GET /health
```

### Workouts
```
GET  /api/workouts      # Get all workouts
POST /api/workouts      # Create a new workout
```

## Environment Variables

### Backend
- `DATABASE_URL` - PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/workout_tracker?sslmode=disable`)
- `PORT` - Server port (default: `8080`)

## Production Build

### Backend
```bash
docker build -t workout-tracker-backend .
docker run -p 8080:8080 workout-tracker-backend
```

### Frontend
```bash
cd frontend
npm run build
# Serve the 'dist' directory with your preferred static server
```

## License

MIT
