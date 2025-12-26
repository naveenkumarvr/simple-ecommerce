# Simple Dockerfile for frontend-nodejs (Express)

# Use a small Node image
FROM node:18-alpine

# Create app directory
WORKDIR /app

# Install dependencies
COPY package.json package-lock.json* ./
RUN npm install --production || npm install

# Copy application source
COPY . .

# Environment
ENV NODE_ENV=production
# API gateway base URL; override at runtime if needed
ENV API_BASE=http://api-gateway-go:8000/api

EXPOSE 3000

CMD ["npm", "start"]
