INSERT INTO users (user_id, name) VALUES
  ('haruka_01', 'Haruka'), ('naoto_02', 'Naoto'), ('riko_03', 'Riko'), ('takeshi_04', 'Takeshi'), ('suzuka_05', 'Suzuka'), ('minoru_06', 'Minoru'), ('yuka_07', 'Yuka'), ('Taro_08', 'Taro'), ('akari_09', 'Akari'), ('mai_10', 'Mai');

INSERT INTO friend_link (user1_id, user2_id) VALUES
  ('haruka_01', 'naoto_02'), ('haruka_01', 'riko_03'), ('haruka_01', 'akari_09'), ('takeshi_04', 'haruka_01'), ('taro_08', 'akari_09');

INSERT INTO block_list (user1_id, user2_id) VALUES
  ('yuka_07', 'haruka_01'), ('riko_03', 'takeshi_04');

INSERT INTO friend_requests (user1_id, user2_id, status) VALUES
  ('suzuka_05', 'haruka_01', 'pending');
