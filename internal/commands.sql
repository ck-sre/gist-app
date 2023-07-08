create database gistapp character set utf8mb4 collate utf8mb4_unicode_ci;
use gistapp;
create table gists ( id integer not null primary key auto_increment, title varchar(100) not null, content text not null, created datetime not null, expires datetime not null) engine=InnoDB;
create index idx_gists_created on gists(created);

insert into gists (title, content, created, expires) values
  ('Beginning', '3 years ago', now(), now() + interval 365 day),
  ('Body', 'The days are I long, the nights I am in bliss', now(), now() + interval 1 day),
  ('End', 'Thank the parents', now(), now() + interval 7 day);

create user 'gistuser'@'localhost' identified by 'gistpass';
grant select, insert, update, delete on gistapp.* to 'gistuser'@'localhost';
alter user 'gistuser'@'localhost' identified by '<redacted>';