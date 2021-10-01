# Alfred

General purpose modular Slack bot for Workday Orbit 2021 Hackathon.

## Services

Services facilitate commands through HTTP communication. When a service is registered through the bot, it has a trigger command
that can be used to interact with the service.

Services can be registered using the `register` command.
```
register [service_name] [trigger] [service_ip] [service_port]
```

Services can be deregistered using the `deregister` command.
```
deregister [trigger]
```

Commands currently usable can be access through the `commands` command.

### Current Services

|          Service        |                          URL                          |                                      Description                                         |
| ----------------------- | ----------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| Resource Service        | https://github.com/JasonVanRaamsdonk/example_rest_api | A service that gives easy access to handy internal Workday resources                     |
| Covid Return to Office  | https://github.com/mcDevittMaya5/hackathon            | A service that provides utilities for returning to the office such as temperature checks |
| Health Reminder Service | https://github.com/GohEeEn/health-reminder-service    | A service that reminds users to practice healthy activities while working from home      |