-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chapters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    "order" INTEGER DEFAULT 0,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lessons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    content TEXT,
    "order" INTEGER DEFAULT 0,
    chapter_id INTEGER REFERENCES chapters(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO courses (name, description) VALUES
('Golang для начинающих', 'Изучаем основы Go с нуля'),
('Microservices на Go', 'Продвинутый курс по архитектуре'),
('Java для начинающих', 'Изучаем основы Java с нуля'),
('Python для начинающих', 'Изучаем основы Python с нуля');

INSERT INTO chapters (course_id, name, description, "order") VALUES
(1, 'Введение', 'Обзор языка и подготовка среды', 1),
(1, 'Синтаксис', 'Базовые конструкции языка', 2),
(2, 'Docker и контейнеры', 'Упаковка Go приложений', 1),
(3, 'Java основы', 'Разберем как устроен язык программирования', 1),
(3, 'Переменные', 'Какие типы переменных используется', 2),
(4, 'Python основы', 'От объявления переменных до Django', 1);


INSERT INTO lessons (chapter_id, name, description, content, "order") VALUES
(1, 'Установка Go', 'Инструкция', 'Скачайте Go с официального сайта...', 1),
(2, 'Переменные и типы', 'Разбор типов', 'В Go строгая типизация...', 2),
(3, 'Dockerfile', 'Создание образа', 'Используйте multi-stage build...', 3),
(4, 'Объявление строки', 'Научимся выводить данные на экран', 'Для того чтобы вывести на экран текст необходимо...', 1),
(5, 'Тип переменных int', 'Где и для чего используется', 'Пропишите название переменных ипосле знака равно...', 1),
(5, 'Тип переменных bool', 'Поймем где нам нужен тип переменных bool', 'Когда нам нужно понять истенно ли значение...', 2),
(6, 'Строки', 'Особенности типа переменных строки', 'КОгда необходимо вывести текст мы используем...', 1);

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS courses;
