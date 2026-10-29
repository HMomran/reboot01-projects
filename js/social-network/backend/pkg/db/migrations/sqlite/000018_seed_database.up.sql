-- Seed database with sample data for demo purposes

-- Insert sample users (passwords are bcrypt hash of "password123")
INSERT OR IGNORE INTO users (id, email, password, first_name, last_name, date_of_birth, nickname, about_me, profile_type)
VALUES
  ('seed-user-001', 'alice@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Alice', 'Johnson', '1995-03-15', 'alice_j', 'Hey there! I love coding and coffee.', 'public'),
  ('seed-user-002', 'bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Bob', 'Smith', '1992-07-22', 'bob_smith', 'Adventure seeker and tech enthusiast.', 'public'),
  ('seed-user-003', 'carol@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Carol', 'Williams', '1998-11-08', 'carol_w', 'Artist and dreamer.', 'private'),
  ('seed-user-004', 'dave@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Dave', 'Brown', '1990-01-30', 'dave_b', 'Music lover and gamer.', 'public');

-- Follow relationships (Alice <-> Bob mutual, Alice -> Carol accepted, Dave -> Alice)
INSERT OR IGNORE INTO followers (id, follower_id, following_id, status)
VALUES
  ('seed-follow-001', 'seed-user-001', 'seed-user-002', 'accepted'),
  ('seed-follow-002', 'seed-user-002', 'seed-user-001', 'accepted'),
  ('seed-follow-003', 'seed-user-001', 'seed-user-003', 'accepted'),
  ('seed-follow-004', 'seed-user-004', 'seed-user-001', 'accepted');

-- Sample posts
INSERT OR IGNORE INTO posts (id, user_id, content, privacy)
VALUES
  ('seed-post-001', 'seed-user-001', 'Hello everyone! Excited to be here on the social network. 🎉', 'public'),
  ('seed-post-002', 'seed-user-002', 'Just finished a great hike in the mountains. The view was incredible!', 'public'),
  ('seed-post-003', 'seed-user-001', 'Working on a new project. Stay tuned for updates! 🚀', 'followers'),
  ('seed-post-004', 'seed-user-003', 'New painting finished today. Art is my escape. 🎨', 'public'),
  ('seed-post-005', 'seed-user-004', 'Anyone else excited for the weekend? Let me know your plans!', 'public');

-- Sample comments
INSERT OR IGNORE INTO comments (id, post_id, user_id, content)
VALUES
  ('seed-comment-001', 'seed-post-001', 'seed-user-002', 'Welcome Alice! Great to have you here! 👋'),
  ('seed-comment-002', 'seed-post-002', 'seed-user-001', 'That sounds amazing, Bob! Where did you go?'),
  ('seed-comment-003', 'seed-post-004', 'seed-user-001', 'Beautiful work, Carol! You are so talented.');

-- Sample group
INSERT OR IGNORE INTO groups (id, creator_id, title, description)
VALUES
  ('seed-group-001', 'seed-user-001', 'Tech Enthusiasts', 'A group for people who love technology, coding, and innovation.');

-- Group members (Alice is creator, Bob is a member)
INSERT OR IGNORE INTO group_members (id, group_id, user_id, status)
VALUES
  ('seed-gm-001', 'seed-group-001', 'seed-user-001', 'accepted'),
  ('seed-gm-002', 'seed-group-001', 'seed-user-002', 'accepted');

-- Sample group post
INSERT OR IGNORE INTO group_posts (id, group_id, user_id, content)
VALUES
  ('seed-gpost-001', 'seed-group-001', 'seed-user-001', 'Welcome to the Tech Enthusiasts group! Feel free to share anything tech-related here.');
