SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `t_comments`
-- ----------------------------
DROP TABLE IF EXISTS `t_comments`;
CREATE TABLE `t_comments` (
  `id` bigint(20) NOT NULL DEFAULT '0',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `imgid` bigint(20) NOT NULL DEFAULT '0',
  `username` varchar(24) NOT NULL DEFAULT '',
  `content` varchar(1024) NOT NULL DEFAULT '',
  `created_time` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_comments
-- ----------------------------

-- ----------------------------
-- Table structure for `t_images`
-- ----------------------------
DROP TABLE IF EXISTS `t_images`;
CREATE TABLE `t_images` (
  `id` bigint(20) NOT NULL DEFAULT '0',
  `md5` varchar(32) NOT NULL DEFAULT '',
  `title` varchar(256) NOT NULL DEFAULT '',
  `url` varchar(256) NOT NULL DEFAULT '',
  `comment_count` int(11) NOT NULL DEFAULT '0',
  `like_count` int(11) NOT NULL DEFAULT '0',
  `created_time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `md5` (`md5`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_images
-- ----------------------------

-- ----------------------------
-- Table structure for `t_likes`
-- ----------------------------
DROP TABLE IF EXISTS `t_likes`;
CREATE TABLE `t_likes` (
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `imgid` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`,`imgid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_likes
-- ----------------------------

-- ----------------------------
-- Table structure for `t_users`
-- ----------------------------
DROP TABLE IF EXISTS `t_users`;
CREATE TABLE `t_users` (
  `id` bigint(20) NOT NULL DEFAULT '0',
  `name` varchar(24) NOT NULL DEFAULT '',
  `password` varchar(24) NOT NULL DEFAULT '',
  `remark` varchar(255) NOT NULL DEFAULT '',
  `created_time` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_users
-- ----------------------------
