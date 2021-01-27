CREATE TABLE IF NOT EXISTS `tag` (
  `id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `user_id` varchar(64) NOT NULL,
  `todo_id` varchar(64) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE utf8_bin;