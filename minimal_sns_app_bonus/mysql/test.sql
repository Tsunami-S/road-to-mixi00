INSERT INTO users (user_id, name) VALUES
  ('id01', 'user01'), ('id02', 'user02'), ('id03', 'user03'), ('id04', 'user04'), ('id05', 'user05'), ('id06', 'user06'), ('id07', 'user07'), ('id08', 'user08'), ('id09', 'user09'), ('id10', 'user10'), ('id11', 'user11'), ('id12', 'user12');

INSERT INTO friend_link (user1_id, user2_id) VALUES
  ('id01', 'id02'), ('id01', 'id03'), ('id04', 'id01'), ('id05', 'id01'), ('id02', 'id03'),  ('id02', 'id06'), ('id02', 'id07'), ('id08', 'id02'), ('id09', 'id02');

INSERT INTO block_list (user1_id, user2_id) VALUES
  ('id01', 'id07'), ('id09', 'id01');

INSERT INTO friend_requests (user1_id, user2_id, status) VALUES
  ('id01', 'id10', 'pending'), ('id11', 'id01', 'pending');
