INSERT INTO users (user_id, name) VALUES
  (1, 'user01'), (2, 'user02'), (3, 'user03'), (4, 'user04'), (5, 'user05'), (6, 'user06'), (7, 'user07'), (8, 'user08'), (9, 'user09'), (10, 'user10'), (11, 'user11'), (12, 'user12');

INSERT INTO friend_link (user1_id, user2_id) VALUES
  (1, 2), (1, 3), (4, 1), (5, 1), (2, 3),  (2, 6), (2, 7), (8, 2), (9, 2), (2, 1), (1, 1);

INSERT INTO block_list (user1_id, user2_id) VALUES
  (1, 7), (9, 1), (1, 9), (1, 10), (1, 10), (1, 1);
