# Multi-stage build for production deployment
FROM node:18-alpine as backend-build

# Set working directory for backend
WORKDIR /app/backend

# Copy backend package files
COPY backend/package*.json ./

# Install backend dependencies
RUN npm ci --only=production

# Copy backend source code
COPY backend/ ./

# Frontend build stage
FROM node:18-alpine as frontend-build

# Set working directory for frontend
WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./

# Install frontend dependencies
RUN npm ci

# Copy frontend source code
COPY frontend/ ./

# Build frontend
ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=$REACT_APP_API_URL
RUN npm run build

# Production stage
FROM node:18-alpine as production

# Set working directory
WORKDIR /app

# Copy backend from backend-build stage
COPY --from=backend-build /app/backend ./backend

# Copy frontend build from frontend-build stage
COPY --from=frontend-build /app/frontend/build ./frontend/build

# Install serve to serve static files
RUN npm install -g serve

# Create startup script
RUN echo '#!/bin/sh' > start.sh && \
    echo 'serve -s frontend/build -l ${FRONTEND_PORT:-3000} &' >> start.sh && \
    echo 'cd backend && node server.js' >> start.sh && \
    chmod +x start.sh

# Expose ports
EXPOSE 5000 3000

# Set environment variables
ENV NODE_ENV=production

# Start both frontend and backend
CMD ["./start.sh"]
