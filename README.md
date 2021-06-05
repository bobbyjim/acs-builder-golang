# acs-builder

A Traveller5 starship component builder.

# Theory

Relatively small, stateless endpoints.
Incrementally developed.

# API

## POST /sensor/{type}

Create a sensor. Type is a single capital letter from: 

- C: Communicator
- H: HoloVisor
- T: Scope
- V: Visor
- W: CommPlus
- E: EMS
- G: Grav
- N: Neutrino
- R: Radar
- S: Scanner
- A: Activity

Request Body:

- mount (optional): enum (T1 | T2 | T3 | T4 | B1 | B2)
- range (optional): enum (Vl | D | Vd | Or | Fo | G | BR | FR | SR | AR | LR | DS)

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

