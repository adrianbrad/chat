--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: Channels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Channels" (
    "ChannelID" integer NOT NULL,
    "Name" text NOT NULL,
    "Description" text
);


--
-- Name: Channels_ChannelID_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."Channels_ChannelID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Channels_ChannelID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."Channels_ChannelID_seq" OWNED BY public."Channels"."ChannelID";


--
-- Name: Messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Messages" (
    "MessageID" integer NOT NULL,
    "Content" text NOT NULL,
    "RoomID" integer NOT NULL,
    "UserID" integer NOT NULL
);


--
-- Name: Messages_MessageID_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."Messages_MessageID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Messages_MessageID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."Messages_MessageID_seq" OWNED BY public."Messages"."MessageID";


--
-- Name: Permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Permissions" (
    "PermissionID" integer NOT NULL,
    "Name" text,
    "Description" text
);


--
-- Name: Permissions_PermissionID_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."Permissions_PermissionID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Permissions_PermissionID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."Permissions_PermissionID_seq" OWNED BY public."Permissions"."PermissionID";


--
-- Name: Roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Roles" (
    "RoleID" integer NOT NULL,
    "Name" text NOT NULL,
    "Description" text
);


--
-- Name: Roles_Permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Roles_Permissions" (
    "RoleID" integer NOT NULL,
    "PermissionID" integer NOT NULL
);


--
-- Name: Roles_RoleID_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."Roles_RoleID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Roles_RoleID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."Roles_RoleID_seq" OWNED BY public."Roles"."RoleID";


--
-- Name: Rooms; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Rooms" (
    "RoomID" integer NOT NULL,
    "Name" text,
    "Description" text,
    "ChannelID" integer
);


--
-- Name: Rooms_RoomID_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."Rooms_RoomID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Rooms_RoomID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."Rooms_RoomID_seq" OWNED BY public."Rooms"."RoomID";


--
-- Name: Users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Users" (
    "UserID" integer NOT NULL,
    "Name" text NOT NULL,
    "UserData" text,
    "RoleID" integer NOT NULL
);


--
-- Name: Users_Channels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Users_Channels" (
    "UserID" integer NOT NULL,
    "ChannelID" integer NOT NULL,
    "Joined" timestamp without time zone NOT NULL
);


--
-- Name: Users_Rooms; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."Users_Rooms" (
    "UserID" integer NOT NULL,
    "RoomID" integer NOT NULL,
    "Joined" timestamp without time zone NOT NULL,
    "LastActivity" timestamp without time zone
);


--
-- Name: users_UserID_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public."users_UserID_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_UserID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public."users_UserID_seq" OWNED BY public."Users"."UserID";


--
-- Name: Channels ChannelID; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Channels" ALTER COLUMN "ChannelID" SET DEFAULT nextval('public."Channels_ChannelID_seq"'::regclass);


--
-- Name: Messages MessageID; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Messages" ALTER COLUMN "MessageID" SET DEFAULT nextval('public."Messages_MessageID_seq"'::regclass);


--
-- Name: Permissions PermissionID; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Permissions" ALTER COLUMN "PermissionID" SET DEFAULT nextval('public."Permissions_PermissionID_seq"'::regclass);


--
-- Name: Roles RoleID; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Roles" ALTER COLUMN "RoleID" SET DEFAULT nextval('public."Roles_RoleID_seq"'::regclass);


--
-- Name: Rooms RoomID; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Rooms" ALTER COLUMN "RoomID" SET DEFAULT nextval('public."Rooms_RoomID_seq"'::regclass);


--
-- Name: Users UserID; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users" ALTER COLUMN "UserID" SET DEFAULT nextval('public."users_UserID_seq"'::regclass);


--
-- Data for Name: Channels; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Channels" ("ChannelID", "Name", "Description") FROM stdin;
1	mbitcasino	\N
\.


--
-- Data for Name: Messages; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Messages" ("MessageID", "Content", "RoomID", "UserID") FROM stdin;
1	Room2	1	2
2	Room2	1	2
\.


--
-- Data for Name: Permissions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Permissions" ("PermissionID", "Name", "Description") FROM stdin;
1	Test Role	
\.


--
-- Data for Name: Roles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Roles" ("RoleID", "Name", "Description") FROM stdin;
1	Admin	Can do everything
2	Test Role	
\.


--
-- Data for Name: Roles_Permissions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Roles_Permissions" ("RoleID", "PermissionID") FROM stdin;
\.


--
-- Data for Name: Rooms; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Rooms" ("RoomID", "Name", "Description", "ChannelID") FROM stdin;
2	Room2		1
1	Room1		1
\.


--
-- Data for Name: Users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Users" ("UserID", "Name", "UserData", "RoleID") FROM stdin;
1	brad	not	1
2	added	asdf	1
3	bbbb	aaaa	1
6	Nm		1
10	braa		1
11	Jogn		1
13	Yes		1
14	Asdf	stuff	1
\.


--
-- Data for Name: Users_Channels; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Users_Channels" ("UserID", "ChannelID", "Joined") FROM stdin;
3	1	2019-02-26 15:49:54
\.


--
-- Data for Name: Users_Rooms; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public."Users_Rooms" ("UserID", "RoomID", "Joined", "LastActivity") FROM stdin;
\.


--
-- Name: Channels_ChannelID_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Channels_ChannelID_seq"', 1, true);


--
-- Name: Messages_MessageID_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Messages_MessageID_seq"', 2, true);


--
-- Name: Permissions_PermissionID_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Permissions_PermissionID_seq"', 1, true);


--
-- Name: Roles_RoleID_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Roles_RoleID_seq"', 2, true);


--
-- Name: Rooms_RoomID_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."Rooms_RoomID_seq"', 2, true);


--
-- Name: users_UserID_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public."users_UserID_seq"', 35, true);


--
-- Name: Channels Channels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Channels"
    ADD CONSTRAINT "Channels_pkey" PRIMARY KEY ("ChannelID");


--
-- Name: Messages Messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Messages"
    ADD CONSTRAINT "Messages_pkey" PRIMARY KEY ("MessageID");


--
-- Name: Permissions Permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Permissions"
    ADD CONSTRAINT "Permissions_pkey" PRIMARY KEY ("PermissionID");


--
-- Name: Roles_Permissions Role_Permission_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Roles_Permissions"
    ADD CONSTRAINT "Role_Permission_pkey" PRIMARY KEY ("RoleID", "PermissionID");


--
-- Name: Roles Roles_Name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Roles"
    ADD CONSTRAINT "Roles_Name_key" UNIQUE ("Name");


--
-- Name: Roles Roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Roles"
    ADD CONSTRAINT "Roles_pkey" PRIMARY KEY ("RoleID");


--
-- Name: Rooms Rooms_Name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Rooms"
    ADD CONSTRAINT "Rooms_Name_key" UNIQUE ("Name");


--
-- Name: Rooms Rooms_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Rooms"
    ADD CONSTRAINT "Rooms_pkey" PRIMARY KEY ("RoomID");


--
-- Name: Users_Channels User_Channel_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users_Channels"
    ADD CONSTRAINT "User_Channel_pkey" PRIMARY KEY ("UserID", "ChannelID");


--
-- Name: Users_Rooms User_Room_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users_Rooms"
    ADD CONSTRAINT "User_Room_pkey" PRIMARY KEY ("UserID", "RoomID");


--
-- Name: Users Users_Name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "Users_Name_key" UNIQUE ("Name");


--
-- Name: Users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT users_pkey PRIMARY KEY ("UserID");


--
-- Name: fki_FK_Role_User; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "fki_FK_Role_User" ON public."Users" USING btree ("RoleID");


--
-- Name: Users_Channels FK_ChannelID; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users_Channels"
    ADD CONSTRAINT "FK_ChannelID" FOREIGN KEY ("ChannelID") REFERENCES public."Channels"("ChannelID");


--
-- Name: Rooms FK_ChannelID; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Rooms"
    ADD CONSTRAINT "FK_ChannelID" FOREIGN KEY ("ChannelID") REFERENCES public."Channels"("ChannelID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Roles_Permissions FK_Permissions; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Roles_Permissions"
    ADD CONSTRAINT "FK_Permissions" FOREIGN KEY ("PermissionID") REFERENCES public."Permissions"("PermissionID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Users FK_Role_User; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "FK_Role_User" FOREIGN KEY ("RoleID") REFERENCES public."Roles"("RoleID");


--
-- Name: Roles_Permissions FK_Roles; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Roles_Permissions"
    ADD CONSTRAINT "FK_Roles" FOREIGN KEY ("RoleID") REFERENCES public."Roles"("RoleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Users_Rooms FK_RoomID; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users_Rooms"
    ADD CONSTRAINT "FK_RoomID" FOREIGN KEY ("RoomID") REFERENCES public."Rooms"("RoomID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Messages FK_Room_Message; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Messages"
    ADD CONSTRAINT "FK_Room_Message" FOREIGN KEY ("MessageID") REFERENCES public."Rooms"("RoomID");


--
-- Name: Users_Rooms FK_UserID; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users_Rooms"
    ADD CONSTRAINT "FK_UserID" FOREIGN KEY ("UserID") REFERENCES public."Users"("UserID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Users_Channels FK_UserID; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Users_Channels"
    ADD CONSTRAINT "FK_UserID" FOREIGN KEY ("UserID") REFERENCES public."Users"("UserID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Messages FK_User_Message; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."Messages"
    ADD CONSTRAINT "FK_User_Message" FOREIGN KEY ("MessageID") REFERENCES public."Users"("UserID");


--
-- PostgreSQL database dump complete
--

