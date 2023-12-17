# Schema of my databases

Users  

CREATE TABLE IF NOT EXISTS users ( \
		id INT AUTO_INCREMENT PRIMARY KEY, \
		username VARCHAR(255) UNIQUE, \
		password VARCHAR(255)\
	)

Todos

CREATE TABLE IF NOT EXISTS todos ( \
    id INT AUTO_INCREMENT PRIMARY KEY, \
    username VARCHAR(255), \
    todo_text TEXT, \
    due_date DATE, \
    priority TEXT, \
    FOREIGN KEY (username) REFERENCES users(username) \
);


