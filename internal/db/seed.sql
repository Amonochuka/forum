INSERT INTO users (id, email, username, password_hash, created_at) VALUES
(1, 'dev_dan@example.com', 'DevDan', 'hash_password_1', '2026-01-01 10:00:00'),
(2, 'sarah_codes@example.com', 'SarahCodes', 'hash_password_2', '2026-01-02 11:30:00'),
(3, 'tech_guru@example.com', 'TechGuru', 'hash_password_3', '2026-01-05 09:15:00'),
(4, 'alex_m@example.com', 'AlexM', 'hash_password_4', '2026-01-10 14:00:00'),
(5, 'nina_v@example.com', 'NinaV', 'hash_password_5', '2026-01-12 16:45:00'),
(6, 'rust_fan@example.com', 'RustFan', 'hash_password_6', '2026-01-15 08:20:00');

INSERT INTO categories (id, name, created_at) VALUES
(1, 'General', '2026-01-01 00:00:00'),
(2, 'Programming', '2026-01-01 00:00:00'),
(3, 'Hardware', '2026-01-01 00:00:00'),
(4, 'Showcase', '2026-01-01 00:00:00'),
(5, 'Career Advice', '2026-01-01 00:00:00');

INSERT INTO posts (id, user_id, title, content, created_at) VALUES
(1, 1, 'Welcome to the Forum!', 'This is the first post. Welcome everyone to our new community!', '2026-01-01 10:05:00'),
(2, 2, 'React vs Vue in 2026', 'What are you all using for frontend these days? I am leaning towards Vue 4.', '2026-02-10 12:00:00'),
(3, 3, 'New GPU Benchmarks', 'The latest cards are finally out. The performance jump is massive.', '2026-03-01 15:30:00'),
(4, 6, 'Why I love Rust', 'Memory safety without a garbage collector is just beautiful.', '2026-03-05 09:00:00'),
(5, 4, 'Looking for my first job', 'Any tips for junior devs in this market? It feels pretty tough.', '2026-03-15 11:00:00'),
(6, 5, 'Check out my Portfolio', 'I just finished my personal site using Svelte. Feedback appreciated!', '2026-03-20 14:20:00'),
(7, 2, 'State of CSS', 'Are we still using Tailwind or moving back to CSS modules?', '2026-04-01 10:00:00'),
(8, 3, 'Mechanical Keyboards', 'Just got a new 60% board. My typing speed increased instantly.', '2026-04-05 16:00:00'),
(9, 1, 'Forum Update v1.1', 'We added reactions and dark mode today!', '2026-04-10 13:00:00');

INSERT INTO post_categories (post_id, category_id, created_at) VALUES
(1, 1, '2026-01-01 10:05:00'),
(2, 2, '2026-02-10 12:00:00'),
(3, 3, '2026-03-01 15:30:00'),
(4, 2, '2026-03-05 09:00:00'),
(5, 5, '2026-03-15 11:00:00'),
(6, 4, '2026-03-20 14:20:00'),
(6, 2, '2026-03-20 14:25:00'),
(7, 2, '2026-04-01 10:00:00'),
(8, 3, '2026-04-05 16:00:00'),
(9, 1, '2026-04-10 13:00:00');

INSERT INTO comments (id, post_id, user_id, parent_id, content, created_at) VALUES
(1, 2, 3, NULL, 'I am still using React. The ecosystem is just too big to leave.', '2026-02-10 12:30:00'),
(2, 2, 2, 1, 'Fair point, but Vue 4 makes the DX so much better.', '2026-02-10 13:00:00'),
(3, 2, 6, NULL, 'Have you guys tried Leptos? It is incredibly fast.', '2026-02-11 08:00:00'),
(4, 4, 2, NULL, 'Rust is great but the compile times kill me.', '2026-03-05 10:00:00'),
(5, 4, 6, 4, 'True, but the safety guarantees save hours of debugging later!', '2026-03-05 10:15:00'),
(6, 5, 3, NULL, 'Focus on networking. Most jobs are found through referrals.', '2026-03-15 12:00:00'),
(7, 5, 4, 6, 'Thanks, I will try to attend more local meetups.', '2026-03-15 13:00:00'),
(8, 8, 5, NULL, 'Which switches did you get? Linears or Tactiles?', '2026-04-05 17:00:00'),
(9, 8, 3, 8, 'Went with Gateron Yellows. Very smooth.', '2026-04-05 17:30:00');

INSERT INTO reactions (id, user_id, post_id, comment_id, reaction_type, created_at) VALUES
(1, 2, 1, NULL, 1, '2026-01-01 10:10:00'),
(2, 3, 1, NULL, 1, '2026-01-01 10:15:00'),
(3, 1, NULL, 1, 1, '2026-02-10 12:40:00'),
(4, 4, 2, NULL, -1, '2026-02-10 14:00:00'),
(5, 5, NULL, 5, 1, '2026-03-05 11:00:00'),
(6, 6, 3, NULL, 1, '2026-03-01 16:00:00'),
(7, 2, NULL, 8, 1, '2026-04-05 18:00:00');

INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES
('uuid-1111-2222', 1, '2026-04-20 08:00:00', '2026-04-20 09:00:00'),
('uuid-3333-4444', 2, '2026-04-20 10:00:00', '2026-04-21 10:00:00'),
('uuid-5555-6666', 3, '2026-04-20 11:00:00', '2026-04-20 23:00:00'),
('uuid-7777-8888', 6, '2026-04-19 12:00:00', '2026-04-19 20:00:00');

-- additional seed data here

-- Adding missing comments to posts 3, 6, 7, 9
INSERT INTO comments (id, post_id, user_id, parent_id, content, created_at) VALUES
(10, 3, 1, NULL, 'The power draw on these new cards is insane though.', '2026-03-01 16:00:00'),
(11, 3, 5, 10, 'True, you basically need a 1000W PSU now.', '2026-03-01 16:45:00'),
(12, 3, 4, NULL, 'Still waiting for the mid-range benchmarks.', '2026-03-02 09:00:00'),
(13, 6, 1, NULL, 'Love the typography choice here! Very clean.', '2026-03-20 15:00:00'),
(14, 6, 3, NULL, 'The animations are a bit heavy on mobile.', '2026-03-20 15:30:00'),
(15, 6, 5, 14, 'I noticed that too, maybe try CSS transitions instead?', '2026-03-20 16:15:00'),
(16, 7, 4, NULL, 'Tailwind for prototypes, CSS modules for scale.', '2026-04-01 11:00:00'),
(17, 7, 6, NULL, 'I actually went back to vanilla CSS with the new nesting support.', '2026-04-01 12:30:00'),
(18, 9, 2, NULL, 'Dark mode looks amazing, thanks for the hard work!', '2026-04-10 14:00:00'),
(19, 9, 5, NULL, 'Found a small bug in the reaction counts on mobile.', '2026-04-10 15:20:00');

-- Strengthening post-category coverage (every category used at least twice)
INSERT INTO post_categories (post_id, category_id, created_at) VALUES
(2, 4, '2026-02-10 12:05:00'),
(3, 1, '2026-03-01 15:35:00'),
(4, 3, '2026-03-05 09:05:00'),
(5, 1, '2026-03-15 11:05:00'),
(7, 4, '2026-04-01 10:05:00'),
(8, 5, '2026-04-05 16:05:00');

-- Adding more reactions (mix of likes/dislikes)
INSERT INTO reactions (id, user_id, post_id, comment_id, reaction_type, created_at) VALUES
(8, 4, 1, NULL, 1, '2026-01-01 11:00:00'),
(9, 6, NULL, 10, -1, '2026-03-01 17:00:00'),
(10, 2, 4, NULL, 1, '2026-03-05 12:00:00'),
(11, 1, NULL, 13, 1, '2026-03-20 15:10:00'),
(12, 3, NULL, 14, -1, '2026-03-20 16:00:00'),
(13, 5, 7, NULL, 1, '2026-04-01 13:00:00'),
(14, 6, NULL, 18, 1, '2026-04-10 14:30:00'),
(15, 1, 8, NULL, 1, '2026-04-06 09:00:00');

-- Additional sessions (Active and Expired)
INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES
('uuid-9999-0000', 4, '2026-04-20 07:00:00', '2026-04-20 08:30:00'),
('uuid-abab-cded', 5, '2026-04-20 11:00:00', '2026-04-21 11:00:00'),
('uuid-fefe-ghgh', 1, '2026-04-20 11:05:00', '2026-04-20 14:00:00');