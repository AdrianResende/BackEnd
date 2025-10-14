-- PostgreSQL Schema para SmartPicks
-- Execute este script no seu banco PostgreSQL (Neon)

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    data_nascimento DATE NOT NULL,
    perfil VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (perfil IN ('admin', 'user')),
    avatar TEXT NULL DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_cpf ON users (cpf);
CREATE INDEX IF NOT EXISTS idx_users_perfil ON users (perfil);

-- Trigger para updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Inserir usuários padrão (se não existirem)
INSERT INTO users (nome, email, password, cpf, data_nascimento, perfil, avatar) VALUES
(
    'Admin User',
    'admin@smartpicks.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    '00000000000',
    '1990-01-01',
    'admin',
    'https://ui-avatars.com/api/?name=Admin+User&background=0d8abc&color=fff'
),
(
    'User Comum',
    'user@smartpicks.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    '11111111111',
    '1995-06-15',
    'user',
    'https://ui-avatars.com/api/?name=User+Comum&background=28a745&color=fff'
)
ON CONFLICT (email) DO UPDATE SET
    nome = EXCLUDED.nome,
    avatar = EXCLUDED.avatar;

-- Verificar dados
SELECT id, nome, email, perfil, avatar, created_at FROM users;
