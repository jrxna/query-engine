# Hypertool Query Engine

The Hypertool Query Engine is our microservice that provides a high-performance and fault-tolerant interface for request queueing, query execution, query result caching and connection pooling. This API documentation makes use of codenames.

The current version of the query engine (codenamed `Xenon`) is a WIP. The API routes are grouped under `/v1`.

Once complete, the future versions of query engine will be focusing on migration from REST to GraphQL (codenamed `Triton`) and optimizing performance for scalability (codenamed `Photon`). 