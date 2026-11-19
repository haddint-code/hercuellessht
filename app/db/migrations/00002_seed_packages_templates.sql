-- +goose Up
-- Seed packages, package_features and starter templates.

INSERT INTO packages (id, name, tagline, description, price_cents, duration_days, max_photos, max_videos, max_events, coming_soon, available, sort_order) VALUES
('silver', 'Silver', 'Simple & Elegant', 'Undangan digital siap pakai, tinggal isi data.', 15000000, 60, 8, 1, 2, FALSE, TRUE, 1),
('gold',   'Gold',   'Personalized & Premium', 'Template premium yang dapat disesuaikan dengan karakter pasangan.', 35000000, 90, 24, 3, 5, FALSE, TRUE, 2),
('platinum','Platinum','Coming Soon', 'Paket eksklusif dengan fitur lengkap. Segera hadir.', 75000000, 180, 60, 10, 10, TRUE, FALSE, 3);

INSERT INTO package_features (package_id,
    rsvp_enabled, guestbook_enabled, digital_envelope, whatsapp_button, background_music,
    countdown, guest_greeting, couple_story,
    custom_colors, custom_fonts, custom_layout, custom_cover, custom_background,
    custom_sections, custom_wording, custom_music, custom_url,
    multiple_events, map_preview, event_customization,
    live_rsvp, guest_list, qr_code, social_sharing_preview) VALUES
('silver',
    TRUE, TRUE, TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE,
    FALSE, FALSE, FALSE, FALSE, FALSE,
    FALSE, FALSE, FALSE, FALSE,
    FALSE, TRUE, FALSE,
    FALSE, FALSE, FALSE, FALSE),
('gold',
    TRUE, TRUE, TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE, TRUE),
('platinum',
    TRUE, TRUE, TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE,
    TRUE, TRUE, TRUE, TRUE);

INSERT INTO templates (id, name, category, description, thumbnail, package_id, is_active, config, sort_order) VALUES
('S01', 'Floral Elegant', 'floral', 'Nuansa bunga pastel yang lembut dan klasik.', '/static/img/templates/s01.jpg', 'silver', TRUE,
 '{"palette":["#8d6e63","#efeae2","#c9a27c"],"font_heading":"Alex Brush","font_body":"Plus Jakarta Sans","cover_style":"framed"}', 1),
('S02', 'Japanese Garden', 'japanese', 'Minimalisme taman Jepang dengan sakura dan garis bersih.', '/static/img/templates/s02.jpg', 'silver', TRUE,
 '{"palette":["#5b4636","#f5f1e8","#b3574f"],"font_heading":"Alex Brush","font_body":"Plus Jakarta Sans","cover_style":"torii"}', 2),
('S03', 'Minimalist Romance', 'minimalist', 'Putih bersih, tipografi elegan, ruang lega.', '/static/img/templates/s03.jpg', 'silver', TRUE,
 '{"palette":["#3f3f3f","#ffffff","#a8a29e"],"font_heading":"Cormorant Garamond","font_body":"Inter","cover_style":"split"}', 3),
('G01', 'Golden Kirana', 'luxury', 'Resmi dan megah dengan sentuhan emas premium.', '/static/img/templates/g01.jpg', 'gold', TRUE,
 '{"palette":["#1f1a17","#f3ead9","#c9a27c"],"font_heading":"Alex Brush","font_body":"Plus Jakarta Sans","cover_style":"fullbleed","layouts":["classic","modern"]}', 4),
('G02', 'Sakura Premium', 'japanese', 'Taman sakura versi premium dengan layout ganda.', '/static/img/templates/g02.jpg', 'gold', TRUE,
 '{"palette":["#4a3728","#faf3ea","#d8a0a0"],"font_heading":"Cormorant Garamond","font_body":"Plus Jakarta Sans","cover_style":"parallax","layouts":["classic","story-first"]}', 5);

-- +goose Down
DELETE FROM templates;
DELETE FROM package_features;
DELETE FROM packages;
