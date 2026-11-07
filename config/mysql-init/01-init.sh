#!/bin/bash
mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<-EOF
GRANT CREATE, DROP, ALTER, INDEX, INSERT, UPDATE, DELETE ON debut.* TO 'webserver'@'%';
FLUSH PRIVILEGES;
EOF

mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<-EOF
  CREATE DATABASE test_the_model CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<-EOF
CREATE USER IF NOT EXISTS 'testserver'@'%'
IDENTIFIED BY 'iamtest123';
GRANT CREATE, DROP, ALTER, INDEX, INSERT, UPDATE, DELETE ON test_the_model.* TO 'testserver'@'%';
FLUSH PRIVILEGES;
EOF
