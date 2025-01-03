# Build stage
FROM node:22-alpine as build

# Set the working directory
WORKDIR /app

# Install dependencies
COPY client/package.json client/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install

# Copy the frontend source code
COPY client/ ./

# Build the production build
RUN pnpm build

# Serve the frontend with Nginx
FROM nginx:alpine

# Copy the built frontend into the Nginx container
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80 for HTTP
EXPOSE 80

# Run Nginx
CMD ["nginx", "-g", "daemon off;"]
