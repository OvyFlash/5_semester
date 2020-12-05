CREATE TABLE IF NOT EXISTS Users (
    userid BIGINT NOT NULL AUTO_INCREMENT,
    username text,
    firstname text,
    lastname text,
    email varchar(512) not null,
    encrypted_password text not null,
    phone_number int,
    UNIQUE(email),
    PRIMARY KEY(userid)
);
CREATE TABLE IF NOT EXISTS Followers (
    id bigint not null AUTO_INCREMENT,
    userid bigint not null,
    followed_userid bigint not null,
    PRIMARY KEY(id)
);
CREATE TABLE IF NOT EXISTS Comments (
    id bigint not null AUTO_INCREMENT,
    postid bigint not null,
    userid bigint not null,
    commentary text not null,
    PRIMARY KEY(id)
);
CREATE TABLE IF NOT EXISTS Posts (
    postid bigint not null AUTO_INCREMENT,
    routeid bigint not null,
    PRIMARY KEY(postid)
);
CREATE TABLE IF NOT EXISTS Routes (
    routeid bigint not null AUTO_INCREMENT,
    userid bigint not null,
    PRIMARY KEY(routeid)
);
CREATE TABLE IF NOT EXISTS Points (
    id bigint not null AUTO_INCREMENT,
    routeid bigint not null,
    latitude bigint not null,
    longitude bigint not null,
    point_index smallint not null,
    PRIMARY KEY(id)
);
CREATE TABLE IF NOT EXISTS RouteStats (
    id bigint not null AUTO_INCREMENT,
    routeid bigint not null,
    userid bigint not null,
    workout_time smallint not null,
    average_speed smallint not null,
    burned_fats smallint not null,
    PRIMARY KEY(id)
);