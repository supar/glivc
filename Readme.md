# GitLab interface to view commits (glivc)

## Зачем все это нужно

Актуальность данного сервиса возникла из-за ряда ограничений в API GitLab: интерфейс не позволяет ограничить доступ только для просмотра истории коммита в выбранных проектах, а так же слабый фильтр коммитов

## API

### Получить список проектов

GET http://localhost:9999/project

### Информация по одному проекту

GET http://localhost:9999/project/ID

Где ID числовой идентификатор проекта

### Получить список веток

GET http://localhost:9999/project/ID/branches

Где ID числовой идентификатор проекта

### Информация по одной ветке

GET http://localhost:9999/project/ID/branches/BRANCH_NAME

Где ID числовой идентификатор проекта, BRANCH_NAME - имя ветки

### Список коммитов

GET http://localhost:9999/project/ID/branch/commits/BRANCH_NAME

Где ID - числовой идентификатор проекта, BRANCH_NAME - имя ветки

#### Дополнительные параметры запроса

- page - номер страницы (default:1)
- limit - количество записей на странице (default:unlimit)
- author - идентификатор автора. Проверяется вхождение значение этого параметра в подписи автора к коммиту: имя, email. В качестве поиска лучше использовать email. Этот параметр можно передать несколько раз http://localhost:9999/project/ID/branch/commits/BRANCH_NAME?author=user1&author=user2
- merge - исключить/включить коммиты Merge Request (default:0)
- msg - поиск по комментарию коммита