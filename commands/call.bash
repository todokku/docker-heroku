#!/usr/bin/env bash

curl http://localhost:8080/slack -X POST \
-d @- <<'JSON'
{
    "token": "token.value",
    "team_id": "team_id.value",
    "api_app_id": "api_app_id.value",
    "event": {
        "client_msg_id": "client_msg_id.value",
        "type": "message",
        "text": "test",
        "user": "user.value",
        "ts": "1583248428.001900",
        "team": "team.value",
        "blocks": [
            {
                "type": "rich_text",
                "block_id": "block_id.value",
                "elements": [
                    {
                        "type": "rich_text_section",
                        "elements": [
                            {
                                "type": "text",
                                "text": "test"
                            }
                        ]
                    }
                ]
            }
        ],
        "channel": "channel.value",
        "event_ts": "1583248428.001900",
        "channel_type": "channel"
    },
    "type": "event_callback",
    "event_id": "event_id.value",
    "event_time": 1583248428,
    "authed_users": [
        "authed_users.value1",
        "authed_users.value2"
    ]
}
JSON
