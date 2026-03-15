CREATE TABLE if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) unique not null,
    age int,
    gender varchar(50),
    birth_date date,
    is_active boolean default true,
    created_at timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CHECK (user_id != friend_id)
);

-- Заполняем 20 пользователей
INSERT INTO users (name, email, age, gender, birth_date) VALUES
('Alice', 'alice@test.com', 20, 'Female', '2006-01-01'),
('Bob', 'bob@test.com', 21, 'Male', '2005-02-02'),
('Charlie', 'charlie@test.com', 22, 'Male', '2004-03-03'),
('David', 'david@test.com', 23, 'Male', '2003-04-04'),
('Eva', 'eva@test.com', 20, 'Female', '2006-05-05'),
('Frank', 'frank@test.com', 21, 'Male', '2005-06-06'),
('Grace', 'grace@test.com', 22, 'Female', '2004-07-07'),
('Hannah', 'hannah@test.com', 23, 'Female', '2003-08-08'),
('Ivan', 'ivan@test.com', 20, 'Male', '2006-09-09'),
('Jack', 'jack@test.com', 21, 'Male', '2005-10-10'),
('Karen', 'karen@test.com', 22, 'Female', '2004-11-11'),
('Leo', 'leo@test.com', 23, 'Male', '2003-12-12'),
('Mona', 'mona@test.com', 20, 'Female', '2006-01-13'),
('Nina', 'nina@test.com', 21, 'Female', '2005-02-14'),
('Oscar', 'oscar@test.com', 22, 'Male', '2004-03-15'),
('Paul', 'paul@test.com', 23, 'Male', '2003-04-16'),
('Quinn', 'quinn@test.com', 20, 'Female', '2006-05-17'),
('Rose', 'rose@test.com', 21, 'Female', '2005-06-18'),
('Sam', 'sam@test.com', 22, 'Male', '2004-07-19'),
('Tom', 'tom@test.com', 23, 'Male', '2003-08-20');

INSERT INTO user_friends (user_id, friend_id) VALUES
(1, 3), (3, 1),
(1, 4), (4, 1),
(1, 5), (5, 1),
(2, 3), (3, 2),
(2, 4), (4, 2),
(2, 5), (5, 2);