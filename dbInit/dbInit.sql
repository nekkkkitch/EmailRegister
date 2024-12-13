SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;


create schema if not exists public;

alter schema public owner to pg_database_owner;

SET search_path TO public;

create table if not exists public.users(
    id serial not null primary key,
    email text not null unique,
    password text not null,
    verificated boolean not null default false
);