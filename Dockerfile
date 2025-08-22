# Frontend builder
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend .
RUN npm run build && npm run export || echo "ensure next export configured"

# Nginx runtime (proxy to external gateway service)
FROM nginx:alpine
RUN adduser -D -H -s /sbin/nologin wohnfair
COPY infra/nginx/nginx.conf /etc/nginx/nginx.conf
COPY --from=frontend-builder /app/out /usr/share/nginx/html
USER wohnfair
EXPOSE 80
HEALTHCHECK --interval=30s --timeout=10s --start-period=20s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost/healthz || exit 1
CMD ["nginx", "-g", "daemon off;"]
