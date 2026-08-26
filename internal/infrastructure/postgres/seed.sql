INSERT INTO bullets (id, type, signifier, content, created_at, updated_at) VALUES
  (gen_random_uuid(), 'task', 'open', 'Создайте свою первую задачу – нажмите на поле ниже', now() - interval '1 hour', now()),
  (gen_random_uuid(), 'task', 'completed', 'Завершённые задачи', now() - interval '1 day', now()),
  (gen_random_uuid(), 'task', 'open', 'Незавершённые задачи', now() - interval '1 day', now()),

  (gen_random_uuid(), 'event', 'open', 'Так выглядят события', now() - interval '1 day', now()),
  (gen_random_uuid(), 'note', 'open', 'Заметки, мысли или идеи', now() - interval '1 day', now());

