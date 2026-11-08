CREATE TABLE files (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);
INSERT INTO files (title, content, created, expires) 
VALUES ( 'G.O.A.T Quote', 
    'Talk is cheap. Show me the code. - Linus Torvalds', 
    UTC_TIMESTAMP(), 
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
CREATE INDEX idx_files_created ON files(created);
