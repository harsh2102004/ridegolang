--
-- PostgreSQL database dump
--

\restrict J4JBrFN2FcuVl1lcKL3LbKfabyyr4Bg9VIQgp2UREEx4VkycMTmHUDhwmWbgsfa

-- Dumped from database version 16.14 (Debian 16.14-1.pgdg13+1)
-- Dumped by pg_dump version 16.14 (Debian 16.14-1.pgdg13+1)

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

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: connected_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.connected_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    provider text NOT NULL,
    provider_user_id text,
    email text,
    access_token text NOT NULL,
    refresh_token text,
    expires_at timestamp without time zone,
    scopes text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    plan text DEFAULT 'free'::text,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: webhook_triggers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.webhook_triggers (
    id uuid NOT NULL,
    workflow_id uuid NOT NULL,
    user_id uuid NOT NULL,
    webhook_url_id text NOT NULL,
    frontend_node_id text NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: workflow_edges; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workflow_edges (
    id uuid NOT NULL,
    workflow_id uuid NOT NULL,
    source_step_id uuid NOT NULL,
    target_step_id uuid NOT NULL,
    condition_branch text,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: workflow_runs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workflow_runs (
    id uuid NOT NULL,
    workflow_id uuid NOT NULL,
    user_id uuid NOT NULL,
    status text DEFAULT 'pending'::text NOT NULL,
    started_at timestamp without time zone,
    finished_at timestamp without time zone,
    error_message text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: workflow_step_runs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workflow_step_runs (
    id uuid NOT NULL,
    workflow_run_id uuid NOT NULL,
    workflow_step_id uuid NOT NULL,
    status text DEFAULT 'pending'::text NOT NULL,
    input jsonb,
    output jsonb,
    error_message text,
    started_at timestamp without time zone,
    finished_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: workflow_steps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workflow_steps (
    id uuid NOT NULL,
    workflow_id uuid NOT NULL,
    frontend_node_id text NOT NULL,
    step_order integer NOT NULL,
    step_type text NOT NULL,
    config jsonb NOT NULL,
    position_x double precision DEFAULT 0 NOT NULL,
    position_y double precision DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: workflow_usage; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workflow_usage (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    month text NOT NULL,
    workflow_runs integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: workflows; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workflows (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    trigger_type text DEFAULT 'manual'::text NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    schedule_enabled boolean DEFAULT false NOT NULL,
    schedule_type text,
    schedule_value text,
    next_run_at timestamp without time zone,
    last_run_at timestamp without time zone,
    is_schedule_running boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Data for Name: connected_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.connected_accounts (id, user_id, provider, provider_user_id, email, access_token, refresh_token, expires_at, scopes, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, email, password, is_verified, plan, created_at) FROM stdin;
3e609664-8b7f-4b24-86a1-6ee701a4386c	shinasmuhammed.dev@gmail.com	$2a$12$ZW.OrTStIA5X2yJoetm11u0svmL1NVM/QL.P42RHAojPVO/ITKbKq	f	free	2026-05-20 04:58:11.467667
0325fd5c-9517-4a0a-93f5-d503c385e003	shnsmk17@gmail.com	$2a$12$KXVyGX1rt40cToe1xVcaL.LeV7TdeA46KRs5G3a2cHfnvjpcTvE9q	t	free	2026-05-20 10:38:00.193257
9af6d792-0ba6-4907-87a8-f7b5b1313648	shinascontact@gmail.com	$2a$10$HDgHDh.VDUoYUTu9TVuBEOO4xgTZHQWb34Oph8/QPva1TIKYQjl2u	t	free	2026-05-20 05:11:53.241912
6380b958-e1fe-4da2-a091-e0ba526f5d46	yadhukrishnangta@gmail.com	$2a$12$4N9KGcy5USktOJy56LhgSuGA5k7w3EXuqUwr7B4AxSzlWQSnPfeOS	f	free	2026-05-25 04:28:15.737645
\.


--
-- Data for Name: webhook_triggers; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.webhook_triggers (id, workflow_id, user_id, webhook_url_id, frontend_node_id, created_at) FROM stdin;
270cce06-88da-4add-aa2c-5f3f446ccd04	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	4fb20d6c-b8e9-4be8-bbbb-89abab339320	node-1779346443005	2026-05-27 08:21:37.209774
b35714c9-bb3e-4240-b44b-f898c02fbdc8	85def2d8-5d22-482c-8880-c8c5a1584147	3e609664-8b7f-4b24-86a1-6ee701a4386c	ea4c80d6-e27e-40f0-aec7-6350771c69a4	node-1779871871529	2026-05-27 08:51:26.057208
8e930762-bead-4a61-baa1-dfead292874a	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	78c2290e-2399-4e2b-9285-249cf18c059b	node-1779512204825	2026-05-30 12:32:52.603389
\.


--
-- Data for Name: workflow_edges; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.workflow_edges (id, workflow_id, source_step_id, target_step_id, condition_branch, created_at) FROM stdin;
096871d3-c36f-4ce3-956c-b509b72f1b07	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	798e4e7f-7440-450c-9883-3dd9d7f2a003	c20d3dfa-33b3-46c4-a18a-3a68c6fe05e7	\N	2026-05-27 08:21:37.225384
4ac89b61-eecc-4bc1-ac87-90e66445dc77	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	c20d3dfa-33b3-46c4-a18a-3a68c6fe05e7	34dbe629-7f1e-461a-9a4e-6e7cb63006da	true	2026-05-27 08:21:37.226943
493719c1-9d9b-4a87-a272-b86d1c9353b6	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	c20d3dfa-33b3-46c4-a18a-3a68c6fe05e7	3bf57504-48c3-4ea7-9ad5-a2e6b893e1aa	false	2026-05-27 08:21:37.228015
7b66ba4b-e48f-408f-9e12-d20060af0b6b	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	34dbe629-7f1e-461a-9a4e-6e7cb63006da	a3c14370-f205-4c24-85d1-03cc45fa61b6	\N	2026-05-27 08:21:37.228928
0c00ecca-7824-4ad0-afcd-06eeddcfe7e3	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	3bf57504-48c3-4ea7-9ad5-a2e6b893e1aa	a3c14370-f205-4c24-85d1-03cc45fa61b6	\N	2026-05-27 08:21:37.229773
d13f198b-ba69-4d63-bc4f-7960ca60ce17	85def2d8-5d22-482c-8880-c8c5a1584147	e7cf727a-c219-45a9-ac0b-3871aadad7ff	0581c397-39a2-4954-ac5d-78734cdda1ef	\N	2026-05-27 08:51:26.069902
24d560dd-7a7d-40db-a91b-4c1b668b0b85	fbe2c588-83d8-4c30-a742-bc0d9df062ff	406337ae-df59-4ae3-bec7-b51be2302019	e11afa6f-a297-4033-b7a1-8b694870fc5e	\N	2026-05-30 12:32:52.629925
a83b2281-7282-474e-9e53-7d2f33ecd880	fbe2c588-83d8-4c30-a742-bc0d9df062ff	406337ae-df59-4ae3-bec7-b51be2302019	7cb8c287-2b7d-45cb-b08c-a6aa952eaa0c	\N	2026-05-30 12:32:52.632288
\.


--
-- Data for Name: workflow_runs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.workflow_runs (id, workflow_id, user_id, status, started_at, finished_at, error_message, created_at) FROM stdin;
c30ba0e7-2eb3-4beb-8c8d-2e921495a07f	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-30 12:31:54.963995	2026-05-30 12:31:59.175739	\N	2026-05-30 12:31:54.963995
4d352985-566c-462a-8bd9-74b0844c3e4b	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:22:26.526739	2026-05-26 11:22:40.576351	\N	2026-05-26 11:22:26.526739
3d189f58-b0c5-4db3-bbb7-852b62dc26a0	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	failed	2026-05-21 06:53:58.826851	2026-05-21 06:53:58.835625	workflow has no steps	2026-05-21 06:53:58.826851
036eb6f6-c534-4efa-bff9-01338ffbe389	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-21 07:29:07.637743	2026-05-21 07:29:07.678427	\N	2026-05-21 07:29:07.637743
fc6b1a37-9283-4f9f-886f-0650103cd257	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:24:09.123282	2026-05-26 11:24:23.151282	\N	2026-05-26 11:24:09.123282
8640de5b-6cc8-496a-8297-76a27b5a675a	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	failed	2026-05-21 10:39:28.015686	2026-05-21 10:39:28.073703	condition field is required	2026-05-21 10:39:28.015686
176f95c5-7f11-4a49-92c7-f69b229a257c	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	failed	2026-05-21 10:41:54.10184	2026-05-21 10:41:54.117763	condition field is required	2026-05-21 10:41:54.10184
d94e28a9-0293-49e5-9585-dccff1d13077	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:50:54.436286	2026-05-27 04:51:08.699355	\N	2026-05-27 04:50:54.436286
0a9ee0b2-a7d6-41c6-aaab-f742f2c87e49	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-21 10:50:39.856341	2026-05-21 10:50:53.655394	\N	2026-05-21 10:50:39.856341
6ff2b376-e670-4e8f-a74f-11aef11e1913	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:29:18.589918	2026-05-26 11:29:32.748346	\N	2026-05-26 11:29:18.589918
560cb446-ed8e-4d1e-93e3-85b149ac8aa1	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 04:47:08.161169	2026-05-23 04:47:22.307058	\N	2026-05-23 04:47:08.161169
164d50e0-9b79-478a-ae68-f9d16d201fe2	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 04:47:14.515514	2026-05-23 04:47:36.579646	\N	2026-05-23 04:47:14.515514
27f9bbba-6376-4de1-88fb-01fba0997570	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 04:49:26.681701	2026-05-23 04:49:42.822072	\N	2026-05-23 04:49:26.681701
6cf45d49-6ad0-4736-83b7-0c971f1aa304	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:33:44.324987	2026-05-26 11:33:58.292577	\N	2026-05-26 11:33:44.324987
acf54fb7-bb09-4df4-ba83-91cb2a6d6421	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 04:58:10.304239	2026-05-23 04:58:14.790659	\N	2026-05-23 04:58:10.304239
5f1e79ec-d5b1-4830-b63c-fa34e67b766b	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 05:57:23.075369	2026-05-23 05:57:37.291186	\N	2026-05-23 05:57:23.075369
17f47bb2-2c36-4f82-9fd9-e0be78120213	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-30 12:32:21.009794	2026-05-30 12:32:24.782269	\N	2026-05-30 12:32:21.009794
fd7585e7-757a-440b-a6ae-9ce43e90f5ad	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 05:58:29.065935	2026-05-23 05:58:42.832749	\N	2026-05-23 05:58:29.065935
d78cc662-c258-41a1-88e3-825cc231b03e	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:06:04.182897	2026-05-27 04:06:18.784415	\N	2026-05-27 04:06:04.182897
b2434eee-864b-4ff6-ad60-c42281a7e2b3	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 05:59:56.383627	2026-05-23 06:00:10.035189	\N	2026-05-23 05:59:56.383627
cdf3a8ba-742b-454f-b570-d524faa94750	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-23 07:00:09.32176	2026-05-23 07:00:13.138048	\N	2026-05-23 07:00:09.32176
210b6a03-3b1d-4bac-a6fe-66d4fdf1f0d8	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:52:35.226295	2026-05-27 04:52:48.666566	\N	2026-05-27 04:52:35.226295
58081121-ff1b-4ae5-a329-88f5de6c153b	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-25 04:04:01.33276	2026-05-25 04:04:05.601632	\N	2026-05-25 04:04:01.33276
a6ea53ff-0798-4f4b-b633-02ebea6ed401	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:17:24.270353	2026-05-27 04:17:39.1571	\N	2026-05-27 04:17:24.270353
b7c8883e-0b8e-4370-82e4-9d01cbea6988	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-25 04:53:47.063801	2026-05-25 04:53:51.16992	\N	2026-05-25 04:53:47.063801
06bd87bf-e2e0-4a9b-bc98-1ae50b08f6e2	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 10:58:46.96196	2026-05-26 10:59:00.874792	\N	2026-05-26 10:58:46.96196
b8e59ad8-2d2c-41db-ad04-86febe8d9191	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:14:49.714998	2026-05-26 11:15:04.289601	\N	2026-05-26 11:14:49.714998
0cbbb45b-5211-4741-97f1-7096fea3dc91	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:20:13.42012	2026-05-27 04:20:27.355876	\N	2026-05-27 04:20:13.42012
58aa9df1-bc74-4902-a684-5c569f856ba0	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:15:19.476062	2026-05-26 11:15:33.545759	\N	2026-05-26 11:15:19.476062
bccb77fa-7cb3-4b12-b65e-4648b2fb0f28	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:15:34.655095	2026-05-26 11:15:48.745648	\N	2026-05-26 11:15:34.655095
97b708cd-645d-4747-a5e1-075106b1e3a5	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:20:32.00293	2026-05-26 11:20:46.832214	\N	2026-05-26 11:20:32.00293
0e835fe8-6965-4a17-9546-24f59950e5cc	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:30:48.269232	2026-05-27 04:31:03.144571	\N	2026-05-27 04:30:48.269232
cb68a4b4-f819-470a-898b-b4501e633263	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-26 11:21:21.643944	2026-05-26 11:21:35.940477	\N	2026-05-26 11:21:21.643944
d6d57c9c-e7ef-4a94-86a9-7ebdccaf0e66	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 05:37:21.240303	2026-05-27 05:37:25.288523	\N	2026-05-27 05:37:21.240303
b4783b07-62d7-4e9a-a826-81ee96e36607	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:36:01.268649	2026-05-27 04:36:14.954835	\N	2026-05-27 04:36:01.268649
57bf664e-dd5d-4896-b28e-02675532cb3e	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:37:34.898702	2026-05-27 04:37:48.313983	\N	2026-05-27 04:37:34.898702
821b65e8-f7b4-49eb-a864-05f2c3b43895	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:40:14.07908	2026-05-27 04:40:27.375458	\N	2026-05-27 04:40:14.07908
3d3c7076-8b53-42e8-be75-e8f01dd4f714	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 05:54:49.090256	2026-05-27 05:55:02.682976	\N	2026-05-27 05:54:49.090256
2f41babe-5818-4a72-b1e1-9236522e940b	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:41:57.702546	2026-05-27 04:42:11.752075	\N	2026-05-27 04:41:57.702546
2ecb583b-e131-44be-ae06-8e0ce7030626	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:46:06.753274	2026-05-27 04:46:20.596046	\N	2026-05-27 04:46:06.753274
64afcac3-b83b-480d-8cec-3e16130ee22b	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-30 12:32:37.909037	2026-05-30 12:32:41.423448	\N	2026-05-30 12:32:37.909037
081f7cd6-c142-47cf-85ee-af48431561ca	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:48:08.607623	2026-05-27 04:48:22.278869	\N	2026-05-27 04:48:08.607623
5c5fedc7-0845-4903-9813-5f9fea1f2ecb	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 05:58:18.380669	2026-05-27 05:58:32.111402	\N	2026-05-27 05:58:18.380669
4cf2aeca-7302-476b-90f2-a78f98d7b57e	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 04:49:56.873931	2026-05-27 04:50:10.369454	\N	2026-05-27 04:49:56.873931
3bed09e2-45ba-4944-a303-64bcbfba42e2	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 06:03:20.371991	2026-05-27 06:03:34.08602	\N	2026-05-27 06:03:20.371991
d6cf002d-32e6-420e-b71f-79c6c4bd16b2	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-27 06:22:07.277703	2026-05-27 06:22:21.201133	\N	2026-05-27 06:22:07.277703
e6140c66-ce24-4fbc-bbb4-63e9829b725c	fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	failed	2026-05-30 12:32:58.917507	2026-05-30 12:33:02.776003	http request url is required	2026-05-30 12:32:58.917507
640ba18c-ec11-48b7-b95b-931452dc00da	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	success	2026-05-30 12:09:34.574741	2026-05-30 12:09:50.458443	\N	2026-05-30 12:09:34.574741
\.


--
-- Data for Name: workflow_step_runs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.workflow_step_runs (id, workflow_run_id, workflow_step_id, status, input, output, error_message, started_at, finished_at, created_at) FROM stdin;
b8f80b93-c879-4bd0-bd21-8239307db695	e6140c66-ce24-4fbc-bbb4-63e9829b725c	e11afa6f-a297-4033-b7a1-8b694870fc5e	success	{"webhook_url": "http://localhost:8080/webhooks/78c2290e-2399-4e2b-9285-249cf18c059b", "webhook_url_id": "78c2290e-2399-4e2b-9285-249cf18c059b"}	{"to": "shinascontact@gmail.com", "status": "sent", "subject": "{{trigger.status}}"}	\N	2026-05-30 12:32:58.927451	2026-05-30 12:33:02.726825	2026-05-30 12:32:58.927451
faa4b9af-bd2d-4017-abd9-6c1323d9d7f8	640ba18c-ec11-48b7-b95b-931452dc00da	798e4e7f-7440-450c-9883-3dd9d7f2a003	success	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320"}	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320"}	\N	2026-05-30 12:09:34.636061	2026-05-30 12:09:34.644246	2026-05-30 12:09:34.636061
0aa96925-1c1c-4fd0-94e6-bd719a0010bb	640ba18c-ec11-48b7-b95b-931452dc00da	c20d3dfa-33b3-46c4-a18a-3a68c6fe05e7	success	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320"}	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320", "condition_result": false}	\N	2026-05-30 12:09:34.646673	2026-05-30 12:09:34.649153	2026-05-30 12:09:34.646673
06144556-37ac-4afc-9861-d03122fd7aa8	640ba18c-ec11-48b7-b95b-931452dc00da	3bf57504-48c3-4ea7-9ad5-a2e6b893e1aa	success	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320"}	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320"}	\N	2026-05-30 12:09:34.653641	2026-05-30 12:09:44.657027	2026-05-30 12:09:34.653641
dc12c860-0a1d-44f9-a379-869993b0e2c6	640ba18c-ec11-48b7-b95b-931452dc00da	a3c14370-f205-4c24-85d1-03cc45fa61b6	success	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320"}	{"to": "shinssascoentsact@gmail.com", "status": "sent", "subject": "Subject: Workflow Completed"}	\N	2026-05-30 12:09:44.659331	2026-05-30 12:09:50.45558	2026-05-30 12:09:44.659331
651c2bef-45d4-414c-908e-ca812381793d	e6140c66-ce24-4fbc-bbb4-63e9829b725c	7cb8c287-2b7d-45cb-b08c-a6aa952eaa0c	failed	{"webhook_url": "http://localhost:8080/webhooks/78c2290e-2399-4e2b-9285-249cf18c059b", "webhook_url_id": "78c2290e-2399-4e2b-9285-249cf18c059b"}	\N	http request url is required	2026-05-30 12:33:02.731475	2026-05-30 12:33:02.733508	2026-05-30 12:33:02.731475
71b4dcd0-eb0a-4884-ae07-574d46046c28	e6140c66-ce24-4fbc-bbb4-63e9829b725c	406337ae-df59-4ae3-bec7-b51be2302019	success	{"webhook_url": "http://localhost:8080/webhooks/78c2290e-2399-4e2b-9285-249cf18c059b", "webhook_url_id": "78c2290e-2399-4e2b-9285-249cf18c059b"}	{"webhook_url": "http://localhost:8080/webhooks/78c2290e-2399-4e2b-9285-249cf18c059b", "webhook_url_id": "78c2290e-2399-4e2b-9285-249cf18c059b"}	\N	2026-05-30 12:32:58.923781	2026-05-30 12:32:58.925049	2026-05-30 12:32:58.923781
\.


--
-- Data for Name: workflow_steps; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.workflow_steps (id, workflow_id, frontend_node_id, step_order, step_type, config, position_x, position_y, created_at) FROM stdin;
798e4e7f-7440-450c-9883-3dd9d7f2a003	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	node-1779346443005	1	webhook_trigger	{"webhook_url": "http://localhost:8080/webhooks/4fb20d6c-b8e9-4be8-bbbb-89abab339320", "webhook_url_id": "4fb20d6c-b8e9-4be8-bbbb-89abab339320"}	60	345.5	2026-05-27 08:21:37.216003
c20d3dfa-33b3-46c4-a18a-3a68c6fe05e7	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	node-1779359849206	2	condition	{"field": "status", "value": "succes", "operator": "equals"}	429.2300955063188	343.9076530850449	2026-05-27 08:21:37.219193
a3c14370-f205-4c24-85d1-03cc45fa61b6	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	node-1779359864342	3	email	{"to": "shinssascoentsact@gmail.com", "body": "Workflow finished.\\nStatus: {{trigger.status}}\\nMessage: {{trigger.message}}", "subject": "Subject: Workflow Completed"}	1215.2300955063188	347.9076530850449	2026-05-27 08:21:37.220864
34dbe629-7f1e-461a-9a4e-6e7cb63006da	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	node-1779359870136	4	http_request	{"url": "https://jsonplaceholder.typicode.com/todos/1", "body": "", "retry": {"enabled": false, "max_attempts": 3, "delay_seconds": 2}, "method": "GET", "timeout_seconds": 15}	801.2300955063188	197.9076530850449	2026-05-27 08:21:37.222786
3bf57504-48c3-4ea7-9ad5-a2e6b893e1aa	5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	node-1779359877138	5	delay	{"unit": "second", "duration": 10}	807.2300955063188	529.9076530850449	2026-05-27 08:21:37.224197
e7cf727a-c219-45a9-ac0b-3871aadad7ff	85def2d8-5d22-482c-8880-c8c5a1584147	node-1779871871529	1	webhook_trigger	{"webhook_url": "http://localhost:8080/webhooks/ea4c80d6-e27e-40f0-aec7-6350771c69a4", "webhook_url_id": "ea4c80d6-e27e-40f0-aec7-6350771c69a4"}	476.6171480502669	171.2416526400949	2026-05-27 08:51:26.061564
0581c397-39a2-4954-ac5d-78734cdda1ef	85def2d8-5d22-482c-8880-c8c5a1584147	node-1779871877598	2	delay	{"unit": "second", "duration": 5}	878.7787203494961	175.54371844021347	2026-05-27 08:51:26.067866
406337ae-df59-4ae3-bec7-b51be2302019	fbe2c588-83d8-4c30-a742-bc0d9df062ff	node-1779512204825	1	webhook_trigger	{"webhook_url": "http://localhost:8080/webhooks/78c2290e-2399-4e2b-9285-249cf18c059b", "webhook_url_id": "78c2290e-2399-4e2b-9285-249cf18c059b"}	77.2608078704134	242.5213976989214	2026-05-30 12:32:52.608569
e11afa6f-a297-4033-b7a1-8b694870fc5e	fbe2c588-83d8-4c30-a742-bc0d9df062ff	node-1779512207789	2	email	{"to": "shinascontact@gmail.com", "body": "{{trigger.message}}", "subject": "{{trigger.status}}"}	535.6410333463784	25.64951166058006	2026-05-30 12:32:52.618127
7cb8c287-2b7d-45cb-b08c-a6aa952eaa0c	fbe2c588-83d8-4c30-a742-bc0d9df062ff	node-1780144337607	3	http_request	{"url": "", "body": "", "retry": {"enabled": false, "max_attempts": 3, "delay_seconds": 2}, "method": "GET", "timeout_seconds": 15}	756.0787258970279	274.63596745561995	2026-05-30 12:32:52.621199
\.


--
-- Data for Name: workflow_usage; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.workflow_usage (id, user_id, month, workflow_runs, created_at, updated_at) FROM stdin;
52e23d87-2aac-45ab-ac88-42818acffbcd	9af6d792-0ba6-4907-87a8-f7b5b1313648	2026-05	33	2026-05-20 10:31:47.698408	2026-05-30 12:32:58.92199
\.


--
-- Data for Name: workflows; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.workflows (id, user_id, name, description, trigger_type, is_active, schedule_enabled, schedule_type, schedule_value, next_run_at, last_run_at, is_schedule_running, created_at, updated_at) FROM stdin;
5fcf632b-0bf7-4f3e-be52-ef65b85f1d89	9af6d792-0ba6-4907-87a8-f7b5b1313648	send email	Add description...	manual	t	f	\N	\N	\N	\N	f	2026-05-21 06:53:56.831933	2026-05-27 07:26:28.333903
85def2d8-5d22-482c-8880-c8c5a1584147	3e609664-8b7f-4b24-86a1-6ee701a4386c	test		manual	t	f	\N	\N	\N	\N	f	2026-05-27 08:51:09.454839	2026-05-27 08:51:09.454839
fbe2c588-83d8-4c30-a742-bc0d9df062ff	9af6d792-0ba6-4907-87a8-f7b5b1313648	test	Add description...	manual	t	f	\N	\N	\N	\N	f	2026-05-23 04:56:42.353863	2026-05-27 09:12:51.716211
38317513-f9a6-4006-885e-f7269c6e59cd	9af6d792-0ba6-4907-87a8-f7b5b1313648	heyy	heylo	manual	t	f	\N	\N	\N	\N	f	2026-05-30 09:41:52.733911	2026-05-30 09:41:52.733911
\.


--
-- Name: connected_accounts connected_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT connected_accounts_pkey PRIMARY KEY (id);


--
-- Name: connected_accounts connected_accounts_user_id_provider_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT connected_accounts_user_id_provider_key UNIQUE (user_id, provider);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: webhook_triggers webhook_triggers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_triggers
    ADD CONSTRAINT webhook_triggers_pkey PRIMARY KEY (id);


--
-- Name: webhook_triggers webhook_triggers_webhook_url_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_triggers
    ADD CONSTRAINT webhook_triggers_webhook_url_id_key UNIQUE (webhook_url_id);


--
-- Name: workflow_edges workflow_edges_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_edges
    ADD CONSTRAINT workflow_edges_pkey PRIMARY KEY (id);


--
-- Name: workflow_runs workflow_runs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_runs
    ADD CONSTRAINT workflow_runs_pkey PRIMARY KEY (id);


--
-- Name: workflow_step_runs workflow_step_runs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_step_runs
    ADD CONSTRAINT workflow_step_runs_pkey PRIMARY KEY (id);


--
-- Name: workflow_steps workflow_steps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_steps
    ADD CONSTRAINT workflow_steps_pkey PRIMARY KEY (id);


--
-- Name: workflow_usage workflow_usage_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_usage
    ADD CONSTRAINT workflow_usage_pkey PRIMARY KEY (id);


--
-- Name: workflow_usage workflow_usage_user_id_month_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_usage
    ADD CONSTRAINT workflow_usage_user_id_month_key UNIQUE (user_id, month);


--
-- Name: workflows workflows_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflows
    ADD CONSTRAINT workflows_pkey PRIMARY KEY (id);


--
-- Name: connected_accounts connected_accounts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.connected_accounts
    ADD CONSTRAINT connected_accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: webhook_triggers webhook_triggers_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_triggers
    ADD CONSTRAINT webhook_triggers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: webhook_triggers webhook_triggers_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_triggers
    ADD CONSTRAINT webhook_triggers_workflow_id_fkey FOREIGN KEY (workflow_id) REFERENCES public.workflows(id) ON DELETE CASCADE;


--
-- Name: workflow_edges workflow_edges_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_edges
    ADD CONSTRAINT workflow_edges_workflow_id_fkey FOREIGN KEY (workflow_id) REFERENCES public.workflows(id) ON DELETE CASCADE;


--
-- Name: workflow_runs workflow_runs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_runs
    ADD CONSTRAINT workflow_runs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: workflow_runs workflow_runs_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_runs
    ADD CONSTRAINT workflow_runs_workflow_id_fkey FOREIGN KEY (workflow_id) REFERENCES public.workflows(id) ON DELETE CASCADE;


--
-- Name: workflow_step_runs workflow_step_runs_workflow_run_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_step_runs
    ADD CONSTRAINT workflow_step_runs_workflow_run_id_fkey FOREIGN KEY (workflow_run_id) REFERENCES public.workflow_runs(id) ON DELETE CASCADE;


--
-- Name: workflow_step_runs workflow_step_runs_workflow_step_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_step_runs
    ADD CONSTRAINT workflow_step_runs_workflow_step_id_fkey FOREIGN KEY (workflow_step_id) REFERENCES public.workflow_steps(id) ON DELETE CASCADE;


--
-- Name: workflow_steps workflow_steps_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_steps
    ADD CONSTRAINT workflow_steps_workflow_id_fkey FOREIGN KEY (workflow_id) REFERENCES public.workflows(id) ON DELETE CASCADE;


--
-- Name: workflow_usage workflow_usage_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflow_usage
    ADD CONSTRAINT workflow_usage_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: workflows workflows_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workflows
    ADD CONSTRAINT workflows_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict J4JBrFN2FcuVl1lcKL3LbKfabyyr4Bg9VIQgp2UREEx4VkycMTmHUDhwmWbgsfa

