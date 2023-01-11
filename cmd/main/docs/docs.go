// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/download/": {
            "post": {
                "tags": [
                    "ZIP"
                ],
                "summary": "Скачивание желаемых файлов в форме zip-архива",
                "operationId": "download-zip",
                "parameters": [
                    {
                        "description": "Zip Descriptor",
                        "name": "user_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ZipDescriptor"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ZIP file",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/AppError"
                        }
                    }
                }
            }
        },
        "/upload/": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "ZIP"
                ],
                "summary": "Загрузка ZIP файла и его разархивация",
                "operationId": "upload-zip",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "file"
                        },
                        "description": "Zip files to upload",
                        "name": "form_data",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/AppError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "AppError": {
            "type": "object",
            "properties": {
                "developer_message": {
                    "description": "Сообщение для разработчика",
                    "type": "string"
                },
                "message": {
                    "description": "Сообщение",
                    "type": "string"
                }
            }
        },
        "ZipDescriptor": {
            "type": "object",
            "properties": {
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.FileEntry"
                    }
                }
            }
        },
        "dto.FileEntry": {
            "type": "object",
            "properties": {
                "path": {
                    "description": "Путь до файла на хосте",
                    "type": "string"
                },
                "zip_path": {
                    "description": "Желаемый путь до файла в zip архиве",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Zip service API documentation",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
