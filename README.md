# Clash-Style-Go-Prototype


A minimal server-authoritative backend prototype for a Clash-of-Clans-style game.

## Stack
Go 1.26
Gin HTTP framework
PostgreSQL (pgx driver)

## Features
Guest login (UUID-based provider identity)
Server-authoritative player state
Base and building system
Time-based building upgrades (server-validated)
Mock PvP battle validation

## Core Design Rules
Client never writes authoritative state directly
All game logic executed server-side
Database is the single source of truth
Upgrade timing is enforced by server time

## Authentication
No auth providers or JWT
Client sends `X-Player-ID` header after guest login
Server validates player existence on each request

## Endpoints
POST `/guest/login`
GET `/profile`
GET `/base`
POST `/building/upgrade`
POST `/pvp/submit`

## Notes
This is a prototype
No migrations or ORM included
Upgrade completion is not yet processed by a background worker (future improvement)

## Run
1. Start PostgreSQL locally
2. Set `DATABASE_URL` (optional)
3. Run:
   ```bash
   go run ./cmd/server
   ```
