---- duckdb ----
create schema com;
use schema com;

create sequence cmd_seq;

create table cmd(id integer default nextval('cmd_seq'), command varchar, descr varchar);

---- sqlite ----

create table cmd(id integer primary key autoincrement, command text unique, descr text);
