-- =====================================================================
-- pitz_backend - PostgreSQL schema
-- =====================================================================
-- Equivalente al esquema MySQL original.
-- Tipos de columna basados en los modelos Go (models/*.go) y las queries
-- SQL usadas en routes/*.go.
-- =====================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================================
-- USERS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "users" (
    "user_uid"         VARCHAR(36) PRIMARY KEY,
    "username"         VARCHAR(100) NOT NULL UNIQUE,
    "email"            VARCHAR(255) NOT NULL UNIQUE,
    "hashed_password"  VARCHAR(255) NOT NULL,
    "first_name"       VARCHAR(100) NOT NULL,
    "last_name"        VARCHAR(100) NOT NULL,
    "created_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag"      SMALLINT NOT NULL DEFAULT 0
);

-- =====================================================================
-- PLAYERS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "players" (
    "player_uid"       VARCHAR(36) PRIMARY KEY,
    "first_name"       VARCHAR(100) NOT NULL,
    "last_name"        VARCHAR(100) NOT NULL,
    "email"            VARCHAR(255),
    "status"           SMALLINT NOT NULL DEFAULT 1,
    "positions"        JSONB,
    "phone_number"     VARCHAR(30),
    "emergency_phone"  VARCHAR(30),
    "birth_dt"         TIMESTAMP,
    "blood_type"       VARCHAR(10),
    "comments"         TEXT,
    "credential"       VARCHAR(100),
    "address"          TEXT,
    "afiliation"       VARCHAR(100),
    "sex"              VARCHAR(20),
    "curp"             VARCHAR(50),
    "enfermedad"       TEXT,
    "insurance"        BOOLEAN NOT NULL DEFAULT FALSE,
    "insurance_name"   VARCHAR(100),
    "created_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag"      SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS "idx_players_first_name" ON "players" ("first_name");

-- =====================================================================
-- EVENT TYPES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "event_types" (
    "event_type_uid"   VARCHAR(36) PRIMARY KEY,
    "type_name"        VARCHAR(100) NOT NULL
);

-- =====================================================================
-- EVENT STATES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "events_state" (
    "event_state_uid"  VARCHAR(36) PRIMARY KEY,
    "event_state"      VARCHAR(100) NOT NULL
);

-- =====================================================================
-- EVENTS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "events" (
    "event_uid"        VARCHAR(36) PRIMARY KEY,
    "event_type"       VARCHAR(36) NOT NULL,
    "event_state_uid"  VARCHAR(36) NOT NULL,
    "event_name"       VARCHAR(255) NOT NULL,
    "date"             TIMESTAMP NOT NULL,
    "address"          TEXT,
    "coordinates"      VARCHAR(255),
    "created_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm"  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag"      SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_events_event_type"
        FOREIGN KEY ("event_type") REFERENCES "event_types" ("event_type_uid"),
    CONSTRAINT "fk_events_event_state"
        FOREIGN KEY ("event_state_uid") REFERENCES "events_state" ("event_state_uid")
);

-- =====================================================================
-- ASISTANCE TYPES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "asistance_types" (
    "asistance_type_uid" VARCHAR(36) PRIMARY KEY,
    "name"               VARCHAR(100) NOT NULL
);

-- =====================================================================
-- ASISTANCE
-- =====================================================================
CREATE TABLE IF NOT EXISTS "asistance" (
    "player_uid"          VARCHAR(36) NOT NULL,
    "event_uid"           VARCHAR(36) NOT NULL,
    "asistance_type_uid"  VARCHAR(36) NOT NULL,
    "created_at_dttm"     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("player_uid", "event_uid"),
    CONSTRAINT "fk_asistance_player"
        FOREIGN KEY ("player_uid") REFERENCES "players" ("player_uid"),
    CONSTRAINT "fk_asistance_event"
        FOREIGN KEY ("event_uid") REFERENCES "events" ("event_uid"),
    CONSTRAINT "fk_asistance_type"
        FOREIGN KEY ("asistance_type_uid") REFERENCES "asistance_types" ("asistance_type_uid")
);

-- =====================================================================
-- PAYMENT TYPES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "payment_type" (
    "payment_type_uid" VARCHAR(36) PRIMARY KEY,
    "payment_name"     VARCHAR(100) NOT NULL
);

-- =====================================================================
-- PAYMENTS
-- =====================================================================
CREATE TABLE IF NOT EXISTS "payments" (
    "payment_uid"         VARCHAR(36) PRIMARY KEY,
    "player_uid"          VARCHAR(36) NOT NULL,
    "payment_reference"   VARCHAR(255),
    "amount"              NUMERIC(12, 2) NOT NULL,
    "comment"             TEXT,
    "date"                TIMESTAMP NOT NULL,
    "payment_type_uid"    VARCHAR(36) NOT NULL,
    "registered_by_uid"   VARCHAR(36) NOT NULL,
    "created_at_dttm"     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm"     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag"         SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_payments_player"
        FOREIGN KEY ("player_uid") REFERENCES "players" ("player_uid"),
    CONSTRAINT "fk_payments_payment_type"
        FOREIGN KEY ("payment_type_uid") REFERENCES "payment_type" ("payment_type_uid"),
    CONSTRAINT "fk_payments_registered_by"
        FOREIGN KEY ("registered_by_uid") REFERENCES "users" ("user_uid")
);

CREATE INDEX IF NOT EXISTS "idx_payments_date" ON "payments" ("date");
CREATE INDEX IF NOT EXISTS "idx_payments_delete_flag" ON "payments" ("delete_flag");

-- =====================================================================
-- EXPENSES
-- =====================================================================
CREATE TABLE IF NOT EXISTS "expenses" (
    "expense_uid"        VARCHAR(36) PRIMARY KEY,
    "assigned_uid"       VARCHAR(36) NOT NULL,
    "registered_by_uid"  VARCHAR(36) NOT NULL,
    "reason"             TEXT NOT NULL,
    "amount"             NUMERIC(12, 2) NOT NULL,
    "date"               TIMESTAMP NOT NULL,
    "created_at_dttm"    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at_dttm"    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delete_flag"        SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_expenses_assigned"
        FOREIGN KEY ("assigned_uid") REFERENCES "users" ("user_uid"),
    CONSTRAINT "fk_expenses_registered_by"
        FOREIGN KEY ("registered_by_uid") REFERENCES "users" ("user_uid")
);

CREATE INDEX IF NOT EXISTS "idx_expenses_date" ON "expenses" ("date");
CREATE INDEX IF NOT EXISTS "idx_expenses_delete_flag" ON "expenses" ("delete_flag");

-- =====================================================================
-- Usuario admin inicial (cambiar password obligatorio)
-- hash bcrypt de "admin123" - REEMPLAZAR antes de producción
-- =====================================================================
-- INSERT INTO "users" ("user_uid", "username", "email", "hashed_password",
--                     "first_name", "last_name")
-- VALUES (
--     uuid_generate_v4()::text,
--     'admin',
--     'admin@pitz.local',
--     '$2a$10$REEMPLAZAR_CON_HASH_REAL',
--     'Admin',
--     'PITZ'
-- );

-- =====================================================================
-- Catálogos iniciales sugeridos
-- =====================================================================
-- INSERT INTO "event_types" ("event_type_uid", "type_name") VALUES
--     (uuid_generate_v4()::text, 'Entrenamiento'),
--     (uuid_generate_v4()::text, 'Partido'),
--     (uuid_generate_v4()::text, 'Reunión');
--
-- INSERT INTO "events_state" ("event_state_uid", "event_state") VALUES
--     (uuid_generate_v4()::text, 'Pendiente'),
--     (uuid_generate_v4()::text, 'Confirmado'),
--     (uuid_generate_v4()::text, 'Cancelado'),
--     (uuid_generate_v4()::text, 'Finalizado');
--
-- INSERT INTO "asistance_types" ("asistance_type_uid", "name") VALUES
--     (uuid_generate_v4()::text, 'Asistió'),
--     (uuid_generate_v4()::text, 'Faltó'),
--     (uuid_generate_v4()::text, 'Justificado');
--
-- INSERT INTO "payment_type" ("payment_type_uid", "payment_name") VALUES
--     (uuid_generate_v4()::text, 'Mensualidad'),
--     (uuid_generate_v4()::text, 'Inscripción'),
--     (uuid_generate_v4()::text, 'Multa');
