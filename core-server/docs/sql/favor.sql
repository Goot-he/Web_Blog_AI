-- Default personal favorites collection. A row is retained on cancellation so
-- version ordering can safely reject a delayed older request.
CREATE TABLE IF NOT EXISTS `interaction_favor` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `article_id` BIGINT UNSIGNED NOT NULL,
    `folder_id` BIGINT UNSIGNED NOT NULL,
    `status` VARCHAR(16) NOT NULL DEFAULT 'nothing',
    `version` BIGINT NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_favor_user_article` (`user_id`, `article_id`),
    KEY `idx_favor_user_status_version` (`user_id`, `status`, `version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章收藏关系表';

CREATE TABLE IF NOT EXISTS `favorite_folder` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(50) NOT NULL,
    `is_default` TINYINT(1) NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_favorite_folder_user_name` (`user_id`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户收藏夹';

-- Existing deployments: create a default folder for every prior favorite owner,
-- then assign old favorite rows before requiring folder_id to be non-null.
-- ALTER TABLE interaction_favor ADD COLUMN folder_id BIGINT UNSIGNED NULL AFTER article_id;
-- INSERT IGNORE INTO favorite_folder (user_id, name, is_default) SELECT DISTINCT user_id, '默认收藏夹', 1 FROM interaction_favor;
-- UPDATE interaction_favor f JOIN favorite_folder d ON d.user_id = f.user_id AND d.is_default = 1 SET f.folder_id = d.id WHERE f.folder_id IS NULL;
-- ALTER TABLE interaction_favor MODIFY folder_id BIGINT UNSIGNED NOT NULL;
