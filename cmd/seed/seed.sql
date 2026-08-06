TRUNCATE TABLE bullets CASCADE;

INSERT INTO bullets (id, type, signifier, content, created_at, updated_at) VALUES
  (gen_random_uuid(), 'task', 'open', 'Доделать документацию по API', now() - interval '2 days', now() - interval '1 day'),
  (gen_random_uuid(), 'task', 'completed', 'Написать тесты для репозитория', now() - interval '5 days', now() - interval '4 days'),
  (gen_random_uuid(), 'task', 'cancelled', 'Переписать фронтенд на React', now() - interval '10 days', now() - interval '9 days'),
  (gen_random_uuid(), 'task', 'open', 'Настроить CI/CD для продакшена', now() - interval '1 day', now()),
  
  (gen_random_uuid(), 'event', 'scheduled', 'Встреча с командой в 15:00', now() + interval '2 hours', now()),
  (gen_random_uuid(), 'event', 'completed', 'Демо для заказчика прошло успешно', now() - interval '3 days', now() - interval '2 days'),
  (gen_random_uuid(), 'event', 'cancelled', 'Вебинар по Go отменён', now() - interval '1 day', now()),
  (gen_random_uuid(), 'event', 'scheduled', 'Релиз v2.0 запланирован на пятницу', now() + interval '3 days', now()),
  
  (gen_random_uuid(), 'note', 'open', 'Идея для нового фича-флага: A/B тестирование', now() - interval '4 hours', now()),
  (gen_random_uuid(), 'note', 'migrated', 'Перенесённые заметки из старого трекера', now() - interval '7 days', now() - interval '6 days'),
  (gen_random_uuid(), 'note', 'open', 'Список книг по архитектуре: Clean Code, Domain-Driven Design', now() - interval '12 hours', now());