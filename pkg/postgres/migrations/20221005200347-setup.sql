-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE notification_status as enum ('delivered', 'failed');

CREATE TABLE notifiers (
    id uuid PRIMARY KEY not null default uuid_generate_v4(),
    endpoint varchar(256) not null,
    endpoint_method varchar(56) not null,
    client_id text not null,
    client_secret text not null,
    max_attempts smallint not null default 3,
    created_at TIMESTAMP WITH TIME ZONE not null default now()
);

CREATE TABLE notifications (
    id uuid PRIMARY KEY not null default uuid_generate_v4(),
    payload JSON not null,
    notifier_id uuid not null,
    max_attempts smallint not null,
    created_at TIMESTAMP WITH TIME ZONE not null default now(),
    CONSTRAINT fk_notifier FOREIGN KEY(notifier_id) REFERENCES notifiers(id)
);

CREATE TABLE notification_attempts (
    id uuid PRIMARY KEY not null default uuid_generate_v4(),
    state notification_status not null,
    notification_id uuid not null,
    response_status smallint not null,
    response_body smallint not null,
    created_at TIMESTAMP WITH TIME ZONE not null default now(),
    CONSTRAINT fk_notification FOREIGN KEY(notification_id) REFERENCES notifications(id)
);

-- +migrate Down
DROP TABLE notifiers;
DROP TABLE notification_attempts;
DROP TABLE notifications;
DROP TYPE notification_status;