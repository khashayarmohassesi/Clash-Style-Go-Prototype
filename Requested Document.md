# Player Authentication 
For player authentication given that these games usually have a guest player feature, each client can start playing the game without an account that is linked to an email/other socials and they can progress until the game forces it (by locking some features) or leave it open till they choose to connect their account. (to continue their progress on another device).
The authentication could be Oauth/JWT tokens and refresh tokens based on what stage the player is playing with.

# Player Profile
Keeping general data, progression data, cosmetic data, last login, and clan_id.

# Base Save/Load
For start could be as simple as a JSON file passed from server to the client, validating the ownership of the data, and caching with Redis, for saving/persistence a database would be needed (NoSQL or SQL)



# Anti-cheating
Anti-cheating comes down to not trusting the client on the packets/requests that they will send to the server and validating everything, do they have enough resources? is the action they're trying to perform available at their tier? can they employ the forces they're trying to send? are those forces in the correct tier?... 
And all of this largely gets solved by server authoritative architectures, almost all of the actions are validated on the server and only then actions are performed on the server and the client just receives the updates.
Moreover, each game can and will result in different kinds of detecting cheating, the most common one in MMO-style games are how many requests each client can send towards the server given a real user is on the other end and not a cheater.

# Resource Generation Validation 

Basically all the validations and anti-cheating is the result of not trusting the client (clients can send faulty packets) and resource generation is no different, all the resource generations of each player with every bonus that they might have will be validated upon entry of the player/enough time has passed

# Upgrade Validation

Upgrades will be stored onto a database with the start date of the request of upgrade (not trusting the clock of the device from the client) and validated (the client has enough resources, can actually start the upgrade (not upgrading to 2 tiers above the current level)), and with that data, the date of the finished date of upgrade is also captured, so there's no need trusting the client to say when upgrade is done, the server will validate the request, the upgrade timer, and when it's done sends the data back as upgrade is done and the client will receive the update.

# PvP attack result validation

PvP attacks largely depend on the game itself, there are usually 2 categories.
#### Limited player participation in the attack
The players don't really participate in the attack itself, their forces and resources will determine the result of the attack and the logic is completely run on the server.

#### Player driven attacks
The players are the driving force of the result of the attack, the attacks are still largely tied to the resources and forces available to each player, but the attacking player has a lot of choices to make during the attack and skill will make a difference in the result.

Limited player participation in the attacks are already cheating-proof, because almost all of the logic of the attack is running on the server and the only validation part must be done on the forces/resources each player has, (how many units/which tier they are/how much resources are available to plunder/are there shields/bonuses active during the attack/defense) and the only part remaining is just the result to be dispatched. Player driven attacks on the other hand require all of the validation plus other techniques.
Depending on the style of the attack/game, one to multiple techniques are required. most importantly is client action validation and all actions are passing through the server, this can be taken one step further as most games implement replays, and the server can just use that base as a validation technique, player actions are recorded with their timestamps and sent to the server, the validation happens as the replay is re-simulated on the server, and kept for the observers as well.

# Matchmaking
For matchmaking, each player will have MMR rating associated with them, which can be derived from their tier/upgrades/how many successful attacks they had and the server starts searching for players in the same range of MMR as the player requesting for the match(attack). starting range could be ±100 MMR and each ~5 minutes the range would expand for striking a good balance between finding players in equal footing and finding players when there are no matches are being found.

# Alliance System
The server provides different requests from the clients for starting an alliance, adding members, deleting members, join request, chat and war/event participation.

# Leaderboards
Most common implementations are Redis sorted sets, could be based on MMR/global trophies/clans rankings/seasonal wins with periodic syncs to the database.

# Live events/Remote config

Usually is a combination of multiple features coming together, event toggling remote configurations, with these remote configurations usually comes with A/B testing, seasonal boosts could be implemented with live events as well.

# Rest API vs Websockets

Most of the communication between the client and server will be through the REST API, for login, clan API, upgrades, etc...
The less important communications will be handled with websockets, stuff like clan chat, notifications...
# Recommended Tech stack

For backend I'd recommend Golang as the primary backend language, as go offers very good concurrency and has strong tooling for backend implementation, for Rest API, Gin is my pick, for database any SQL database would suffice (or even nosql) my pick would be PostgreSQL because of good JSON support (JSON support for base layout), Caching any type of Redis will be used, and the Unity Client which would be a thin client communicating with the server, with Rest APIs as the main driving force of the majority of the packets(Login, upgrades, economy, clan APIs, battle submission), and Websockets for less security intensive communication.(clan chat, notifications) 

# What can be delivered in 30 days
This is a very broad question and I can't really give an exact roadmap, since it would require me to know much more detail about the project, the scale and scope and design on the game, and the scope of my responsibilities, with a lot of assumptions (which could be false) in the first 30 days the core gameplay and validation of the major actions would be done and the project ready to move forward to other phases.

# Architecture Scaling & Risks
I'll go over this over 2-3 phases : 
First phase: 0-50K DAUs 
Client will always be Unity, going through the API Backend and the backend will be Modular Monolith.
Second Phase: 50K to 500K DAUs:
Split out high load systems, these could be: Battle Validation, Matchmaking, Leaderboards...
Third Phase: 500K DAUs to Million users:
More distribution by sharding, regions, player partitioning... 
This phase is highly speculative and is dependent on the project itself, but the principles remain the same

For the risks: 
The most important risk would be Economy cheating:
Fake resource claims, false loot, fake upgrades, clock hacking, botting

majority of the risks can be mitigated with server authoritative logic, for botting, each game requires its own consideration for anti-botting measures

Battle result cheating: 
Players can fake many things, server authoritative logic can mitigate alot of this, on top of it, replay validation/ deterministic re-simulation will mitigate other risks as well

## players
id PK  
auth_provider  
provider_user_id  
name  
xp  
trophies  
level  
clan_id  
created_at  
last_login
## bases
id PK  
player_id FK  
layout_version  
shield_until  
last_resource_calc
## buildings
(note: buildings and layouts and base could be changed depending on the complexity of the game and use a completely different version)
id PK  
base_id FK  
type  
level  
x  
y  
hp  
state
## resources
player_id FK  
gold  
elixir  
gems  
updated_at
## upgrades
id PK  
building_id FK  
from_level  
to_level  
start_at  
finish_at  
status
## battles
id PK  
attacker_id  
defender_id  
replay_url  
stars  
destruction  
loot_gold  
loot_elixir  
validated  
created_at


## clans
id PK  
name  
owner_id  
description  
trophies


## clan_members

clan_id  
player_id  
role  
joined_at
## leaderboards

id  
scope  
entity_id  
score  
season  
rank

## live events
id  
name  
config_json  
start_at  
end_at  
enabled