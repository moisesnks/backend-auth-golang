// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "moisesnks",
            "url": "https://github.com/moisesnks",
            "email": "moisesnks@utem.cl"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/change-password": {
            "post": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Actualiza la contraseña del usuario autenticado usando Firebase Authentication.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Actualizar contraseña",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Nueva contraseña del usuario",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.ChangePasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Contraseña actualizada correctamente",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "Datos de solicitud inválidos",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Token de autenticación no proporcionado o inválido",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Token de autenticación expirado",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/forgot-password": {
            "post": {
                "description": "Envía un correo electrónico para restablecer la contraseña si el usuario existe en Firebase Authentication.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Restablecer contraseña",
                "parameters": [
                    {
                        "description": "Correo electrónico del usuario",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.ForgotPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Correo de restablecimiento de contraseña enviado",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "Datos de solicitud inválidos",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Inicia sesión utilizando Firebase Authentication",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Inicia sesión de usuario",
                "parameters": [
                    {
                        "description": "Datos de inicio de sesión",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Inicio de sesión exitoso",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "Datos de solicitud inválidos o errores en la solicitud",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Credenciales incorrectas",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Registra un nuevo usuario utilizando Firebase Authentication y envía un correo de verificación.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Registrar usuario",
                "parameters": [
                    {
                        "description": "Datos de registro del usuario",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Respuesta exitosa al registrar usuario",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "Datos de solicitud inválidos",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "El correo electrónico ya está en uso",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/resend-code": {
            "post": {
                "description": "Reenvía el código de verificación al usuario si el código anterior ha expirado.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Reenviar código de verificación",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Respuesta exitosa al reenviar código de verificación",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "El usuario ya está verificado o el código de verificación aún es válido",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Usuario no autorizado",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/update": {
            "patch": {
                "description": "Actualiza el displayName en Firebase Auth y los campos rut y birthdate en Firestore",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Actualiza el perfil de usuario",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token de autorización JWT",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Datos de actualización del perfil",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.UpdateProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Perfil actualizado exitosamente",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "Datos de solicitud inválidos",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Usuario no autorizado",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/upload-photo": {
            "post": {
                "description": "Cargar una foto de perfil para el usuario actual.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Cargar foto de perfil",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token de autorización JWT",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Archivo de imagen",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Foto de perfil cargada correctamente",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "No se ha enviado ningún archivo",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "No autorizado",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error al cargar la foto de perfil",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/verify-code": {
            "post": {
                "description": "Verifica el código de verificación de un usuario en Firestore y Firebase Auth.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verifica el código de verificación de un usuario.",
                "parameters": [
                    {
                        "description": "Datos de la solicitud",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.VerifyCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Usuario verificado",
                        "schema": {
                            "$ref": "#/definitions/httputil.StandardResponse"
                        }
                    },
                    "400": {
                        "description": "Datos de solicitud inválidos",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Código de verificación incorrecto",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error interno del servidor",
                        "schema": {
                            "$ref": "#/definitions/httputil.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.ChangePasswordRequest": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "password": {
                    "type": "string"
                }
            }
        },
        "controllers.ForgotPasswordRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "controllers.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "controllers.RegisterRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "controllers.UpdateProfileRequest": {
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string"
                },
                "displayName": {
                    "type": "string"
                },
                "rut": {
                    "type": "string"
                }
            }
        },
        "controllers.VerifyCodeRequest": {
            "type": "object",
            "properties": {
                "uid": {
                    "type": "string"
                },
                "verificationCode": {
                    "type": "string"
                }
            }
        },
        "httputil.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.StandardResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "backend-autch.tssw.cl",
	BasePath:         "/",
	Schemes:          []string{"https"},
	Title:            "Backend API",
	Description:      "API para el backend de una aplicación web con autenticación de Firebase.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
