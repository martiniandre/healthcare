# --- Estágio de Compilação ---
FROM node:20-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY package.json package-lock.json ./
RUN npm ci

# Copiar todo o código-fonte
COPY . .

# Compilar a aplicação para produção
RUN npm run build

# --- Estágio de Execução (Nginx) ---
FROM nginx:alpine

# Copiar os arquivos estáticos gerados no build para a pasta padrão do Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copiar uma configuração customizada de Nginx se necessário para SPA (React Router)
# Para evitar erros de 404 ao atualizar rotas SPA, criamos um arquivo básico de configuração
RUN echo 'server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html index.htm; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
