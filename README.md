Микросервис для проведения тендеров

# Стек

- Go
- Pgx
- Chi
- Postgres (из задания)

# Зпуск
Сервис запускается на 8080 порту через docker файл, миграции лежат в папке migrates

Старался сделать сервис легко расширяемым и тестируемым. Но тесты написать не успел

# Мои решения
1) Не малая часть реализации была перенесена на базу данных (триггеры, функции) - считаю это позволило сделать код чище, чем могло быть без этого. Но, конечно, функционал базы не так удобно обслуживать и отлаживать
2) Было не совсем понятно: менять ли версию просто при изменении статуса - решил менять при любом изменении
3) Также при /bids/{bidId}/feedback показалось странным, что возврощается не feedback, а bid. Решил быть верным openapi.yml

# Хотелось бы исправить
1) Думаю можно значительно снизить количество кода, если грамотно вынести обработку ошибок и другие куски кода, которые схожи внешне, но разлчины в деталях - не успел сделать
2) А так же в слое сервисов что-то придумать с проверками на существование, не существование, ответственность, отсутствие ответственностие - тоже помогло бы сократить код
