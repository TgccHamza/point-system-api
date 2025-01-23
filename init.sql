CREATE USER 'melkey'@'%' IDENTIFIED BY 'password1234';
GRANT ALL PRIVILEGES ON blueprint.* TO 'melkey'@'%';

-- CREATE USER 'melkey'@'::1' IDENTIFIED BY 'password1234';
-- GRANT ALL PRIVILEGES ON blueprint.* TO 'melkey'@'::1';

FLUSH PRIVILEGES;
