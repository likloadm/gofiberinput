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
        "/login": {
            "post": {
                "description": "login user by name and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "login user by name and password",
                "operationId": "Login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User name",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User password",
                        "name": "pass",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "access JWT token",
                        "schema": {
                            "$ref": "#/definitions/model.AccessTokenJWT"
                        }
                    },
                    "400": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "register user by name and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "register user by name and password",
                "operationId": "Register",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User password",
                        "name": "pass",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User registered",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    },
                    "400": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    }
                }
            }
        },
        "/transaction": {
            "get": {
                "description": "get transaction by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get transaction by id",
                "operationId": "GetTransaction",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Transaction id",
                        "name": "id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/model.TransactionUserSystem"
                        }
                    },
                    "400": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    },
                    "401": {
                        "description": "Bearer auth required",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    }
                }
            },
            "post": {
                "description": "get id by type transaction and amount \u003e 0",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new transaction to the database and change balance",
                "operationId": "transaction",
                "parameters": [
                    {
                        "type": "number",
                        "description": "Transaction amount(\u003e0)",
                        "name": "amount",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "1 input, 2 output",
                        "name": "type",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/model.SystemUserBalance"
                        }
                    },
                    "400": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    },
                    "401": {
                        "description": "Bearer auth required",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "description": "get user by jwt auth",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get user by jwt auth",
                "operationId": "GetUser",
                "responses": {
                    "200": {
                        "description": "user",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    },
                    "401": {
                        "description": "Bearer auth required",
                        "schema": {
                            "$ref": "#/definitions/model.MessageModel"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AccessTokenJWT": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "model.MessageModel": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "error"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "model.SystemUserBalance": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.TransactionUserSystem": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "sender": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                },
                "type": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Bearer",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Go Fiber Input",
	Description:      "API for controlling the receipt and withdrawal of funds",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
