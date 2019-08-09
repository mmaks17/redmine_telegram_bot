-- create database bot;
-- create user bot with password 'bot';
-- grant ALL on  DATABASE bot to bot;

-- \c bot


--
-- PostgreSQL database dump
--

-- Dumped from database version 11.4
-- Dumped by pg_dump version 11.4

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

SET default_tablespace = '';

SET default_with_oids = false;


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: maks
--

CREATE TABLE public.tasks (
    userid integer,
    taskid character(254),
    status character(64),
    subject character(254)
);


ALTER TABLE public.tasks OWNER TO bot;

--
-- Name: tasks2; Type: TABLE; Schema: public; Owner: maks
--

CREATE TABLE public.tasks2 (
    userid integer,
    taskid character(254),
    status character(64),
    subject character(254)
);


ALTER TABLE public.tasks2 OWNER TO bot;

--
-- Name: users; Type: TABLE; Schema: public; Owner: bot
--

CREATE TABLE public.users (
    id integer,
    name character(64),
    chat integer
);


ALTER TABLE public.users OWNER TO bot;

--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: maks
--

COPY public.tasks (userid, taskid, status, subject) FROM stdin;
\.


--
-- Data for Name: tasks2; Type: TABLE DATA; Schema: public; Owner: maks
--

COPY public.tasks2 (userid, taskid, status, subject) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: bot
--

COPY public.users (id, name, chat) FROM stdin;
\.


--
-- Name: inq_id; Type: INDEX; Schema: public; Owner: bot
--

CREATE UNIQUE INDEX inq_id ON public.users USING btree (id);


--
-- PostgreSQL database dump complete
--

