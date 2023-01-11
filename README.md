# ZIP Service
Сервис для работы с ZIP архивами

## Запуск

Нужен предустановленный docker и утилита make

1. Создаем образ
```
make build
```
2. Запускаем контейнер с ограничением RAM в 512 Мб
```
make run
```

## API

Более подробно можно посмотреть, а также протестировать в Swagger-е 
по маршруту <b>/swagger</b>

Все манипуляции с файлами происходят в директории <b>root</b>, название которой
можно изменить в .env файле
### POST /download

Скачивает ZIP архив из файлов загруженных на сервер ранее

Пример CURL запроса
```
curl -X 'POST' \
  'http://0.0.0.0:8080/download/' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "files": [
    {
      "path": "images/photo.png",
      "zip_path": "myzip/photo1.png"
    },
    {
      "path": "upload/archive/first4.png",
      "zip_path": "myzip/photo2.png"
    }
  ]
}'
```
В теле указываю массив из файлов, где поле <b>path</b> - местоположение файла на хосте,
а поле <b>zip_path</b> желаемое местонахождение файла в архиве

В результате выполнения запроса получаем zip-архив с именем <b>archive.zip</b>

### POST /upload

Загружает zip-архивы на сервер в формате multipart-formdata
и распаковывает его в папку <b>upload</b> (название можно изменить в .env)

Пример CURL запроса
```
curl -X 'POST' \
  'http://0.0.0.0:8080/upload/' \
  -H 'accept: application/json' \
  -H 'Content-Type: multipart/form-data' \
  -F 'user_id=@archive_2.zip;type=application/zip'
```

В положительном результате получаем код ответ 204 (No content)
