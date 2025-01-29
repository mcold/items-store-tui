create schema com;
use schema com;

create sequence cmd_seq;

create table cmd(id integer default nextval('cmd_seq'), command varchar, descr varchar);