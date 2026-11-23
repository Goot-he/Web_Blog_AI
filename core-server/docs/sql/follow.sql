CREATE TABLE IF NOT EXISTS user_follow
(
    id           BIGINT UNSIGNED              NOT NULL AUTO_INCREMENT,
    created_at   DATETIME(3)                  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3)                  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    follower_id  BIGINT UNSIGNED              NOT NULL,
    following_id BIGINT UNSIGNED              NOT NULL,
    status       ENUM ('following','nothing') NOT NULL DEFAULT 'nothing',
    version      BIGINT                       NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uk_follow_pair (follower_id, following_id),
    KEY idx_following_status (following_id, status),
    KEY idx_follower_status (follower_id, status)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
