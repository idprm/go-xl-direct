INSERT INTO services
(category, code, name, price, program_id, sid_optin, sid_mt, renewal_day, trial_day, url_telco, url_portal, url_callback, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback, url_postback_billable)
VALUES
('DIGMAGZ', 'DIGMAGZ', 'DIGMAGZ', 2000, 'REGDIGMAGZ', 'REGDIGMAGZ', '00087197', 2, 0, 'https://api.digitalcore.telkomsel.com', 'https://tsel.digmagz.com', 'https://tsel.digmagz.com', 'https://tsel.digmagz.com/api/subscription/subscribe', 'https://tsel.digmagz.com/api/subscription/unsubscribe', 'https://tsel.digmagz.com/api/subscription/renewal', 'https://kbtools.net/pass-tsel.php', 'https://kbtools.net/pass-tsel.php');


INSERT INTO services
(category, code, name, price, program_id, sid_optin, sid_mt, renewal_day, trial_day, url_telco, url_portal, url_callback, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback, url_postback_billable)
VALUES
('GAMESIK', 'GAMESIK', 'GAMESIK', 2000, 'GAMESIK', 'GAMGENRKBGAMESIK_Subs', 'GAMGENRKBGAMESIK_Subs', 2, 0, 'https://api.digitalcore.telkomsel.com', 'https://tsel.gamesik.mobi', 'https://tsel.gamesik.mobi', 'https://tsel.gamesik.mobi/api/notification/subscribe', 'https://tsel.gamesik.mobi/api/notification/unsubscribe', 'https://tsel.gamesik.mobi/api/notification/renewal', 'https://kbtools.net/kb-tsel.php', 'https://kbtools.net/kb-tsel.php');



INSERT INTO contents
(service_id, name, value, tid)
VALUES
(1, 'FIRSTPUSH', 'Rp2220 DIGMAGZ Ada Artikel seru nih disini, klik: https://tsel.digmagz.com (Berlaku Tarif Internet) PIN : @pin CS: http://bit.ly/3Ly8Seq', '2220'),
(1, 'RENEWAL', 'Rp2220 DIGMAGZ Ada Artikel seru nih disini, klik: https://tsel.digmagz.com (Berlaku Tarif Internet) PIN : @pin CS: http://bit.ly/3Ly8Seq', '2220');


INSERT INTO contents
(service_id, name, value, tid)
VALUES
(1, 'FIRSTPUSH', 'Rp2220 GAMESIK Games bisa kamu mainkan, klik: https://tsel.gamesik.mobi (Berlaku Tarif Internet) PIN : @pin CS: http://bit.ly/CS-KB-Tsel', '2220'),
(1, 'RENEWAL', 'Rp2220 GAMESIK Games bisa kamu mainkan, klik: https://tsel.gamesik.mobi (Berlaku Tarif Internet) PIN : @pin CS: http://bit.ly/CS-KB-Tsel', '2220');


INSERT INTO schedules
(id, name, publish_at, unlocked_at, is_unlocked)
VALUES
(1, 'RENEWAL', NOW(), NOW(), false),
(2, 'RETRY_INSUFF', NOW(), NOW(), false),
(3, 'RETRY_INSUFF', NOW(), NOW(), false),
(4, 'RETRY_INSUFF', NOW(), NOW(), false),
(5, 'RETRY_INSUFF', NOW(), NOW(), false);