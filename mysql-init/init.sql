USE flashcards;

CREATE TABLE flashcards_oxford (
    id int primary key auto_increment,
    word varchar(255) unique
);

CREATE TABLE users (
    id int primary key auto_increment,
    username varchar(255)
);

CREATE TABLE flashcards (
    id int primary key auto_increment,
    user_id int,
    word varchar(255),
    foreign key (user_id) references users(id)
);
