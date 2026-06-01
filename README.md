# System Design Compendium

## Overview

This repo is a Markdown-first learning compendium for system design concepts. Each topic is designed to be quick to read, easy for LLMs to navigate, and practical enough to turn into small implementation exercises.

The goal is to pair concise concept notes and curated resources with focused exercises and reference solutions.

## Topics

- [Networking](networking/README.md): protocols, service communication patterns, and small API/RPC exercises.
- [Databases](databases/README.md): Postgres fundamentals, indexing, locking, partitioning, Cassandra modeling, and Elasticsearch search.
- [Caching](caching/README.md): cache-aside, TTLs, counters, rate limiting, and ephemeral state.

## Repository Pattern

Each topic follows this shape:

```text
<topic>/
├── README.md
├── resources.md
└── exercises/
    └── <exercise>/
        ├── README.md
        ├── IMPLEMENTATION_GUIDE.md
        └── solution/
```

Not every exercise has a solution yet. README files define the exercise and learning goals; implementation guides provide a step-by-step path; solution folders hold reference implementations when available.

## How To Use

1. Start with a topic README.
2. Read the exercise README for the problem statement.
3. Try the implementation yourself.
4. Use the implementation guide if you get stuck.
5. Compare with the reference solution when one exists.
