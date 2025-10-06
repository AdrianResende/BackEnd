-- Versão Simples do Setup do Banco SmartPicks
-- Execute este arquivo em um banco MySQL vazio

CREATE DATABASE IF NOT EXISTS smartpicks;
USE smartpicks;

-- Criar tabela users completa com avatar
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    data_nascimento DATE NOT NULL,
    perfil ENUM('admin', 'user') NOT NULL DEFAULT 'user',
    avatar TEXT NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   
    INDEX idx_email (email),
    INDEX idx_cpf (cpf),
    INDEX idx_perfil (perfil)
);

-- Inserir usuários padrão para teste
-- Senhas: "123456" (hash bcrypt)
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
ON DUPLICATE KEY UPDATE 
    nome = VALUES(nome),
    avatar = VALUES(avatar);

-- Mostrar estrutura da tabela
DESCRIBE users;

-- Mostrar usuários criados
SELECT id, nome, email, perfil, avatar, created_at FROM users;