-- Admin seed, pw: AdminPassword123!
INSERT INTO admin (
    username, email, password_hash, is_active, created_at, updated_at
) VALUES (
    'superadmin',
    'admin@example.com',
    '$2b$12$by3Payo5txAFt.I5BaVHUO7D9SctkP5buIkTdcIKdERHe6/EOYnh.',
    TRUE,
    NOW(),
    NOW()
);

-- Client seed
INSERT INTO client (
    phone, email, username, is_active, created_at, updated_at
) VALUES (
    '+85362818821',
    'client@example.com',
    'clientdemo',
    TRUE,
    NOW(),
    NOW()
);