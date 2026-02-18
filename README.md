# Проект - Менеджер задач

## Короткое описание

REST API для управления пользователями и задачами.
Реализована поддержка назначения нескольких исполнителей на одну задачу.

## Стек технологий

- Go 1.25
- PostgreSQL
- Chi (router)
- slog (logging)
- UUID (github.com/google/uuid)
- go-playground/validator

## Архитектура приложения

- domain - бизнес-сущности и ошибки
- usecase - бизнес-логика
- repository - работа с БД
- http-server/handler - HTTP слой
- lib - вспомогательные пакеты (logger, response)

## Схема БД

### users
- id (UUID)
- email
- name

### tasks
- id (UUID)
- title
- description
- status
- created_at

### user_task (many-to-many)
- user_id
- task_id

## Конечные точки API

### Users

POST /users

<img width="438" height="385" alt="image" src="https://github.com/user-attachments/assets/e44908e8-b6eb-4d8b-b778-104687d1107d" />

GET /users/{id} 

<img width="546" height="337" alt="image" src="https://github.com/user-attachments/assets/f86215f6-9000-4b71-ba92-0acac88874a4" />

DELETE /users/{id}

<img width="552" height="284" alt="image" src="https://github.com/user-attachments/assets/3f85ac59-0fbd-43c3-8829-84d2cbaa277f" />

### Tasks

POST /tasks 

<img width="460" height="498" alt="image" src="https://github.com/user-attachments/assets/b6309566-3233-46ad-ba49-cdf8624d5e80" />

GET /tasks/{id}

<img width="540" height="339" alt="image" src="https://github.com/user-attachments/assets/e84ba19a-8295-417d-9052-9570fc06d8a1" />

DELETE /tasks/{id}

<img width="561" height="233" alt="image" src="https://github.com/user-attachments/assets/a6cd3513-f48a-4225-8f78-0c82442106dc" />
