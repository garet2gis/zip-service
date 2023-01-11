# ZIP Service
Сервис для работы с ZIP архивами

## Запуск

Нужен предустановленный docker и утилита make

1. Создаем образ
```
make build
```
2. Запускаем контейнер с ограничением RAM в 512 Мб (боллее подробно команды можно посмотреть в файле Makefile)
```
make run
```

## API

Более подробно можно посмотреть, а также протестировать в Swagger-е 
по маршруту <b>/swagger</b>

<img width="621" alt="image" src="https://user-images.githubusercontent.com/42912280/211901479-308ff62e-5222-4404-b44d-2807c87df96e.png">


Все манипуляции с файлами происходят в директории <b>root</b>.
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

В результате выполнения запроса получаем zip-архив с именем <b>archive.zip</b> следующей структуры (2 картинки в папке myzip):

<img width="699" alt="image" src="https://user-images.githubusercontent.com/42912280/211898212-3a29f3ff-8330-4751-9b87-0aa8b23b6995.png">

Реализация использует буффер, размер которого можно задать. Таким образом, файл не загружается целиком в память, а только его часть. Было успешно протестировано на файлах размером в несколько Гб.

### POST /upload

Загружает zip-архивы на сервер в формате multipart-formdata
и распаковывает его в директорию <b>root/uploads</b>.

Пример CURL запроса
```
curl -X 'POST' \
  'http://0.0.0.0:8080/upload/' \
  -H 'accept: application/json' \
  -H 'Content-Type: multipart/form-data' \
  -F 'user_id=@archive.zip;type=application/zip'
```

В положительном результате получаем код ответа 204 (No content)
и следующую структуру файлов: 

<img width="239" alt="image" src="https://user-images.githubusercontent.com/42912280/211903695-8d1d87f7-0724-402b-9831-438c7362377f.png">

Реализация использует функцию ParseMultipartForm, которая нарезает поля multipart-formdata
на куски заданного размера, в моем случае около 4 Мб и загружает только их в RAM, остальные части сохраняются на диске
во временных файлах.

После чего загруженные zip-архивы распаковываются в отдельных горутинах в папку <b>uploads</b>.

Было также успешно протестировано на файлах размером в несколько Гб.

### GET /
Листинг файлов папки <b>root</b>:

<img width="752" alt="image" src="https://user-images.githubusercontent.com/42912280/211902405-723c407d-19a9-4ba2-a652-023fc607c5c8.png">

Содержимое папки <b>uploads</b>:

<img width="752" alt="image" src="https://user-images.githubusercontent.com/42912280/211902507-8c2a521d-2e8a-46bb-bfd7-f1ffe18163a2.png">

Содержимое папки <b>myzip</b>:

<img width="752" alt="image" src="https://user-images.githubusercontent.com/42912280/211902562-e446c6df-69c7-4794-929a-6b0900edc8e1.png">






