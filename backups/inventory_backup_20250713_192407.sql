--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15 (Homebrew)
-- Dumped by pg_dump version 14.15 (Homebrew)

-- Started on 2025-07-13 19:24:09 WIB

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

DROP DATABASE IF EXISTS inventory_system;
--
-- TOC entry 3762 (class 1262 OID 28989)
-- Name: inventory_system; Type: DATABASE; Schema: -; Owner: fmaulana
--

CREATE DATABASE inventory_system WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'C';


ALTER DATABASE inventory_system OWNER TO fmaulana;

\connect inventory_system

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

SET default_table_access_method = heap;

--
-- TOC entry 226 (class 1259 OID 29133)
-- Name: activity_logs; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.activity_logs (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    action text NOT NULL,
    resource text,
    resource_id bigint,
    details text,
    ip_address text,
    user_agent text,
    created_at timestamp with time zone
);


ALTER TABLE public.activity_logs OWNER TO fmaulana;

--
-- TOC entry 225 (class 1259 OID 29132)
-- Name: activity_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.activity_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.activity_logs_id_seq OWNER TO fmaulana;

--
-- TOC entry 3763 (class 0 OID 0)
-- Dependencies: 225
-- Name: activity_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.activity_logs_id_seq OWNED BY public.activity_logs.id;


--
-- TOC entry 212 (class 1259 OID 29005)
-- Name: products; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    name text NOT NULL,
    sku text NOT NULL,
    description text,
    category text,
    price numeric NOT NULL,
    cost numeric NOT NULL,
    quantity bigint DEFAULT 0,
    min_stock bigint DEFAULT 10,
    max_stock bigint DEFAULT 1000,
    location text,
    supplier text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    hpp numeric(10,2) DEFAULT 0
);


ALTER TABLE public.products OWNER TO fmaulana;

--
-- TOC entry 3764 (class 0 OID 0)
-- Dependencies: 212
-- Name: COLUMN products.hpp; Type: COMMENT; Schema: public; Owner: fmaulana
--

COMMENT ON COLUMN public.products.hpp IS 'Harga Pokok Penjualan';


--
-- TOC entry 211 (class 1259 OID 29004)
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.products_id_seq OWNER TO fmaulana;

--
-- TOC entry 3765 (class 0 OID 0)
-- Dependencies: 211
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- TOC entry 224 (class 1259 OID 29113)
-- Name: purchase_order_items; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.purchase_order_items (
    id bigint NOT NULL,
    purchase_order_id bigint NOT NULL,
    product_id bigint NOT NULL,
    quantity_ordered bigint NOT NULL,
    quantity_received bigint DEFAULT 0,
    unit_cost numeric NOT NULL,
    total numeric NOT NULL
);


ALTER TABLE public.purchase_order_items OWNER TO fmaulana;

--
-- TOC entry 223 (class 1259 OID 29112)
-- Name: purchase_order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.purchase_order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.purchase_order_items_id_seq OWNER TO fmaulana;

--
-- TOC entry 3766 (class 0 OID 0)
-- Dependencies: 223
-- Name: purchase_order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.purchase_order_items_id_seq OWNED BY public.purchase_order_items.id;


--
-- TOC entry 222 (class 1259 OID 29090)
-- Name: purchase_orders; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.purchase_orders (
    id bigint NOT NULL,
    po_number text NOT NULL,
    user_id bigint NOT NULL,
    status text DEFAULT 'pending'::text,
    total numeric NOT NULL,
    notes text,
    order_date timestamp with time zone,
    expected_date timestamp with time zone,
    received_date timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    supplier text NOT NULL
);


ALTER TABLE public.purchase_orders OWNER TO fmaulana;

--
-- TOC entry 221 (class 1259 OID 29089)
-- Name: purchase_orders_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.purchase_orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.purchase_orders_id_seq OWNER TO fmaulana;

--
-- TOC entry 3767 (class 0 OID 0)
-- Dependencies: 221
-- Name: purchase_orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.purchase_orders_id_seq OWNED BY public.purchase_orders.id;


--
-- TOC entry 218 (class 1259 OID 29060)
-- Name: sale_items; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.sale_items (
    id bigint NOT NULL,
    sale_id bigint NOT NULL,
    product_id bigint NOT NULL,
    quantity bigint NOT NULL,
    price numeric NOT NULL,
    total numeric NOT NULL
);


ALTER TABLE public.sale_items OWNER TO fmaulana;

--
-- TOC entry 217 (class 1259 OID 29059)
-- Name: sale_items_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.sale_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sale_items_id_seq OWNER TO fmaulana;

--
-- TOC entry 3768 (class 0 OID 0)
-- Dependencies: 217
-- Name: sale_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.sale_items_id_seq OWNED BY public.sale_items.id;


--
-- TOC entry 216 (class 1259 OID 29040)
-- Name: sales; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.sales (
    id bigint NOT NULL,
    sale_number text NOT NULL,
    user_id bigint NOT NULL,
    customer_name text,
    subtotal numeric NOT NULL,
    tax numeric DEFAULT 0,
    discount numeric DEFAULT 0,
    total numeric NOT NULL,
    payment_method text NOT NULL,
    status text DEFAULT 'completed'::text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.sales OWNER TO fmaulana;

--
-- TOC entry 215 (class 1259 OID 29039)
-- Name: sales_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.sales_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sales_id_seq OWNER TO fmaulana;

--
-- TOC entry 3769 (class 0 OID 0)
-- Dependencies: 215
-- Name: sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.sales_id_seq OWNED BY public.sales.id;


--
-- TOC entry 214 (class 1259 OID 29021)
-- Name: stock_movements; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.stock_movements (
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    user_id bigint NOT NULL,
    type text NOT NULL,
    quantity bigint NOT NULL,
    reference text,
    notes text,
    created_at timestamp with time zone
);


ALTER TABLE public.stock_movements OWNER TO fmaulana;

--
-- TOC entry 213 (class 1259 OID 29020)
-- Name: stock_movements_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.stock_movements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stock_movements_id_seq OWNER TO fmaulana;

--
-- TOC entry 3770 (class 0 OID 0)
-- Dependencies: 213
-- Name: stock_movements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.stock_movements_id_seq OWNED BY public.stock_movements.id;


--
-- TOC entry 220 (class 1259 OID 29079)
-- Name: suppliers; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.suppliers (
    id bigint NOT NULL,
    name text NOT NULL,
    email text,
    phone text,
    address text,
    contact_person text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.suppliers OWNER TO fmaulana;

--
-- TOC entry 219 (class 1259 OID 29078)
-- Name: suppliers_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.suppliers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.suppliers_id_seq OWNER TO fmaulana;

--
-- TOC entry 3771 (class 0 OID 0)
-- Dependencies: 219
-- Name: suppliers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.suppliers_id_seq OWNED BY public.suppliers.id;


--
-- TOC entry 210 (class 1259 OID 28991)
-- Name: users; Type: TABLE; Schema: public; Owner: fmaulana
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    name text NOT NULL,
    role text DEFAULT 'employee'::text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO fmaulana;

--
-- TOC entry 209 (class 1259 OID 28990)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: fmaulana
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO fmaulana;

--
-- TOC entry 3772 (class 0 OID 0)
-- Dependencies: 209
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: fmaulana
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 3559 (class 2604 OID 29136)
-- Name: activity_logs id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.activity_logs ALTER COLUMN id SET DEFAULT nextval('public.activity_logs_id_seq'::regclass);


--
-- TOC entry 3541 (class 2604 OID 29008)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 3557 (class 2604 OID 29116)
-- Name: purchase_order_items id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_order_items ALTER COLUMN id SET DEFAULT nextval('public.purchase_order_items_id_seq'::regclass);


--
-- TOC entry 3555 (class 2604 OID 29093)
-- Name: purchase_orders id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_orders ALTER COLUMN id SET DEFAULT nextval('public.purchase_orders_id_seq'::regclass);


--
-- TOC entry 3552 (class 2604 OID 29063)
-- Name: sale_items id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sale_items ALTER COLUMN id SET DEFAULT nextval('public.sale_items_id_seq'::regclass);


--
-- TOC entry 3548 (class 2604 OID 29043)
-- Name: sales id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sales ALTER COLUMN id SET DEFAULT nextval('public.sales_id_seq'::regclass);


--
-- TOC entry 3547 (class 2604 OID 29024)
-- Name: stock_movements id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.stock_movements ALTER COLUMN id SET DEFAULT nextval('public.stock_movements_id_seq'::regclass);


--
-- TOC entry 3553 (class 2604 OID 29082)
-- Name: suppliers id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.suppliers ALTER COLUMN id SET DEFAULT nextval('public.suppliers_id_seq'::regclass);


--
-- TOC entry 3538 (class 2604 OID 28994)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3756 (class 0 OID 29133)
-- Dependencies: 226
-- Data for Name: activity_logs; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.activity_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, created_at) FROM stdin;
\.


--
-- TOC entry 3742 (class 0 OID 29005)
-- Dependencies: 212
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.products (id, name, sku, description, category, price, cost, quantity, min_stock, max_stock, location, supplier, is_active, created_at, updated_at, deleted_at, hpp) FROM stdin;
11	ASD	ASD		ASd	15000	1000	998	10	1000			t	2025-07-13 18:25:43.554685+07	2025-07-13 18:44:43.650835+07	\N	0.00
10	Kabel UTP 1 Meter	KBL-UTP		KABEL	13000	10000	10	10	1000			t	2025-07-13 18:05:48.72817+07	2025-07-13 18:44:43.652421+07	\N	0.00
\.


--
-- TOC entry 3754 (class 0 OID 29113)
-- Dependencies: 224
-- Data for Name: purchase_order_items; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.purchase_order_items (id, purchase_order_id, product_id, quantity_ordered, quantity_received, unit_cost, total) FROM stdin;
9	10	10	15	15	10000	150000
\.


--
-- TOC entry 3752 (class 0 OID 29090)
-- Dependencies: 222
-- Data for Name: purchase_orders; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.purchase_orders (id, po_number, user_id, status, total, notes, order_date, expected_date, received_date, created_at, updated_at, deleted_at, supplier) FROM stdin;
10	PO-20250713-4748	1	pending	150000		2025-07-13 07:00:00+07	2025-07-20 07:00:00+07	\N	2025-07-13 18:05:48.725848+07	2025-07-13 18:05:48.732504+07	\N	PT Surya
\.


--
-- TOC entry 3748 (class 0 OID 29060)
-- Dependencies: 218
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.sale_items (id, sale_id, product_id, quantity, price, total) FROM stdin;
6	4	10	2	13000	26000
7	5	10	1	13000	13000
8	5	11	1	15000	15000
9	6	10	1	13000	13000
10	7	11	1	15000	15000
11	7	10	1	13000	13000
\.


--
-- TOC entry 3746 (class 0 OID 29040)
-- Dependencies: 216
-- Data for Name: sales; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.sales (id, sale_number, user_id, customer_name, subtotal, tax, discount, total, payment_method, status, created_at, updated_at, deleted_at) FROM stdin;
4	SALE-20250713-1752406102	1		26000	0	0	26000	cash	completed	2025-07-13 18:28:22.689943+07	2025-07-13 18:28:22.689943+07	\N
5	SALE-20250713-1752406414	1		28000	0	0	28000	cash	completed	2025-07-13 18:33:34.68686+07	2025-07-13 18:33:34.68686+07	\N
6	SALE-20250713-1752406616	1	mantap	13000	0	0	13000	cash	completed	2025-07-13 18:36:56.688021+07	2025-07-13 18:36:56.688021+07	\N
7	SALE-20250713-1752407083	1		28000	0	0	28000	cash	completed	2025-07-13 18:44:43.653092+07	2025-07-13 18:44:43.653092+07	\N
\.


--
-- TOC entry 3744 (class 0 OID 29021)
-- Dependencies: 214
-- Data for Name: stock_movements; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.stock_movements (id, product_id, user_id, type, quantity, reference, notes, created_at) FROM stdin;
16	10	1	in	15	PO-20250713-4748	Purchase Order PO-20250713-4748 - PT Surya	2025-07-13 18:05:48.72954+07
17	10	1	out	2	SALE-20250713-1752406102	Sale transaction	2025-07-13 18:28:22.689233+07
18	10	1	out	1	SALE-20250713-1752406414	Sale transaction	2025-07-13 18:33:34.680541+07
19	11	1	out	1	SALE-20250713-1752406414	Sale transaction	2025-07-13 18:33:34.686337+07
20	10	1	out	1	SALE-20250713-1752406616	Sale transaction	2025-07-13 18:36:56.6877+07
21	11	1	out	1	SALE-20250713-1752407083	Sale transaction	2025-07-13 18:44:43.651341+07
22	10	1	out	1	SALE-20250713-1752407083	Sale transaction	2025-07-13 18:44:43.652719+07
\.


--
-- TOC entry 3750 (class 0 OID 29079)
-- Dependencies: 220
-- Data for Name: suppliers; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.suppliers (id, name, email, phone, address, contact_person, is_active, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- TOC entry 3740 (class 0 OID 28991)
-- Dependencies: 210
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: fmaulana
--

COPY public.users (id, email, password, name, role, is_active, created_at, updated_at, deleted_at) FROM stdin;
1	admin@inventory.com	$2a$10$634Tys7jT4kv1MQk85mBoOSEnKtP/gOIIpN6qIljJ4vD/X/XUNVBK	System Administrator	admin	t	2025-07-10 21:16:34.073493+07	2025-07-10 21:16:34.073493+07	\N
\.


--
-- TOC entry 3773 (class 0 OID 0)
-- Dependencies: 225
-- Name: activity_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.activity_logs_id_seq', 1, false);


--
-- TOC entry 3774 (class 0 OID 0)
-- Dependencies: 211
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.products_id_seq', 11, true);


--
-- TOC entry 3775 (class 0 OID 0)
-- Dependencies: 223
-- Name: purchase_order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.purchase_order_items_id_seq', 9, true);


--
-- TOC entry 3776 (class 0 OID 0)
-- Dependencies: 221
-- Name: purchase_orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.purchase_orders_id_seq', 10, true);


--
-- TOC entry 3777 (class 0 OID 0)
-- Dependencies: 217
-- Name: sale_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.sale_items_id_seq', 11, true);


--
-- TOC entry 3778 (class 0 OID 0)
-- Dependencies: 215
-- Name: sales_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.sales_id_seq', 7, true);


--
-- TOC entry 3779 (class 0 OID 0)
-- Dependencies: 213
-- Name: stock_movements_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.stock_movements_id_seq', 22, true);


--
-- TOC entry 3780 (class 0 OID 0)
-- Dependencies: 219
-- Name: suppliers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.suppliers_id_seq', 1, false);


--
-- TOC entry 3781 (class 0 OID 0)
-- Dependencies: 209
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: fmaulana
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- TOC entry 3590 (class 2606 OID 29140)
-- Name: activity_logs activity_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.activity_logs
    ADD CONSTRAINT activity_logs_pkey PRIMARY KEY (id);


--
-- TOC entry 3567 (class 2606 OID 29016)
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- TOC entry 3588 (class 2606 OID 29121)
-- Name: purchase_order_items purchase_order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_order_items
    ADD CONSTRAINT purchase_order_items_pkey PRIMARY KEY (id);


--
-- TOC entry 3584 (class 2606 OID 29098)
-- Name: purchase_orders purchase_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_orders
    ADD CONSTRAINT purchase_orders_pkey PRIMARY KEY (id);


--
-- TOC entry 3578 (class 2606 OID 29067)
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (id);


--
-- TOC entry 3574 (class 2606 OID 29050)
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- TOC entry 3571 (class 2606 OID 29028)
-- Name: stock_movements stock_movements_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.stock_movements
    ADD CONSTRAINT stock_movements_pkey PRIMARY KEY (id);


--
-- TOC entry 3581 (class 2606 OID 29087)
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- TOC entry 3569 (class 2606 OID 29018)
-- Name: products uni_products_sku; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT uni_products_sku UNIQUE (sku);


--
-- TOC entry 3586 (class 2606 OID 29100)
-- Name: purchase_orders uni_purchase_orders_po_number; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_orders
    ADD CONSTRAINT uni_purchase_orders_po_number UNIQUE (po_number);


--
-- TOC entry 3576 (class 2606 OID 29052)
-- Name: sales uni_sales_sale_number; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT uni_sales_sale_number UNIQUE (sale_number);


--
-- TOC entry 3562 (class 2606 OID 29002)
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- TOC entry 3564 (class 2606 OID 29000)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3565 (class 1259 OID 29019)
-- Name: idx_products_deleted_at; Type: INDEX; Schema: public; Owner: fmaulana
--

CREATE INDEX idx_products_deleted_at ON public.products USING btree (deleted_at);


--
-- TOC entry 3582 (class 1259 OID 29111)
-- Name: idx_purchase_orders_deleted_at; Type: INDEX; Schema: public; Owner: fmaulana
--

CREATE INDEX idx_purchase_orders_deleted_at ON public.purchase_orders USING btree (deleted_at);


--
-- TOC entry 3572 (class 1259 OID 29058)
-- Name: idx_sales_deleted_at; Type: INDEX; Schema: public; Owner: fmaulana
--

CREATE INDEX idx_sales_deleted_at ON public.sales USING btree (deleted_at);


--
-- TOC entry 3579 (class 1259 OID 29088)
-- Name: idx_suppliers_deleted_at; Type: INDEX; Schema: public; Owner: fmaulana
--

CREATE INDEX idx_suppliers_deleted_at ON public.suppliers USING btree (deleted_at);


--
-- TOC entry 3560 (class 1259 OID 29003)
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: fmaulana
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- TOC entry 3599 (class 2606 OID 29141)
-- Name: activity_logs fk_activity_logs_user; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.activity_logs
    ADD CONSTRAINT fk_activity_logs_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3597 (class 2606 OID 29122)
-- Name: purchase_order_items fk_purchase_order_items_product; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_order_items
    ADD CONSTRAINT fk_purchase_order_items_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- TOC entry 3598 (class 2606 OID 29127)
-- Name: purchase_order_items fk_purchase_orders_items; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_order_items
    ADD CONSTRAINT fk_purchase_orders_items FOREIGN KEY (purchase_order_id) REFERENCES public.purchase_orders(id);


--
-- TOC entry 3596 (class 2606 OID 29106)
-- Name: purchase_orders fk_purchase_orders_user; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.purchase_orders
    ADD CONSTRAINT fk_purchase_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3595 (class 2606 OID 29073)
-- Name: sale_items fk_sale_items_product; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT fk_sale_items_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- TOC entry 3594 (class 2606 OID 29068)
-- Name: sale_items fk_sales_items; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT fk_sales_items FOREIGN KEY (sale_id) REFERENCES public.sales(id);


--
-- TOC entry 3593 (class 2606 OID 29053)
-- Name: sales fk_sales_user; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT fk_sales_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3591 (class 2606 OID 29029)
-- Name: stock_movements fk_stock_movements_product; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.stock_movements
    ADD CONSTRAINT fk_stock_movements_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- TOC entry 3592 (class 2606 OID 29034)
-- Name: stock_movements fk_stock_movements_user; Type: FK CONSTRAINT; Schema: public; Owner: fmaulana
--

ALTER TABLE ONLY public.stock_movements
    ADD CONSTRAINT fk_stock_movements_user FOREIGN KEY (user_id) REFERENCES public.users(id);


-- Completed on 2025-07-13 19:24:09 WIB

--
-- PostgreSQL database dump complete
--

