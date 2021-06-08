# acs-builder

A Traveller5 starship component builder.

# Theory

Relatively small, stateless endpoints.
Incrementally developed.

# API

## GET /sensors

Returns a list of all sensors in the database.

## POST /sensors

Add a sensor to the database.

## POST /sensors/{type}

Build a sensor from existing mounts, ranges, and types.

Request Body:

- mount (optional): enum (T1 | T2 | T3 | T4 | B1 | B2)
- range (optional): enum (Vl | D | Vd | Or | Fo | G | BR | FR | SR | AR | LR | DS)

## GET /weapons

Returns a list of all weapons in the database.

## POST /weapons

Add a weapon to the database.

## POST /weapons/{type}

Build a weapon from existing mounts, ranges, and types.  Uses the same request body as the sensor.

## GET /mounts

Returns a list of all mounts in the database.

## POST /mounts

Add a mount to the database.

## GET /ranges

Returns a list of all ranges in the database.

## POST /ranges

Add a range to the database.

# To Do

- Swagger doc.
- Decouple methods.
- Fill out maps.
- Externalize maps.
- Weapons.
- Defenses.
- Hull.
- Drives.
- Etc.

