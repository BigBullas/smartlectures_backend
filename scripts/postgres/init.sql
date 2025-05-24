CREATE TABLE IF NOT EXISTS users
(
    user_id  BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email    VARCHAR(300) UNIQUE NOT NULL,
    password BYTEA               NOT NULL,
    name     VARCHAR(128) DEFAULT '',
    surname  VARCHAR(128) DEFAULT ''
);


CREATE TABLE IF NOT EXISTS dirs
(
    dir_id     BIGSERIAL PRIMARY KEY,
    name       VARCHAR(128) NOT NULL,
    user_id    BIGINT REFERENCES users (user_id) NOT NULL,
    repeated_num BIGINT DEFAULT 0,
    parent_dir BIGINT REFERENCES dirs (dir_id) ON DELETE CASCADE DEFAULT NULL,
    icon_url VARCHAR(512) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS notes
(
    note_id     BIGSERIAL PRIMARY KEY,
    name        VARCHAR(256) NOT NULL,
    body        TEXT         NOT NULL           DEFAULT '' NOT NULL,
    created_at  TIMESTAMP                       DEFAULT NOW() NOT NULL,
    last_update TIMESTAMP                       DEFAULT NOW() NOT NULL,
    parent_dir  BIGINT REFERENCES dirs (dir_id)  ON DELETE CASCADE DEFAULT NULL,
    repeated_num BIGINT NOT NULL DEFAULT 0,
    user_id     BIGINT REFERENCES users (user_id) NOT NULL
);


CREATE TABLE IF NOT EXISTS snippets
(
    snippet_id  BIGSERIAL PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    description TEXT DEFAULT '' NOT NULL, 
    body        TEXT         NOT NULL DEFAULT '',
    user_id     BIGINT REFERENCES users (user_id)  ON DELETE CASCADE DEFAULT NULL 
);


-- TRIGGERS
create or replace function update_note_repeated_num()
returns trigger
as
    $$
BEGIN
    NEW.repeated_num = (SELECT COUNT(*) FROM notes WHERE name = NEW.name AND user_id = NEW.user_id AND parent_dir = NEW.parent_dir);
    RETURN NEW;
END;
$$
language plpgsql
;


CREATE TRIGGER update_note_repeated_num_trigger
BEFORE UPDATE ON notes
FOR EACH ROW
EXECUTE FUNCTION update_note_repeated_num();

create or replace function update_dir_repeated_num()
returns trigger
as
    $$
BEGIN
    NEW.repeated_num = (SELECT COUNT(*) FROM dirs WHERE name = NEW.name AND user_id = NEW.user_id AND parent_dir = NEW.parent_dir);
    RETURN NEW;
END;
$$
language plpgsql
;


CREATE TRIGGER update_dir_repeated_num_trigger
BEFORE UPDATE ON dirs
FOR EACH ROW
EXECUTE FUNCTION update_dir_repeated_num();

create or replace function create_default_note()
returns trigger
as
    $$
BEGIN
    INSERT INTO notes(name, user_id, repeated_num, parent_dir, body)
    VALUES ('Введение в приложение', NEW.user_id, 0, NULL,
      '
# Руководство по функционалу\n\n

 ## Содержание\n
1. Введение\n
2. Создание заметок\n
    - Создание заметки с нуля\n
    - Загрузка изображения\n
    - Импорт файла\n
3. Редактирование заметок\n
    - Основы\n
    - Добавление файлов в документ\n\n

# 1. Введение\n
## EasyTex - приложение для редактирования конспектов. Вот его основной функционал:
1. Распознавание рукописных конспектов и изображений, находящихся в нередактируемых форматах, а именно:
    - преобразование фотографии рукописного текста в Markdown-формат
    - импортирование PDF- и Markdown-файлов
    - папочная иерархия для удобной систематизации конспектов

2. Сниппеты для упрощенной работы с языком LaTeX


> Примечание: чтобы полностью раскрыть функционал нашего приложения, стоит знать, как работать с языками Markdown и LaTeX. Подробнее про них можно прочитать здесь:
>- Markdown - [тыц](https://doka.guide/tools/markdown/)
>- LaTeX - [тыц](https://www.overleaf.com/learn/latex/Learn_LaTeX_in_30_minutes)


# 2. Создание заметок
Создание любой заметки начинается с перехода в меню создания заметок с помощью нажатия на иконку "плюс" в верхней части экрана: 

![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/59f4670e-4328-4e30-8a0c-20bd9f4785ed)

или через клик ПКМ по папке в панели навигации: 

![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/7e92f5c1-4ac6-4401-92b6-f2c7d98f1371)

После этого вы перейдете в поле выбора способа создания документа:
 ![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/3c71d16f-7df4-4226-a422-2ec65487368b)
## Создание заметки с нуля
Для создания заметки с нуля используйте соответствующую кнопку, при нажатии на которую совершится автоматический переход в редактор.

## Загрузка изображения
Для распознавания фотографии конспекта достаточно нажатия на среднюю кнопку и перетащить нужный файл в поле.

> Обратите внимание: до нажатия на кнопку "Загрузка фотографий" поле для перетаскивания файла будет неактивно!

## Сохранение заметки
Для сохранения заметки **обязательно** нажать на кнопку "Сохранить" в правом верхнем углу приложения  
![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/65a3e100-cea1-4302-99ce-f1acf3333e17)

## Импорт файла
Для импорта файла сценарий аналогичен загрузке изображения за исколючением первого шага, на котором нужно нажать кнопку "Импорт файлов".

> В настоящее время доступен импорт из PDF и Markdown.

# 3. Редактирование заметок
## Основы
> Редактирование основано на использовании таких языков как Markdown и LaTex. Быстро ознакомиться с ними можно по ссылкам в начале мануала.

Перейти к редактированию заметки можно, нажав на нее в панели навигации:  
![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/22154358-9af8-4705-9bba-f30681044407)

Откроется редактор. Он имеет три представления, переключаться между которыми можно с помощью пиктограмм в правом верхнем углу окна редактора:  
![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/e0cf56ca-be67-4039-8b17-014f65ebbc0a)

Для редактирования кода рекомендуем использовать режим с превью итогового файла, который выглядит следующим образом:  
![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/6ec2accf-98e2-42b4-9fb3-379fdaf9c417)

## Добавление файлов в документ
Если у вас есть какое-то изображение, которое вам бы хотелось прикрепить к текущему документу, в окне редактирования при нажатии ПКМ откроется диалоговое окошко, которое предложит вам сделать это несколькими способами:
![image](https://hb.ru-msk.vkcs.cloud/easytex/attachments/3f6cefe2-a3fc-423d-9e88-08806bb4ac6d)
> Вставить фото - вставляет файл без изменений - просто как вложение

> Преобразовать - можно добавить распознанное фото текста или формулы. Они поместятся в документ как кусок текста.

> Добавить фото конспекта - распознает смешанные фотографии, т.е. те, на которых есть и текст, и формулы');

RETURN NEW; -- This was missing in the original

END;
$$
language plpgsql
;

-- note_id     BIGSERIAL PRIMARY KEY,
--     name        VARCHAR(256) NOT NULL,
--     body        TEXT         NOT NULL           DEFAULT '' NOT NULL,
--     created_at  TIMESTAMP                       DEFAULT NOW() NOT NULL,
--     last_update TIMESTAMP                       DEFAULT NOW() NOT NULL,
--     parent_dir  BIGINT REFERENCES dirs (dir_id)  ON DELETE CASCADE DEFAULT NULL,
--     repeated_num BIGINT NOT NULL DEFAULT 0,
--     user_id     BIGINT REFERENCES users (user_id) NOT NULL
-- );

CREATE TRIGGER create_default_note_trigger
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION create_default_note();
