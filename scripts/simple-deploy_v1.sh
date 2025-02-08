#!/bin/bash

# 🔹 Загрузка переменных из .env
if [ -f .env ]; then
    echo "🚀 Загрузка переменных из .env файла."
    export $(cat .env | grep -v '#' | xargs)
fi

# Проверка наличия всех переменных
if [ -z "$SSH_USER" ] || [ -z "$SSH_HOST" ] || [ -z "$SSH_PORT" ] || [ -z "$DOCKERHUB_USER" ] || [ -z "$IMAGE_NAME" ] || [ -z "$TAG" ] || [ -z "$CONTAINER_NAME" ]; then
    echo "❗ Ошибка: Отсутствуют необходимые переменные в .env файле."
    exit 1
fi

# 🔹 1. Сборка Docker-образа
echo "🚀 Сборка Docker-образа..."
docker build --platform linux/amd64 -t $DOCKERHUB_USER/$IMAGE_NAME:$TAG -f docker/golang/Dockerfile .
#docker build --platform linux/arm64 -t $DOCKERHUB_USER/$IMAGE_NAME:$TAG -f docker/golang/Dockerfile .

# 🔹 2. Логин в Docker Hub
echo "🔑 Вход в Docker Hub..."
echo "Введите пароль для Docker Hub:"
docker login docker.io -u $DOCKERHUB_USER

# 🔹 3. Пушим образ в Docker Hub
echo "📤 Пушим образ в Docker Hub..."
docker push $DOCKERHUB_USER/$IMAGE_NAME:$TAG

# 🔹 4. Деплой на сервер по SSH
echo "🔗 Подключение к серверу по SSH и обновление контейнера..."
ssh -tt -p $SSH_PORT $SSH_USER@$SSH_HOST << EOF
    echo "🛑 Остановка старого контейнера (если он есть)..."
    docker stop $CONTAINER_NAME || true
    docker rm $CONTAINER_NAME || true

    echo "📥 Pull образа..."
    docker pull $DOCKERHUB_USER/$IMAGE_NAME:$TAG

    echo "🚀 Запуск нового контейнера..."
    docker run -d --rm --name $CONTAINER_NAME -p $DOCKER_RSS_PORT:$DOCKER_RSS_INTERNAL_PORT $DOCKERHUB_USER/$IMAGE_NAME:$TAG
    exit
EOF

echo "✅ Деплой завершен!"
