DROP TABLE IF EXISTS block_list;
DROP TABLE IF EXISTS friend_link;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id BIGINT(20) NOT NULL AUTO_INCREMENT,
  user_id int(11) NOT NULL,
  name VARCHAR(64) NOT NULL DEFAULT '',
  PRIMARY KEY (id)
);

CREATE TABLE friend_link (
  id BIGINT(20) NOT NULL AUTO_INCREMENT,
  user1_id int(11) NOT NULL,
  user2_id int(11) NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE block_list (
  id BIGINT(20) NOT NULL AUTO_INCREMENT,
  user1_id int(11) NOT NULL,
  user2_id int(11) NOT NULL,
  PRIMARY KEY (id)
);
