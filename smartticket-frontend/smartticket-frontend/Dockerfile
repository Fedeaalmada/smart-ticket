# ── Etapa 1: compilar el frontend React ──────────────────────
FROM node:20-alpine AS builder

WORKDIR /app

# Copiar dependencias primero (mejor cache)
COPY package.json package-lock.json* ./
RUN npm install

# Copiar el resto y compilar
COPY . .
RUN npm run build

# ── Etapa 2: servir con Nginx ─────────────────────────────────
FROM nginx:alpine

# Copiar la build al directorio de Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Configuración para que React Router funcione y proxy al backend
RUN printf 'server {\n\
    listen 80;\n\
    location /api/ {\n\
        proxy_pass http://backend:8080/;\n\
        proxy_set_header Host $host;\n\
        proxy_set_header X-Real-IP $remote_addr;\n\
    }\n\
    location / {\n\
        root /usr/share/nginx/html;\n\
        index index.html;\n\
        try_files $uri $uri/ /index.html;\n\
    }\n\
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
