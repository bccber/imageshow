SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `t_images`
-- ----------------------------
DROP TABLE IF EXISTS `t_images`;
CREATE TABLE `t_images` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `md5` varchar(32) NOT NULL,
  `title` varchar(256) NOT NULL,
  `url` varchar(256) NOT NULL,
  `state` bit(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `md5` (`md5`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_images
-- ----------------------------
