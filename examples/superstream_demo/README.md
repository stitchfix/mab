# Contextual multi-armed bandit

With microservice reward source

## Starting the services:

`docker compose build && docker compose up`

## Services

### Bandit

The bandit service is a stateless app that uses Thompson sampling to select an arm for a contextual multi-armed bandit.
It depends on the reward service, which provides reward estimates for each arm depending on context.

The bandit service does not use the context directly, but just passes it to the reward service. The reward service is
responsible for validating the context.

For example:

`curl -XPOST localhost:1338/select_arm -d '{"unit": "visitor_id:12345", "context": {"source_id": 1}}'`

The bandit service will pass the value under the "context" key as a top-level JSON object in the request to the reward
service.

### Reward

The reward service is a stateful service that provides reward estimates given a context. In this basic example the
rewards are hard-coded, but a real reward service would be connected to a DB.

You can query the reward service directly with:

`curl -i -XPOST localhost:1337/rewards -d '{"source_id": 1}'`

The reward service returns an error if the context is invalid or there are no reward estimates for the given context,
otherwise it returns the reward estimate for each arm.