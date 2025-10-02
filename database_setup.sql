CREATE DATABASE IF NOT EXISTS smartpicks;
USE smartpicks;

CREATE TABLE IF NOT EXISTS users (
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

-- Inserir usuários padrão (senha para ambos: "123456")
INSERT INTO users (nome, email, password, cpf, data_nascimento, perfil) VALUES
('Admin User', 'admin@smartpicks.com', '$2a$10$YQkxB5K5z5Z5Z5Z5Z5Z5Z.Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5', '00000000000', '1990-01-01', 'admin'),
('User Comum', 'user@smartpicks.com', '$2a$10$YQkxB5K5z5Z5Z5Z5Z5Z5Z.Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5', '11111111111', '1995-06-15', 'user')
ON DUPLICATE KEY UPDATE id=id;

DESCRIBE users;