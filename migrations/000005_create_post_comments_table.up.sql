CREATE TABLE `post_comments` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `post_id`     BIGINT UNSIGNED NOT NULL COMMENT '帖子业务ID(post.post_id)',
    `user_id`     BIGINT UNSIGNED NOT NULL COMMENT '评论用户ID',
    `content`     TEXT NOT NULL COMMENT '评论内容',
    `like_count`  INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数（可选，先不加也行）',
    `status`      TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1正常, 2删除',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',

    PRIMARY KEY (`id`),
    KEY `idx_post_time` (`post_id`, `create_time` DESC),   -- 按帖子 + 时间倒序查一级评论列表
    KEY `idx_user_id` (`user_id`)                          -- 用户评论列表（可选）
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci 
  COMMENT='帖子一级评论表';