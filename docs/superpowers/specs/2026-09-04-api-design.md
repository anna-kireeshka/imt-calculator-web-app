# API для Telegram Mini App «Weight Secret»

Дата: 2026-09-04
Статус: утверждён, ждёт плана реализации

## 1. Цель

Дать мини-аппу внутри Telegram возможность сохранять замеры роста и веса,
показывать историю и получать производные числа — ИМТ, границы здорового веса
и суточную норму калорий. Сейчас страница считает ИМТ в браузере и держит
последние цифры в `localStorage`; после этой работы данные живут на сервере и
привязаны к пользователю Telegram.

## 2. Решения, зафиксированные до начала

| Вопрос | Решение |
|---|---|
| Первый клиент | Мини-апп в Telegram. Бот и внешние клиенты — вне объёма |
| Объём v1 | Профиль и замеры |
| Зависимости | Только стандартная библиотека плюс уже подключённый `pgx` |
| Раскладка | Два слоя: хендлер → репозиторий; формулы — чистые функции в `domain` |

Третий слой (`service`) сознательно не вводится: на текущем объёме он состоял бы
из однострочных прослоек. Переезд из двух слоёв в три дешёвый, обратный — нет.

## 3. Авторизация

Единственный источник личности пользователя — подпись `initData`, которую
Telegram выдаёт мини-аппу. Ничему остальному, что пришло от клиента, доверять
нельзя: телу запроса, заголовкам, параметрам пути.

Алгоритм проверки:

1. Секретный ключ: `HMAC_SHA256(key: "WebAppData", data: bot_token)`.
2. `data_check_string`: все пары `ключ=значение` из `initData`, кроме `hash`,
   отсортированные по ключу и соединённые через `\n`.
3. Вычислить `HMAC_SHA256(key: секрет, data: data_check_string)` и сравнить с
   `hash` при помощи `hmac.Equal` — не через `==` и не через `strings.Compare`.
4. Отвергнуть, если `auth_date` старше 24 часов.

Транспорт: заголовок `Authorization: tma <initData>`.

Мидлварь `Auth`:

- проверяет подпись, при неудаче отвечает `401 unauthenticated`;
- выполняет `INSERT INTO users (telegram_user_id, name) VALUES ($1, $2)
  ON CONFLICT (telegram_user_id) DO UPDATE SET name = EXCLUDED.name
  RETURNING user_id, ...` — один запрос, идемпотентно, заодно освежает имя;
- кладёт пользователя в `r.Context()` под неэкспортированным ключом.

Токен бота приходит в `server.New` из `main` (переменная `TELEGRAM_API_TOKEN`),
пакет `server` сам в окружение не ходит.

## 4. HTTP-контракт

Все ответы — `application/json; charset=utf-8`. Все пути под `/api` требуют
авторизации. Путь `/` остаётся отдачей лендинга.

### GET /api/me

`200 OK`

```json
{
  "user": {
    "name": "Анна",
    "gender": "woman",
    "birth_date": "1996-04-12",
    "activity_coefficient": 1.375,
    "profile_complete": true
  },
  "last_measurement": {
    "measurement_id": 42,
    "height_cm": 170,
    "weight_kg": 64.5,
    "imt": 22.3,
    "measured_at": "2026-09-04T09:12:00Z"
  },
  "derived": {
    "healthy_weight_kg": { "min": 54, "max": 72 },
    "daily_calories": 1920
  }
}
```

- `profile_complete` — истина, когда заполнены `gender`, `birth_date` и
  `activity_coefficient`.
- Нет ни одного замера: и `last_measurement`, и `derived` равны `null` целиком —
  не объектами с пустыми полями.
- Замер есть, профиль неполон: `healthy_weight_kg` заполнен (он зависит только
  от роста), `daily_calories` равен `null`.

Отдельный `/api/summary` не нужен: производные числа возвращаются здесь.

### PUT /api/me

Полная замена профиля. Тело:

```json
{
  "name": "Анна",
  "gender": "woman",
  "birth_date": "1996-04-12",
  "activity_coefficient": 1.375
}
```

- `name` может быть пустой строкой, остальные три поля обязательны.
- `gender` — строго `man` или `woman`.
- `birth_date` — `YYYY-MM-DD`, не в будущем, возраст в пределах 10–120 лет.
- `activity_coefficient` — от 1.0 до 2.5.

Ответ `200 OK` — объект `user` в том же виде, что в `GET /api/me`.

### POST /api/measurements

```json
{ "height_cm": 170, "weight_kg": 64.5 }
```

- `height_cm` — от 100 до 250, `weight_kg` — от 20 до 400.
- `measured_at` клиент не задаёт: время ставит база.

Ответ `201 Created` — созданный замер целиком, включая вычисленный базой `imt`
(через `RETURNING`).

### GET /api/measurements

Параметры: `limit` (по умолчанию 50, максимум 200) и `before` — курсор.

```json
{
  "items": [ { "measurement_id": 42, "height_cm": 170, "weight_kg": 64.5,
               "imt": 22.3, "measured_at": "2026-09-04T09:12:00Z" } ],
  "next_before": 17
}
```

Порядок — `measurement_id DESC`. Так как в v1 клиент не задаёт `measured_at`,
порядок по идентификатору совпадает с хронологическим. Курсор `before` — это
`measurement_id`: выдаются записи со строго меньшим идентификатором.
`next_before` равен `null`, когда страница последняя.

## 5. Формат ошибок

Единственный формат наружу:

```json
{ "error": { "code": "invalid_argument", "message": "рост должен быть от 100 до 250 см", "field": "height_cm" } }
```

`field` присутствует только у ошибок валидации.

| code | HTTP | Когда |
|---|---|---|
| `unauthenticated` | 401 | нет заголовка, битая подпись, просроченный `auth_date` |
| `invalid_argument` | 400 | не разобралось тело, не прошла валидация |
| `not_found` | 404 | неизвестный путь |
| `method_not_allowed` | 405 | метод не подходит под маршрут |
| `payload_too_large` | 413 | тело больше лимита |
| `internal` | 500 | всё остальное |

Правила:

- перевод доменных ошибок в HTTP живёт в одном месте — `writeError`, а не
  расползается по хендлерам;
- наружу уходит общее сообщение, подробности пишутся в лог с идентификатором
  запроса;
- `log.Fatal` не встречается нигде, кроме `main`;
- тело читается через `http.MaxBytesReader` с лимитом 64 КБ;
- `json.Decoder` с `DisallowUnknownFields`, чтобы опечатка в имени поля не
  проходила молча.

Отсутствие замеров — не ошибка: `GET /api/me` отвечает `200` с `null`.

## 6. Структура пакетов

```
internal/telegram/initdata.go       Validate(initData, token) (TelegramUser, error)
internal/domain/user.go             User, Gender
internal/domain/measurement.go      Measurement
internal/domain/calories.go         BMI, HealthyWeightRange, DailyCalories — чистые функции
internal/database/pool.go           NewPool (уже написан)
internal/database/migrations.go     Migrate(ctx, pool) + //go:embed миграций
internal/database/users.go          UserRepo:        Upsert, Get, UpdateProfile
internal/database/measurements.go   MeasurementRepo: Create, ListByUser, LastWithProfile
internal/server/server.go           New(deps) *http.Server, Router()
internal/server/middleware.go       Auth, Recover, RequestLog
internal/server/json.go             decodeJSON, writeJSON, writeError
internal/server/handlers_me.go      GET/PUT /api/me
internal/server/handlers_measurements.go
database/migrations/0001_init.sql
```

Ключевое изменение против текущего кода: `CalcUserDailyCalorie` разделяется
надвое. `MeasurementRepo.LastWithProfile` одним `JOIN` достаёт последний замер
вместе с полом, датой рождения и коэффициентом активности;
`domain.DailyCalories(weightKg, heightCm, age, gender, activity)` считает по
Миффлину-Сан Жеору из готовых чисел. Формула перестаёт зависеть от Postgres и
покрывается обычными табличными тестами.

Интерфейсы репозиториев объявляются в пакете `server` — там, где они
используются, а не там, где реализуются.

Функции `CreateTableUser` и `CreateTableMeasurements` удаляются: их заменяют
миграции.

## 7. Схема и миграции

`CREATE TABLE IF NOT EXISTS` не меняет уже существующую таблицу — из-за этого
код и база в проекте уже разошлись. Вместо него: пронумерованные `.sql` в
`database/migrations/`, вшитые через `//go:embed`, и таблица учёта.

Раннер: читает файлы в лексикографическом порядке; для каждого проверяет запись
в `schema_migrations`; непринятые выполняет в транзакции вместе с вставкой
записи. Около пятидесяти строк, внешние зависимости не нужны.

Текущие таблицы созданы из промежуточных версий DDL и содержат мусор — в них
нет `activity_coefficient`, зато есть лишний `age`. Данных в них нет, поэтому
база пересоздаётся один раз через `docker compose down -v`, дальше только
миграции.

`0001_init.sql`:

```sql
CREATE TABLE users (
    user_id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    telegram_user_id     BIGINT       NOT NULL UNIQUE,
    name                 TEXT,
    birth_date           DATE,
    gender               TEXT         CHECK (gender IN ('man', 'woman')),
    activity_coefficient NUMERIC(4,3) CHECK (activity_coefficient BETWEEN 1 AND 2.5),
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE measurements (
    measurement_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id        BIGINT        NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    height_cm      NUMERIC(5,1)  NOT NULL CHECK (height_cm BETWEEN 100 AND 250),
    weight_kg      NUMERIC(5,2)  NOT NULL CHECK (weight_kg BETWEEN 20 AND 400),
    imt            NUMERIC(4,1)  GENERATED ALWAYS AS (weight_kg / ((height_cm / 100) ^ 2)) STORED,
    measured_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX measurements_user_id_idx ON measurements (user_id, measurement_id DESC);
```

Возраст не хранится: он выводится из `birth_date` запросом
`date_part('year', age(m.measured_at, u.birth_date))::int`, то есть считается на
дату замера, а не на «сейчас». Хранимый возраст протух бы через год.

`activity_coefficient` — именно `NUMERIC(4,3)`: при одном знаке после запятой
стандартный коэффициент 1.375 округляется до 1.4, что даёт ошибку в несколько
десятков килокалорий.

## 8. Тестирование

- `internal/telegram` — таблица случаев: валидная подпись, подделанная, без
  `hash`, просроченный `auth_date`, мусор вместо query-строки. Токен
  фиктивный, сеть не нужна.
- `internal/domain` — таблица на формулы: мужчина, женщина, границы возраста,
  неизвестный пол, границы диапазонов роста и веса.
- `internal/server` — `httptest` поверх роутера с подставными репозиториями:
  коды ответов, форматы ошибок, отказ без заголовка авторизации.
- `internal/database` — интеграционные тесты против Postgres из `docker
  compose`; пропускаются при `go test -short`.

## 9. Порядок реализации

Сквозным срезом, а не слоями: пока первый запрос не проходит насквозь, не
видно, сходятся ли куски друг с другом.

1. Миграции, раннер, `0001_init.sql`, пересоздание базы.
2. `internal/telegram/initdata.go` и тесты к нему.
3. `GET /api/me` целиком: роутер, мидлварь `Auth`, `json.go`, `UserRepo`,
   `MeasurementRepo.LastWithProfile`, `domain.DailyCalories`, тесты. Заодно
   перестаёт быть пустым `internal/server/server.go`, который сейчас ломает
   сборку.
4. `PUT /api/me`, `POST /api/measurements`, `GET /api/measurements` по образцу.
5. Перевод лендинга с `localStorage` на API.

## 10. Вне объёма v1

Удаление и редактирование замеров, экспорт данных, вебхуки, ограничение частоты
запросов, метрики, публичное API с собственными учётными записями,
локализация, вызовы этого API из хендлеров бота.
