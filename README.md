# acs-builder

A Traveller5 starship component builder.

# Theory

Relatively small, stateless endpoints.
Incrementally developed.

# API

## GET /mounts - Returns a list of all mounts in the database.
## GET /ranges - Returns a list of all ranges in the database.
## GET /sensors - Returns a list of all sensor types in the database.
## GET /weapons - Returns a list of all weapon types in the database.
## GET /drives - Returns a list of all drive types in the database.
## GET /hulls - Returns a list of all hull configurations in the database.

## POST /mounts - Add a mount to the database.
## POST /ranges - Add a range to the database.
## POST /sensors - Add a sensor type to the database.
## POST /weapons - Add a weapon type to the database.

## POST /sensors/{type}

Build a sensor from existing mounts, ranges, and types.

Request Body:

- mount (optional): enum (T1 | T2 | T3 | T4 | B1 | B2)
- range (optional): enum (Vl | D | Vd | Or | Fo | G | BR | FR | SR | AR | LR | DS)

## POST /weapons/{type}
 
Build a weapon from existing mounts, ranges, and types.  Uses the same request body as the sensor.

## POST /drives/{type}

Build a drive from existing types.  

Request Body:

- rating: numeric drive rating (e.g. 2 for Jump-2 or Maneuver-2)
- targetHullTons: hull tons (e.g. 200 for a 200t hull; 10000 for a 10 kiloton hull)

## POST /hulls/{config}

Build a hull with the given configuration.

Request Body:

- tons: volume of hull, in tons
- tl: TL of the hull


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

