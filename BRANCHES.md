# Estrutura de Branches - SmartPicks Backend

## 📋 Branches do Projeto

### 🚀 `production`
**Ambiente:** Vercel (Produção)  
**URL:** https://backend-smart-picks.vercel.app  
**Banco de dados:** PostgreSQL (Neon/Vercel)  

**Uso:**
- Branch principal para deploy em produção
- Conectada automaticamente com a Vercel
- Apenas código testado e aprovado

**Deploy:**
```bash
git checkout production
git merge main  # ou development após testes
git push origin production
```

A Vercel fará o deploy automático ao detectar push nesta branch.

---

### 💻 `development`
**Ambiente:** Local (Desenvolvimento)  
**URL:** http://localhost:8080  
**Banco de dados:** PostgreSQL local ou MySQL local  

**Uso:**
- Branch para desenvolvimento local
- Testes e desenvolvimento de novas features
- Usa o arquivo `main.go` para rodar localmente

**Como usar:**
```bash
# Mudar para a branch de desenvolvimento
git checkout development

# Rodar o servidor local
go run main.go

# Ou compilar e executar
go build -o main.exe
./main.exe
```

---

### 🌿 `main`
**Ambiente:** Branch principal  
**Sincronização:** Com production e development  

**Uso:**
- Branch principal do repositório
- Recebe merges de development após aprovação
- Fonte para production

**Workflow:**
```bash
development -> main -> production
```

---

## 🔄 Workflow Recomendado

### Desenvolvendo uma nova feature:

1. **Criar branch a partir de development:**
```bash
git checkout development
git pull origin development
git checkout -b feature/nome-da-feature
```

2. **Desenvolver e testar localmente:**
```bash
go run main.go
# Testar a aplicação em http://localhost:8080
```

3. **Commit e merge para development:**
```bash
git add .
git commit -m "feat: descrição da feature"
git checkout development
git merge feature/nome-da-feature
```

4. **Testar em development e depois merge para main:**
```bash
git checkout main
git merge development
git push origin main
```

5. **Deploy para production:**
```bash
git checkout production
git merge main
git push origin production
```

---

## 📁 Estrutura do Projeto

```
backend-SmartPicks/
├── api/
│   └── index.go          # Handler para Vercel (usado em production)
├── internal/
│   ├── database/         # Conexões com banco de dados
│   ├── handlers/         # Handlers das rotas (usado em development)
│   ├── models/          # Modelos de dados
│   └── routes/          # Definição de rotas (usado em development)
├── main.go              # Arquivo principal para desenvolvimento local
├── go.mod               # Dependências Go
├── vercel.json          # Configuração da Vercel
└── .vercelignore        # Arquivos ignorados no deploy
```

---

## 🔧 Variáveis de Ambiente

### Production (Vercel):
- `DATABASE_URL` - String de conexão PostgreSQL

### Development (Local):
Criar arquivo `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua_senha
DB_NAME=smartpicks
PORT=8080
```

---

## 📝 Comandos Úteis

```bash
# Ver todas as branches
git branch -a

# Mudar de branch
git checkout <nome-da-branch>

# Ver diferenças entre branches
git diff development production

# Enviar todas as branches para o repositório
git push origin --all

# Atualizar branch local
git pull origin <nome-da-branch>
```

---

## ✅ Branches Atuais

- ✅ `main` - Branch principal
- ✅ `production` - Deploy Vercel (produção)
- ✅ `development` - Desenvolvimento local

---

## 🚀 APIs Disponíveis

### Produção (Vercel)
- Base URL: `https://backend-smart-picks.vercel.app`

### Desenvolvimento Local
- Base URL: `http://localhost:8080`

### Endpoints:
- `GET /` - Status da API
- `GET /health` - Health check
- `POST /api/register` - Cadastro
- `POST /api/login` - Login
- `GET /api/users` - Listar usuários
- `GET /api/users/permissions` - Verificar permissões
- `GET /api/users/profile` - Buscar por perfil
- `POST /api/users/avatar` - Atualizar avatar
- `DELETE /api/users/avatar` - Remover avatar
