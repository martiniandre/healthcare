CREATE TABLE notifications (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type           VARCHAR(50) NOT NULL,
    priority       VARCHAR(20) NOT NULL DEFAULT 'medium',
    title          VARCHAR(255) NOT NULL,
    body           TEXT NOT NULL DEFAULT '',
    actor_id       UUID,
    resource_type  VARCHAR(50) NOT NULL DEFAULT '',
    resource_id    VARCHAR(255) NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notification_recipients (
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL,
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    read_at         TIMESTAMPTZ,
    PRIMARY KEY (notification_id, user_id)
);

CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notification_recipients_user_id ON notification_recipients(user_id);
CREATE INDEX idx_notification_recipients_unread ON notification_recipients(user_id, is_read) WHERE is_read = FALSE;
