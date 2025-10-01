# 🚀 SmartPicks Backend API

API RESTful para o sistema SmartPicks com controle de usuários e perfis administrativos.

## 📋 Índice

- [Características](#-características)
- [Pré-requisitos](#-pré-requisitos) 
- [Instalação](#-instalação)
- [Configuração do Banco](#-configuração-do-banco)
- [Executando](#-executando)
- [Documentação da API](#-documentação-da-api)
- [Endpoints](#-endpoints)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Tecnologias](#-tecnologias)

## ✨ Características

- ✅ **Sistema de Autenticação** com login/registro
- ✅ **Controle de Perfis** (admin/user) com validação ENUM
- ✅ **Criptografia de Senhas** com bcrypt
- ✅ **Validação de Permissões** por perfil
- ✅ **CORS Configurado** para frontend
- ✅ **Documentação Swagger** automática
- ✅ **Validação de Dados** robusta
- ✅ **MySQL** como banco de dados

## 📋 Pré-requisitos

- **Go** 1.19+ 
- **MySQL** 8.0+
- **Git** (opcional)

## 🔧 Instalação

1. **Clone o repositório:**
```bash
git clone <url-do-repositorio>
cd SmartPicks/BackEnd
```

2. **Instale as dependências:**
```bash
go mod tidy
```

3. **Instale o swag (para gerar documentação):**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## 🗄️ Configuração do Banco

1. **Execute o MySQL e crie o banco:**
```sql
CREATE DATABASE smartpicks;
```

2. **Execute o script de setup:**
```bash
mysql -u root -p < database_setup.sql
```

**Ou execute manualmente:**
```sql
USE smartpicks;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    data_nascimento DATE NOT NULL,
    perfil ENUM('admin', 'user') NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_email (email),
    INDEX idx_cpf (cpf),
    INDEX idx_perfil (perfil)
);
```

## 🚀 Executando

1. **Gerar documentação Swagger:**
```bash
swag init
```

2. **Executar o servidor:**
```bash
go run main.go
```

3. **Ou compilar e executar:**
```bash
go build
./smartpicks-backend.exe
```

**Servidor rodará em:** `http://localhost:8080`

## � Documentação da API

### **Swagger UI:** 
🌐 **http://localhost:8080/swagger/**

### **Health Checks:**
- **Status:** `GET /`
- **Health:** `GET /health`

## 🔗 Endpoints

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
| `GET` | `/api/users/check` | Status da tabela (debug) | - |

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

- **[Go](https://golang.org/)** - Linguagem de programação
- **[Gorilla Mux](https://github.com/gorilla/mux)** - Roteador HTTP
- **[MySQL](https://www.mysql.com/)** - Banco de dados
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Criptografia de senhas
- **[Swagger](https://swagger.io/)** - Documentação da API
- **[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)** - Driver MySQL

## 🔒 Perfis e Permissões

### **Perfil: `user` (padrão)**
- ✅ Acesso básico ao sistema
- ✅ Visualização de dados próprios
- ❌ Não pode gerenciar outros usuários

### **Perfil: `admin`**
- ✅ Todas as permissões de `user`
- ✅ Gerenciar todos os usuários
- ✅ Acesso a funcionalidades administrativas
- ✅ Visualizar relatórios e logs

## 🐛 Debug e Desenvolvimento

### **Verificar tabela users:**
```bash
GET http://localhost:8080/api/users/check
```

### **Logs do servidor:**
O servidor exibe logs detalhados no console, incluindo:
- Conexão com banco de dados
- Requisições HTTP
- Erros e avisos

### **Regenerar documentação Swagger:**
```bash
swag init
```

## 📞 Suporte

Para dúvidas ou problemas:
- 📧 Email: suporte@smartpicks.com
- 🐛 Issues: [GitHub Issues](link-para-issues)

---

**Desenvolvido com ❤️ para SmartPicks**
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