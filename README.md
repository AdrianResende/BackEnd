# SmartPicks Backend

API backend para o sistema SmartPicks desenvolvida em Go.

## 🚀 Como executar

### Pré-requisitos
- Go 1.21 ou superior
- MySQL 5.7 ou superior

### Configuração do Banco de Dados
1. Execute o script SQL que está no arquivo `database_setup.sql`
2. Ou execute manualmente os comandos SQL abaixo:

```sql
CREATE DATABASE IF NOT EXISTS meuprojeto;
USE meuprojeto;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    data_nascimento DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_email (email),
    INDEX idx_cpf (cpf)
);
```

### Configuração do Ambiente
1. Copie o arquivo `.env.example` para `.env`
2. Configure suas credenciais do banco de dados no arquivo `.env`

### Executando a aplicação
```bash
# Baixar dependências
go mod tidy

# Compilar
go build

# Executar
go run main.go
```

A API estará disponível em `http://localhost:8080`

## 📡 Endpoints

### Health Check
- **GET** `/` - Retorna status da API
- **GET** `/health` - Health check

### Autenticação
- **POST** `/api/login` - Login do usuário
- **POST** `/api/register` - Cadastro de novo usuário

Exemplo de payload para login:
```json
{
    "email": "usuario@exemplo.com",
    "password": "sua_senha"
}
```

Exemplo de payload para registro:
```json
{
    "nome": "João Silva",
    "email": "joao@exemplo.com",
    "password": "minhasenha123",
    "cpf": "123.456.789-10",
    "data_nascimento": "1990-05-15"
}
```

## 🔧 Estrutura do Projeto

```
BackEnd/
├── main.go                 # Ponto de entrada da aplicação
├── go.mod                  # Dependências do Go
├── go.sum                  # Checksums das dependências
├── .env.example            # Exemplo de configuração
└── internal/
    ├── database/
    │   └── db.go           # Configuração do banco
    ├── handlers/
    │   └── auth.go         # Handlers de autenticação
    ├── models/
    │   └── user.go         # Modelo do usuário
    └── routes/
        └── routes.go       # Configuração das rotas
```

## ⚠️ Notas de Segurança

- As senhas devem ser armazenadas com hash bcrypt
- Configure variáveis de ambiente para credenciais do banco
- Implemente autenticação JWT para sessões
- Adicione rate limiting para proteção contra ataques

## 🛠️ Próximos Passos

- [ ] Implementar JWT para autenticação
- [ ] Adicionar endpoint de registro
- [ ] Implementar middleware de logging
- [ ] Adicionar testes unitários
- [ ] Configurar Docker