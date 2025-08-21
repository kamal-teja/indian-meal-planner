# Multi-stage build for production deployment

# Backend build stage
FROM golang:1.21-alpine as backend-build

# Set working directory for backend
WORKDIR /app/backend

# Copy backend go module files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build the Go application
RUN go build -o main cmd/server/main.go

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
FROM alpine:latest as production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy backend binary from backend-build stage
COPY --from=backend-build /app/backend/main ./backend/

# Copy frontend build from frontend-build stage
COPY --from=frontend-build /app/frontend/build ./frontend/build

# Install serve to serve static files
FROM node:18-alpine
WORKDIR /app
COPY --from=frontend-build /app/frontend/build ./frontend/build
COPY --from=backend-build /app/backend/main ./backend/
RUN npm install -g serve

# Create startup script
RUN echo '#!/bin/sh' > start.sh && \
    echo 'serve -s frontend/build -l ${FRONTEND_PORT:-3000} &' >> start.sh && \
    echo './backend/main' >> start.sh && \
    chmod +x start.sh

# Expose ports
EXPOSE 5000 3000

# Set environment variables
ENV ENVIRONMENT=production

# Start both frontend and backend
CMD ["./start.sh"]
