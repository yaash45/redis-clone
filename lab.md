# 🧠 Build Your Own Redis — Project Plan

## 📌 Project Philosophy

This is not about cloning Redis feature-by-feature.

You are:

* Building a **mental model of a database system**
* Understanding **systems tradeoffs**
* Progressing from **simple → complex** deliberately


# 🏗️ High-Level Milestones

1. Core Engine (in-memory KV store)
2. Networking (TCP server)
3. Protocol (command parsing)
4. Event Loop (non-blocking I/O)
5. Persistence (RDB + AOF concepts)
6. Expiry system
7. Data structures
8. Memory management
9. Performance analysis
10. Optional advanced features


# 📦 Stage 0 — Groundwork

## Objective

Understand what you're building before writing code.

## Tasks

* Read about:

  * In-memory databases
  * Event-driven systems
  * Why Redis is (mostly) single-threaded
* Skim architecture docs for Redis

## Deliverable

Write a short design doc:

* What guarantees will your system provide?
* What will it explicitly NOT support?


# 🧱 Stage 1 — Core Key-Value Store

## Objective

Build a minimal in-memory database.

## Requirements

* Keys: strings
* Values: strings
* Operations:

  * `SET key value`
  * `GET key`
  * `DEL key`

## Constraints

* No networking
* No persistence
* Single-threaded

## Key Concepts

* Hash maps
* API design

## Deliverable

A CLI-like interface:

```
SET x 10
GET x
→ 10
```


# 🌐 Stage 2 — TCP Server

## Objective

Expose your database over a network.

## Requirements

* Open a TCP socket
* Accept multiple clients
* Read and write raw bytes

## Constraints

* Still single-threaded
* Blocking I/O allowed (for now)

## Key Concepts

* Sockets
* Connection lifecycle

## Deliverable

Interact using:

* `telnet`
* `nc`


# 🧾 Stage 3 — Protocol Design

## Objective

Define how clients communicate with your server.

## Requirements

* Parse commands:

  ```
  SET foo bar
  GET foo
  ```
* Return structured responses:

  * Success
  * Errors
  * Values

## Stretch Goal

* Implement a simplified RESP-like protocol

## Key Concepts

* Parsing
* Protocol design


# 🔁 Stage 4 — Event Loop (Critical)

## Objective

Handle multiple clients efficiently.

## Requirements

* Replace blocking I/O with:

  * `select`, `poll`, or equivalent
* Single-threaded multiplexing

## Key Concepts

* Event loops
* Non-blocking I/O

## Deliverable

Server handles multiple clients without blocking.


# 💾 Stage 5 — Persistence

## Objective

Survive restarts.


## Part A — Snapshotting (RDB-style)

### Requirements

* Periodically write full dataset to disk


## Part B — Append-Only Log (AOF-style)

### Requirements

* Log every write operation


## Key Questions

* When do you flush to disk?
* What happens on crash?

## Key Concepts

* Durability
* Tradeoffs between safety and performance


# ⏳ Stage 6 — Expiry System

## Objective

Support time-based key expiration.

## Requirements

* `SET key value EX seconds`
* Automatic key expiration

## Design Decisions

* Lazy vs active expiration
* Data structures for tracking TTL

## Key Concepts

* Time-based eviction


# 🧮 Stage 7 — Data Structures

## Objective

Expand beyond simple strings.

## Add Support For

* Lists
* Sets
* Hashes

## Example Commands

* `LPUSH`, `RPUSH`, `LPOP`
* `SADD`, `SREM`
* `HSET`, `HGET`

## Key Concepts

* Internal representations
* Tradeoffs between structures


# 🧠 Stage 8 — Memory Management

## Objective

Handle limited memory.

## Features

* Max memory limit
* Eviction policies:

  * LRU
  * (Optional) LFU

## Key Concepts

* Cache eviction strategies
* Approximation vs precision


# ⚡ Stage 9 — Performance Analysis

## Objective

Understand system bottlenecks.

## Tasks

* Measure:

  * Throughput
  * Latency
* Identify:

  * Hot paths
  * Memory overhead

## Key Concepts

* Real-world vs theoretical performance
* CPU cache effects


# 🔁 Stage 10 — Optional Advanced Features

## Replication

* Leader → follower synchronization

## Pub/Sub

* Clients subscribe to channels


# 🧪 Testing Strategy

At every stage:

## Write tests for:

* Correctness
* Edge cases

## Simulate:

* Concurrent clients
* Crashes (for persistence)


# 🧭 Suggested Timeline

| Week | Focus                   |
| ---- | ----------------------- |
| 1    | Stage 0–1               |
| 2    | Stage 2–3               |
| 3    | Stage 4                 |
| 4    | Stage 5                 |
| 5    | Stage 6–7               |
| 6+   | Optimization + advanced |


# 🧩 Development Discipline

For every stage:

1. Write a short design doc (1–2 pages)
2. List assumptions
3. Identify edge cases
4. Implement incrementally
5. Test before moving on
