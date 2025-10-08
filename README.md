# 🚀 SmartPicks Backend API

API RESTful para o sistema SmartPicks com controle de usuários, perfis administrativos e upload de avatar.

## 📋 Índice

- [🚀 Como Fazer o Clone e Setup](#-como-fazer-o-clone-e-setup)
- [✨ Características](#-características)
- [📋 Pré-requisitos](#-pré-requisitos) 
- [🔧 Instalação Passo a Passo](#-instalação-passo-a-passo)
- [🗄️ Configuração do Banco](#️-configuração-do-banco)
- [⚙️ Configuração de Ambiente](#️-configuração-de-ambiente)
- [🚀 Executando o Projeto](#-executando-o-projeto)
- [📚 Documentação da API](#-documentação-da-api)
- [🔍 Endpoints Disponíveis](#-endpoints-disponíveis)
- [🏗️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🛠️ Tecnologias](#️-tecnologias)
- [🔧 Troubleshooting](#-troubleshooting)

## 🚀 Como Fazer o Clone e Setup

### 1. **Clone o Repositório**
```bash
git clone https://github.com/AdrianResende/backend-SmartPicks.git
cd backend-SmartPicks
```

### 2. **Instale o Go** (se ainda não tiver)
- Baixe e instale Go 1.19+ em: https://golang.org/dl/
- Verifique a instalação: `go version`

### 3. **Instale o MySQL** (se ainda não tiver)
- MySQL 8.0+: https://dev.mysql.com/downloads/mysql/
- Ou use Docker: `docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -p 3306:3306 -d mysql:8.0`

### 4. **Configure o Projeto** (siga os passos abaixo)

## ✨ Características

- ✅ **Sistema de Autenticação** com login/registro
- ✅ **Controle de Perfis** (admin/user) com validação ENUM
- ✅ **Upload de Avatar** com suporte a Base64 e URLs
- ✅ **Criptografia de Senhas** com bcrypt
- ✅ **Validação de Permissões** por perfil
- ✅ **CORS Configurado** para frontend
- ✅ **Documentação Swagger** automática e versionada
- ✅ **Validação de Dados** robusta
- ✅ **MySQL** como banco de dados
- ✅ **Handlers Organizados** por funcionalidade
- ✅ **Migração Automática** de schema no startup

## 📋 Pré-requisitos

- **Go** 1.19+ ([Download](https://golang.org/dl/))
- **MySQL** 8.0+ ([Download](https://dev.mysql.com/downloads/mysql/))
- **Git** (opcional, mas recomendado)

## 🔧 Instalação Passo a Passo

### **Passo 1: Baixar Dependências**
```bash
go mod tidy
```

### **Passo 2: Instalar Swag (Opcional - para regenerar docs)**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🗄️ Configuração do Banco

### **Passo 3: Criar o Banco de Dados**
```sql
-- Conecte ao MySQL como root
mysql -u root -p

-- Crie o banco de dados
CREATE DATABASE smartpicks CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### **Passo 4: Execute o Script de Setup**
```bash
# Opção A: Script para primeira instalação (recomendado)
mysql -u root -p smartpicks < database_setup_simple.sql

# Opção B: Script completo com condicionais
mysql -u root -p < database_setup.sql
```

### **Estrutura da Tabela Criada**
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    data_nascimento DATE NOT NULL,
    perfil ENUM('admin', 'user') NOT NULL DEFAULT 'user',
    avatar MEDIUMTEXT,  -- Suporte a Base64 grande
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_email (email),
    INDEX idx_cpf (cpf),
    INDEX idx_perfil (perfil)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## ⚙️ Configuração de Ambiente

### **Passo 5: Criar Arquivo .env**
```bash
# Copie o exemplo e configure
cp .env.example .env
```

**Edite o arquivo `.env`:**
```env
# Configuração do Banco de Dados
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=sua_senha_mysql
DB_NAME=smartpicks

# Configuração do Servidor
PORT=8080

# Configuração CORS (ajuste conforme seu frontend)
CORS_ORIGIN=http://localhost:9000
```

## 🚀 Executando o Projeto

### **Passo 6: Testar a Conexão**
```bash
go run main.go
```

**Se tudo estiver correto, você verá:**
```
2025/10/08 10:00:00 Servidor rodando na porta 8080
2025/10/08 10:00:00 Documentação Swagger disponível em: http://localhost:8080/swagger/
```

### **Passo 7: Verificar se Funcionou**
```bash
# Teste básico
curl http://localhost:8080/

# Deve retornar: {"status": "API rodando", "version": "1.0.0"}
```

### **Comandos Úteis:**
```bash
# Executar em modo desenvolvimento
go run main.go

# Compilar e executar
go build
./smartpicks-backend.exe    # Windows
./smartpicks-backend        # Linux/Mac

# Executar com port customizada
PORT=8081 go run main.go
```

## 📚 Documentação da API

### **Swagger UI:** 
🌐 **http://localhost:8080/swagger/**

### **Health Checks:**
- **Status:** `GET /` → `{"status": "API rodando", "version": "1.0.0"}`
- **Health:** `GET /health` → `{"status": "healthy"}`

## � Endpoints Disponíveis

### 🔐 **Autenticação**

| Método | Endpoint | Descrição | Body |
|--------|----------|-----------|------|
| `POST` | `/api/login` | Login de usuário | `{email, password}` |
| `POST` | `/api/register` | Cadastro de usuário | `{nome, email, password, cpf, data_nascimento, perfil?}` |

### 👥 **Usuários**

| Método | Endpoint | Descrição | Parâmetros |
|--------|----------|-----------|------------|
| `GET` | `/api/users` | Listar todos usuários | - |
| `GET` | `/api/users/permissions` | Verificar permissões | `?email=usuario@email.com` |
| `GET` | `/api/users/profile` | Usuários por perfil | `?profile=admin` ou `?profile=user` |

### 🖼️ **Avatar (Upload de Imagem)**

| Método | Endpoint | Descrição | Body |
|--------|----------|-----------|------|
| `POST` | `/api/users/avatar` | Upload/criar avatar | `{user_id: 1, avatar: "base64_ou_url"}` |
| `PUT` | `/api/users/avatar` | Atualizar avatar | `{user_id: 1, avatar: "base64_ou_url"}` |
| `DELETE` | `/api/users/avatar` | Remover avatar | `{user_id: 1}` |

### 📝 **Exemplos de Requisições**

**Cadastro:**
```json
POST /api/register
{
  "nome": "João Silva",
  "email": "joao@exemplo.com",
  "password": "senha123",
  "cpf": "123.456.789-10",
  "data_nascimento": "1990-05-15",
  "perfil": "user"
}
```

**Login:**
```json
POST /api/login
{
  "email": "joao@exemplo.com",
  "password": "senha123"
}
```

**Resposta de Sucesso:**
```json
{
  "id": 1,
  "nome": "João Silva",
  "email": "joao@exemplo.com",
  "cpf": "123.456.789-10",
  "data_nascimento": "1990-05-15",
  "perfil": "user",
  "is_admin": false,
  "has_permission": true,
  "created_at": "2025-09-30T20:11:21Z",
  "updated_at": "2025-09-30T20:11:21Z"
}
```

## � Estrutura do Projeto

```
BackEnd/
├── 📄 main.go                 # Ponto de entrada
├── 📄 go.mod                  # Dependências Go
├── 📄 go.sum                  # Lock das dependências
├── 📄 database_setup.sql      # Script do banco
├── 📄 .gitignore             # Arquivos ignorados
├── 📄 README.md              # Esta documentação
├── 📁 internal/
│   ├── 📁 database/
│   │   └── 📄 db.go          # Conexão com MySQL
│   ├── 📁 handlers/
│   │   └── 📄 auth.go        # Handlers da API
│   ├── 📁 models/
│   │   └── 📄 user.go        # Modelos de dados
│   └── 📁 routes/
│       └── 📄 routes.go      # Configuração de rotas
└── 📁 docs/                  # Documentação Swagger (gerada)
    ├── 📄 docs.go
    ├── 📄 swagger.json
    └── 📄 swagger.yaml
```

## 🛠️ Tecnologias

- **[Go 1.19+](https://golang.org/)** - Linguagem de programação moderna e performática
- **[Gorilla Mux](https://github.com/gorilla/mux)** - Roteador HTTP robusto
- **[MySQL 8.0+](https://www.mysql.com/)** - Banco de dados relacional
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Criptografia segura de senhas
- **[Swag](https://github.com/swaggo/swag)** - Geração automática de documentação OpenAPI
- **[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)** - Driver oficial MySQL

## � Troubleshooting

### **Problema: "Tabela users não encontrada"**
```bash
# Execute o script de setup do banco
mysql -u root -p smartpicks < database_setup_simple.sql
```

### **Problema: "connection refused" ou "Error 1045"**
```bash
# Verifique se o MySQL está rodando
mysql -u root -p

# Verifique suas credenciais no .env
DB_USER=root
DB_PASSWORD=sua_senha_real
```

### **Problema: "Port 8080 already in use"**
```bash
# Use uma porta diferente
PORT=8081 go run main.go

# Ou mate o processo na porta 8080
netstat -ano | findstr :8080
taskkill /PID XXXX /F
```

### **Problema: Documentação Swagger não aparece**
```bash
# Reinstale swag e regenere
go install github.com/swaggo/swag/cmd/swag@latest
swag init
go run main.go
```

### **Problema: CORS errors no frontend**
```bash
# Ajuste o CORS_ORIGIN no .env para o seu frontend
CORS_ORIGIN=http://localhost:3000  # React
CORS_ORIGIN=http://localhost:9000  # Sua aplicação
```

### **Regenerar Documentação Swagger (Opcional)**
```bash
# Se você modificar os comentários dos handlers
swag init -g main.go -o docs
```

## 🔒 Perfis e Permissões

### **Perfil: `user` (padrão)**
- ✅ Acesso básico ao sistema
- ✅ Upload/atualização do próprio avatar
- ✅ Visualização de dados próprios

### **Perfil: `admin`**
- ✅ Todas as permissões de `user`
- ✅ Gerenciar todos os usuários e avatars
- ✅ Acesso a funcionalidades administrativas

## 🚀 Para Produção

### **Variáveis de Ambiente Recomendadas:**
```env
# .env para produção
DB_HOST=seu-servidor-mysql.com
DB_PORT=3306
DB_USER=smartpicks_user
DB_PASSWORD=senha_segura_aleatoria
DB_NAME=smartpicks

PORT=8080
CORS_ORIGIN=https://seufrontend.com

# Opcional: configurações de conexão
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
```

### **Compile para Produção:**
```bash
# Linux/Ubuntu
GOOS=linux GOARCH=amd64 go build -o smartpicks-backend

# Windows
GOOS=windows GOARCH=amd64 go build -o smartpicks-backend.exe
```

---

**✅ API SmartPicks Backend pronta para uso!**  
**🌐 Acesse a documentação completa em: http://localhost:8080/swagger/**
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