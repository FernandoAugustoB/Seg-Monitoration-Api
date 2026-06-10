-- Criar permissões básicas
INSERT INTO permissoes (nome, descricao) VALUES 
('relatorio.view', 'Visualizar relatórios operacionais'),
('relatorio.create', 'Criar novos relatórios'),
('admin.gerenciar', 'Acesso total ao sistema');

-- Criar grupos (Cargos)
INSERT INTO grupos (nome, cor_hex, ordem) VALUES 
('Admin', '#FF0000', 1),
('Membro', '#00FF00', 2);

-- Associar permissões ao grupo Admin (Tudo)
INSERT INTO grupo_permissoes (grupo_id, permissao_id)
SELECT (SELECT id FROM grupos WHERE nome = 'Admin'), id FROM permissoes;

-- Associar apenas visualização ao grupo Membro
INSERT INTO grupo_permissoes (grupo_id, permissao_id)
SELECT (SELECT id FROM grupos WHERE nome = 'Membro'), id FROM permissoes WHERE nome = 'relatorio.view';