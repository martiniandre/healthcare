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

> [!WARNING]
> **Implicações de Segurança**: O DinD exige o parâmetro `privileged: true`. Isso dá ao container acesso quase total ao hardware do host. Use com cautela em ambientes de produção.

---

## Opção 2: Docker-out-of-Docker (DooD)
Se o seu objetivo é apenas conseguir rodar comandos `docker` dentro de um container (por exemplo, em um container do Jenkins para gerar imagens), a abordagem DooD costuma ser a mais recomendada. 

Em vez de rodar outro daemon, nós **compartilhamos o socket do Docker do host** com o container. Isso significa que qualquer container criado por ele rodará como um "irmão" no seu próprio host, e não dentro dele.

### Exemplo de Configuração DooD (`docker-compose.dood.yml`):

```yaml
version: '3.8'

services:
  dood-client:
    image: docker:latest
    container_name: dood-client
    volumes:
      # Monta o socket do Docker do seu computador no container
      - /var/run/docker.sock:/var/run/docker.sock
    # Mantém o container vivo
    entrypoint: ["sh", "-c", "echo 'Pronto! Acessando o Docker do host...' && docker ps && tail -f /dev/null"]
```

### Vantagens do DooD:
* **Segurança**: Não precisa de `privileged: true`.
* **Desempenho**: Compartilha o cache de imagens do host.
* **Simplicidade**: Qualquer container gerado é visível diretamente na sua máquina local com `docker ps`.
