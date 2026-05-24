# Docker-in-Docker (DinD) vs Docker-out-of-Docker (DooD)

Este diretório contém uma configuração pronta para executar o Docker dentro do Docker. Dependendo do seu caso de uso (como CI/CD, sandboxing ou desenvolvimento local), existem duas abordagens principais.

---

## Opção 1: Docker-in-Docker (DinD) Real
Nesta abordagem, criamos um daemon do Docker totalmente isolado e independente rodando **dentro** do container. Ele não compartilha nenhuma imagem ou container com o seu Docker do host.

O arquivo `docker-compose.yml` neste diretório já está configurado com esta abordagem:

```yaml
version: '3.8'

services:
  docker-daemon:
    image: docker:dind
    privileged: true
    environment:
      DOCKER_TLS_CERTDIR: ""
    ports:
      - "2375:2375"
...
```

### Como Executar:
1. Navegue até esta pasta:
   ```bash
   cd docker-in-docker
   ```
2. Inicie os containers:
   ```bash
   docker compose up -d
   ```
3. Acompanhe os logs para ver se o cliente conseguiu se comunicar com o daemon interno:
   ```bash
   docker compose logs dind-client
   ```
4. Para entrar no container cliente e rodar comandos do Docker (como baixar imagens ou subir containers de teste de forma isolada):
   ```bash
   docker compose exec docker-client sh
   # Dentro do container:
   docker run hello-world
   ```

---

## Opção 3: Orquestrando Frontend e Backend Separados (Aninhado em DinD)

Se você quer subir a sua arquitetura completa (Frontend em React, Backend em Go, Banco de Dados Postgres e Redis) dentro do ambiente isolado do Docker-in-Docker, a estrutura já está pronta para isso!

Os arquivos `backend/Dockerfile`, `frontend/Dockerfile` e `docker-in-docker/docker-compose.nested.yml` foram criados para gerenciar todo o ecossistema de forma aninhada.

### Como funciona?
O Docker do seu computador (Host) monta o código fonte na pasta `/src` dos dois contêineres principais (`dind-daemon` e `dind-client`). Como eles compartilham a mesma pasta no mesmo caminho exato, o cliente consegue solicitar ao daemon que compile os Dockerfiles locais!

### Como Rodar o Ecossistema Aninhado:

1. **Suba os contêineres do DinD principal:**
   ```bash
   cd docker-in-docker
   docker compose up -d
   ```

2. **Rode o Docker Compose aninhado de dentro do contêiner cliente:**
   ```bash
   docker exec -it dind-client docker compose -f /src/docker-in-docker/docker-compose.nested.yml up --build -d
   ```

3. **Verifique os contêineres aninhados rodando dentro do daemon interno:**
   ```bash
   docker exec -it dind-client docker ps
   ```

4. **Acesse as aplicações:**
   * **Frontend:** Acesse no seu navegador em `http://localhost:3000` (porta 3000 do host redireciona para a porta 3000 do daemon interno, que aponta para o contêiner do frontend).
   * **API Backend (HTTP):** Acesse em `http://localhost:8080` do seu computador.
   * **API Backend (gRPC):** Porta `50051`.

### Para derrubar o ecossistema aninhado:
```bash
docker exec -it dind-client docker compose -f /src/docker-in-docker/docker-compose.nested.yml down
```

