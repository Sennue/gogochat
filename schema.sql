-- psql -f schema.sql
-- \i schema.sql

-- Remove Role and Database
DROP OWNED BY gogochat CASCADE;
DROP DATABASE IF EXISTS gogochat;
DROP ROLE IF EXISTS gogochat;

-- Create Role and Database
CREATE ROLE gogochat PASSWORD 'gogochat' NOSUPERUSER NOCREATEDB NOCREATEROLE INHERIT LOGIN;
CREATE DATABASE gogochat OWNER gogochat;
\c gogochat
CREATE EXTENSION pgcrypto;
SET ROLE gogochat;

-- Create Schema
DROP SCHEMA IF EXISTS gogochat CASCADE;
CREATE SCHEMA IF NOT EXISTS gogochat AUTHORIZATION gogochat;
SET search_path TO gogochat,public;

-- Create Table
DROP TABLE IF EXISTS account CASCADE;
CREATE TABLE IF NOT EXISTS account(
  account_id BIGSERIAL PRIMARY KEY,
  active     BOOLEAN NOT NULL DEFAULT true,
  name       TEXT NOT NULL,
  email      VARCHAR(320) NOT NULL UNIQUE,
  password   VARCHAR(128) NOT NULL,
  salt       VARCHAR(128) NOT NULL
);

DROP TABLE IF EXISTS device CASCADE;
CREATE TABLE IF NOT EXISTS device(
  -- iOS UIDevice.currentDevice().identifierForVendor.UUIDString = 32+4
  -- Android telephonyManager.getDeviceId() = 16
  -- Need prefix, may support other platforms
  device_id  VARCHAR(128) PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES account(account_id)
);

DROP TABLE IF EXISTS room CASCADE;
CREATE TABLE IF NOT EXISTS room(
  room_id     BIGSERIAL PRIMARY KEY,
  name        TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL
);

DROP TABLE IF EXISTS message CASCADE;
CREATE TABLE IF NOT EXISTS message(
  message_id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES account(account_id),
  room_id    BIGINT NOT NULL REFERENCES room(room_id),
  body       TEXT NOT NULL
);

-- Account Creation
CREATE OR REPLACE FUNCTION add_account(device_id varchar, name varchar, email varchar, password varchar, active boolean)
	RETURNS TABLE(success BOOLEAN, account_id BIGINT) AS
$$
DECLARE
	salt VARCHAR;
BEGIN
	salt := gen_salt('bf', 8);
	password := crypt(password, salt);
	INSERT INTO account (name, email, password, salt) VALUES (name, email, password, salt);
	account_id := currval('account_account_id_seq');
	INSERT INTO device (device_id, account_id) VALUES (device_id, account_id);
	success := TRUE;
	RETURN NEXT;
EXCEPTION
	WHEN unique_violation THEN
		account_id := 0;
		success := FALSE;
		RETURN NEXT;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION add_room(name varchar, description varchar)
	RETURNS TABLE(success BOOLEAN, account_id BIGINT) AS
$$
BEGIN
	INSERT INTO room (name, description) VALUES (name, description);
	account_id := currval('room_room_id_seq');
	success := TRUE;
	RETURN NEXT;
EXCEPTION
	WHEN unique_violation THEN
		account_id := 0;
		success := FALSE;
		RETURN NEXT;
END;
$$ LANGUAGE 'plpgsql';

-- Verify Creation
\dt gogochat.*

