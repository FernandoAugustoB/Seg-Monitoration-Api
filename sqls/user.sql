-- Habilitar extensão para UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Tabela de Usuários
CREATE TABLE usuarios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tabela de Grupos (Cargos/Roles)
CREATE TABLE grupos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(50) UNIQUE NOT NULL,
    ordem INTEGER DEFAULT 0
);

-- 3. Tabela de Permissões Específicas
-- Exemplos de nome: "relatorio.criar", "relatorio.deletar", "usuario.gerenciar"
CREATE TABLE permissoes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(100) UNIQUE NOT NULL,
    descricao TEXT
);

-- 4. Relação Grupo <-> Permissões (Muitos para Muitos)
CREATE TABLE grupo_permissoes (
    grupo_id UUID REFERENCES grupos(id) ON DELETE CASCADE,
    permissao_id UUID REFERENCES permissoes(id) ON DELETE CASCADE,
    PRIMARY KEY (grupo_id, permissao_id)
);

-- 5. Relação Usuário <-> Grupos (Muitos para Muitos)
CREATE TABLE usuario_grupos (
    usuario_id UUID REFERENCES usuarios(id) ON DELETE CASCADE,
    grupo_id UUID REFERENCES grupos(id) ON DELETE CASCADE,
    PRIMARY KEY (usuario_id, grupo_id)
);