create table if not exists users (
    id uuid not null,
    name varchar(60),
    firstName varchar(60),
    lastName varchar(60),
    bio varchar(90),
    phoneNumbers text[],
    createdAt timestamp,
    updateAt timestamp,
    deletedAt timestamp,
    status varchar(60),
    PRIMARY KEY(id)
    );

create table adress (
    id uuid,
    users_id uuid,
    country varchar(70),
    city varchar(70),
    district varchar(70),
    postalCodes int,
    PRIMARY KEY(id),
    CONSTRAINT fk_users
    FOREIGN KEY(users_id) 
    REFERENCES users(id)
    ON DELETE CASCADE 
); 




