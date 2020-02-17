/*
 Navicat MySQL Data Transfer

 Source Server         : webx
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : localhost:3306
 Source Schema         : blog

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 13/02/2020 17:54:45
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for blog_article
-- ----------------------------
DROP TABLE IF EXISTS `blog_article`;
CREATE TABLE `blog_article`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '标签ID',
  `title` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '文章标题',
  `desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '简述',
  `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `created_on` int(11) NULL DEFAULT NULL,
  `created_by` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '修改时间',
  `modified_by` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '修改人',
  `state` tinyint(3) UNSIGNED NULL DEFAULT 1 COMMENT '状态 0为禁用1为启用',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '文章管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for blog_auth
-- ----------------------------
DROP TABLE IF EXISTS `blog_auth`;
CREATE TABLE `blog_auth`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '账号',
  `password` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '密码',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of blog_auth
-- ----------------------------
INSERT INTO `blog_auth` VALUES (1, 'test', 'test123456');

-- ----------------------------
-- Table structure for blog_tag
-- ----------------------------
DROP TABLE IF EXISTS `blog_tag`;
CREATE TABLE `blog_tag`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '标签名称',
  `created_on` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '创建时间',
  `created_by` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '修改时间',
  `modified_by` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '修改人',
  `state` tinyint(3) UNSIGNED NULL DEFAULT 1 COMMENT '状态 0为禁用、1为启用',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '文章标签管理' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
