-- Remove seed data
DELETE FROM group_posts WHERE id LIKE 'seed-%';
DELETE FROM group_members WHERE id LIKE 'seed-%';
DELETE FROM groups WHERE id LIKE 'seed-%';
DELETE FROM comments WHERE id LIKE 'seed-%';
DELETE FROM posts WHERE id LIKE 'seed-%';
DELETE FROM followers WHERE id LIKE 'seed-%';
DELETE FROM users WHERE id LIKE 'seed-%';
