-- Down migration para reverter a criação do schema de usuários

DROP TABLE IF EXISTS usuario_grupos CASCADE;
DROP TABLE IF EXISTS grupo_permissoes CASCADE;
DROP TABLE IF EXISTS permissoes CASCADE;
DROP TABLE IF EXISTS grupos CASCADE;
DROP TABLE IF EXISTS usuarios CASCADE;