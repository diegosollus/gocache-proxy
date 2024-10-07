# GoCache Proxy

GoCache Proxy é um servidor proxy reverso escrito em Go que inclui funcionalidades para:

- **Bloquear IPs específicos**.
- **Manipular subdomínios**, redirecionando URLs com o subdomínio `www` para um caminho específico.
- **Verificar e bloquear conteúdo malicioso** em query strings e payloads (HTML malicioso, como `<script>` tags).

## Funcionalidades

- **Bloqueio de IPs**: Detecta e bloqueia requisições de IPs que estão em uma lista de IPs bloqueados.
- **Manipulação de subdomínios**: Reescreve o caminho da URL para requisições que vêm do subdomínio `www`, redirecionando para um caminho específico.
- **Proteção contra HTML malicioso**: Detecta e bloqueia requisições que contenham tags HTML maliciosas na query string ou no corpo da requisição (para POST e PUT).

## Estrutura do Projeto

```bash
gocache-proxy/
├── cmd/
│   └── main.go               # Ponto de entrada da aplicação
├── data/
│   └── config.yaml           # Estrutura de servidores para proxy reverso
├── db/
│   └── db.go                 # Configuração do banco de dados (SQLite)
├── internal/
│   ├── utils/
│       └── domain.go         # Funções utilitárias para IPs e subdomínios
│       └── ipfilter.go 
│   ├── server/
│   │   └── server.go         # Implementação principal do proxy reverso
│
├── security/
│   └── htmltags.go           # Verificação de HTML malicioso com regex
│
└── README.go                 # Documentação da aplicação
```

## Instalação

1. **Clone o repositório**:

   ```bash
   git clone https://github.com/diegosollus/gocache-proxy.git
   cd gocache-proxy
   ```

2. **Instale as dependências**:
   ```bash
   go mod tidy
   ```

3. **Configure a aplicação**:

## Uso

### Rodando o servidor

Para iniciar o servidor proxy reverso:

```bash
go run ./cmd/main.go
```

O servidor será iniciado e escutará na porta padrão (ex: `localhost:8080`). Ele redirecionará as requisições para o host de destino configurado no arquivo config.yaml.

### Principais Endpoints

A aplicação age como um proxy, portanto, qualquer requisição HTTP/HTTPS pode ser passada pelo servidor. A lógica de bloqueio de IP, manipulação de subdomínios e detecção de HTML malicioso será aplicada antes que a requisição seja redirecionada para o backend.

### Exemplo de manipulação de subdomínio

- Requisição para `http://www.exemplo.com/teste` será redirecionada para `/site/www/teste`.

### Exemplo de bloqueio de HTML malicioso

- Qualquer requisição contendo tags maliciosas, como `<script>`, na query string ou no corpo (em POST ou PUT), será bloqueada com uma resposta HTTP `403 Forbidden`.

## Testes

Para rodar os testes:

```bash
go test ./...
```

Isso executará os testes em todos os pacotes da aplicação.
