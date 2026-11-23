-- ============================================================
-- web-blog-AI 数据库结构
--
-- 由 core-server 的 ent schema 生成（entgo.io/ent v0.14.6）
-- 生成工具: core-server/cmd/ddlgen（不连接业务库，只读取 MySQL 版本）
--
-- 后处理说明:
--   1. 关闭了外键生成：article.category_id 用 0 表示“未分类”，
--      开启外键会导致 0 值插入失败。仓库层由应用保证一致性。
--   2. ent 默认把 time.Time 映射为 MySQL timestamp，这里统一改为
--      datetime(3)：避免 MySQL 5.7 的隐式默认值行为，也避免 timestamp
--      只支持 1970-2038 的问题；user.birthday 改为 date。
--   3. 字符集统一为 utf8mb4 / utf8mb4_unicode_ci（与原手工 SQL 一致，
--      保证用户名/邮箱等比较是大小写不敏感的）。
--   4. ai_chat_session / ai_chat_message 属于 gateway 的 GORM 模型，
--      不在 core-server 的 ent schema 中，这里按原 ai_chat.sql 补充。
--
-- 如果导入到已存在的数据库，可跳过下面的 CREATE DATABASE / USE。
-- ============================================================

CREATE DATABASE IF NOT EXISTS `web_blog_ai`
    DEFAULT CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE `web_blog_ai`;

SET NAMES utf8mb4;

-- ------------------------------------------------------------
-- 用户表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NOT NULL,
    `updated_at` datetime(3) NOT NULL,
    `deleted_at` datetime(3) NULL,
    `name` varchar(50) NOT NULL,
    `password` varchar(20) NOT NULL,
    `phone` varchar(20) NOT NULL DEFAULT '',
    `avatar` varchar(255) NOT NULL DEFAULT '',
    `email` varchar(50) NOT NULL DEFAULT '',
    `role` enum('admin','user','unknown') NOT NULL DEFAULT 'user',
    `sex` enum('male','female','secret','unknown') NOT NULL DEFAULT 'unknown',
    `birthday` date NULL,
    `signature` varchar(255) NOT NULL DEFAULT '',
    `status` enum('approved','blocked','unknown') NOT NULL DEFAULT 'approved',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_name` (`name`),
    INDEX `idx_email` (`email`),
    INDEX `idx_phone` (`phone`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ------------------------------------------------------------
-- 用户统计表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_stat` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NOT NULL,
    `updated_at` datetime(3) NOT NULL,
    `deleted_at` datetime(3) NULL,
    `article_count` bigint unsigned NOT NULL DEFAULT 0,
    `view_count` bigint unsigned NOT NULL DEFAULT 0,
    `followers_count` bigint unsigned NOT NULL DEFAULT 0,
    `following_count` bigint unsigned NOT NULL DEFAULT 0,
    `like_count` bigint unsigned NOT NULL DEFAULT 0,
    `receive_like_count` bigint unsigned NOT NULL DEFAULT 0,
    `favor_count` bigint unsigned NOT NULL DEFAULT 0,
    `receive_favor_count` bigint unsigned NOT NULL DEFAULT 0,
    `receive_comment_count` bigint unsigned NOT NULL DEFAULT 0,
    `user_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_user_stat_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户统计表';

ALTER TABLE `user_stat`
    ADD CONSTRAINT `fk_user_stat_user_id`
        FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

-- ------------------------------------------------------------
-- 分类表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `category` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NOT NULL,
    `updated_at` datetime(3) NOT NULL,
    `parent_id` bigint unsigned NOT NULL DEFAULT 0,
    `name` varchar(64) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_category_parent_id_name` (`parent_id`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

CREATE index idx_category_parent_id on category (parent_id);


-- ------------------------------------------------------------
-- 文章表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `article` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NOT NULL,
    `updated_at` datetime(3) NOT NULL,
    `deleted_at` datetime(3) NULL,
    `published_at` datetime(3) NULL,
    `title` varchar(200) NOT NULL,
    `summary` varchar(500) NOT NULL DEFAULT '',
    `content` longtext NOT NULL,
    `cover_image` varchar(500) NOT NULL DEFAULT '',
    `is_top` bool NOT NULL DEFAULT false,
    `is_published` bool NOT NULL DEFAULT true,
    `view_count` bigint NOT NULL DEFAULT 0,
    `like_count` bigint NOT NULL DEFAULT 0,
    `favor_count` bigint NOT NULL DEFAULT 0,
    `comment_count` bigint NOT NULL DEFAULT 0,
    `category_id` bigint unsigned NULL DEFAULT 0,
    `author_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_category_id` (`category_id`),
    INDEX `idx_user_id` (`author_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='博客文章表';

-- ------------------------------------------------------------
-- 评论表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `comment` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NOT NULL,
    `deleted_at` datetime(3) NULL,
    `root_id` bigint unsigned NOT NULL DEFAULT 0,
    `reply_to_id` bigint unsigned NOT NULL DEFAULT 0,
    `content` longtext NOT NULL,
    `like_count` int unsigned NOT NULL DEFAULT 0,
    `child_count` int unsigned NOT NULL DEFAULT 0,
    `article_id` bigint unsigned NOT NULL,
    `user_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_article_id` (`article_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_root_id` (`root_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- ------------------------------------------------------------
-- 收藏以及收藏夹表
-- ------------------------------------------------------------
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

-- ------------------------------------------------------------
-- 点赞交互表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `interaction_like` (
    `id` BIGINT UNSIGNED auto_increment NOT NULL COMMENT '主键ID',
    `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
    `object_type` VARCHAR(32) NOT NULL COMMENT '对象类型（如：post, comment, video等）',
    `object_id` VARCHAR(64) NOT NULL COMMENT '对象ID',
    `status` VARCHAR(32) NOT NULL DEFAULT 'thumb_up' COMMENT '状态：thumb_up-已点赞, nothing-已取消',
    `version` BIGINT NOT NULL DEFAULT 0 COMMENT '版本号（乐观锁）',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_object` (`user_id`, `object_type`, `object_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞交互表';
-- ------------------------------------------------------------
-- 关注关系表
-- ------------------------------------------------------------
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

-- ------------------------------------------------------------
-- AI 助手会话表（gateway GORM 模型，非 ent 管理）
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ai_chat_session` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NULL,
    `updated_at` datetime(3) NULL,
    `user_id` varchar(64) NOT NULL,
    `session_id` varchar(64) NOT NULL,
    `title` varchar(256) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_ai_chat_session_id` (`session_id`),
    INDEX `idx_ai_chat_session_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI助手会话';

-- ------------------------------------------------------------
-- AI 助手对话消息表（gateway GORM 模型，非 ent 管理）
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ai_chat_message` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NULL,
    `session_id` varchar(64) NOT NULL,
    `user_id` varchar(64) NOT NULL,
    `role` varchar(16) NOT NULL,
    `content` text NOT NULL,
    `ai_model` varchar(64) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    INDEX `idx_ai_chat_message_session_created` (`session_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI助手对话消息';
