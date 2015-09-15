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
  email      VARCHAR(320) NOT NULL,
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
  name        TEXT NOT NULL,
  description TEXT NOT NULL
);

DROP TABLE IF EXISTS message CASCADE;
CREATE TABLE IF NOT EXISTS message(
  message_id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES account(account_id),
  room_id    BIGINT NOT NULL REFERENCES room(room_id),
  body       TEXT NOT NULL
);

-- Verify Creation
\dt gogochat.*

