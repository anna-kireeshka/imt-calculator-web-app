# env

```
APP_PORT=8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_NAME=weightsecret
DB_USER=user
DB_PASSWORD=weightsecret

DATABASE_URL=postgres://user:weightsecret@localhost:5432/weightsecret
TELEGRAM_API_TOKEN=подставить с прода
```

История


weight float64 json wight
IMT float 64 json imt
Goal float64 json:"goal"
CreatedAt time.Time    `json:"-"`
UpdatedAt sql.NullTime `json:"-"`


Nutrions

Fiber
Fats
carbohydrates
BMRTotal
BMRCurrent



1 экран

Привет, для расчета вашего идеального веса и калорий нам нужно получить нескольк ваших ханных
Кнопка начать

2 экран 

Рост

3 экран 

вес

4 экран 

Дата рождения

Примечание - необходимо для расчета базого кол-ва калорий

5 экран
Пол

6 Цель
 - похудение
 - набор
 - удержание

экран с IMT и Каллориями объеденить

1 блок - информация о IMT


2 блок 
- информация о том сколько нужно потреблять калорий п
- информация о рекомендуемом БЖУ
- Кнопка изменить данные - вес, рост, активность,  цель

Страницв профиля

- ваш вес, при нажатии открывается экран с историей веса и возможностью обновить данные
- IMT меняется в рантайме
- Цель по весу, по дефолту не казана, но можно нажать на кнопку и откроется мини форма с водом желаемого веса
- Объемы тела при желании
История

- Дата
- Калоий в день
- Вес в сравнение с целью если есть
- Изменения в объемах если есть такая информация


POST

body {
    gender
    date_of_birth
    height
    weight
    goal
    activity_type
}

Insert inito user gender, date_of_birth

Insert into m height,  weight,  goal, activity_type



GET

BMR

response {
    IMT
    weight
    ccal
    macros: {
        fiber
        fat
        carbohydrates
    }
}

