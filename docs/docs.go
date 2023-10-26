// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/auth/signup-or-login": {
            "post": {
                "description": "signup or login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Signup Or Login",
                "parameters": [
                    {
                        "description": "signup or login",
                        "name": "signupOrLogin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth_handler.SignupOrLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth_handler.SignupOrLoginResponse"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "description": "get user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authoriztion",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user_handler.GetUserResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth_handler.GoogleToken": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "idToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "auth_handler.ProviderToken": {
            "type": "object",
            "properties": {
                "googleToken": {
                    "$ref": "#/definitions/auth_handler.GoogleToken"
                }
            }
        },
        "auth_handler.SignUpDetails": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        },
        "auth_handler.SignupOrLoginRequest": {
            "type": "object",
            "properties": {
                "deviceId": {
                    "type": "string"
                },
                "deviceOs": {
                    "type": "string"
                },
                "providerToken": {
                    "$ref": "#/definitions/auth_handler.ProviderToken"
                },
                "signUpDetails": {
                    "description": "TODO: will be removed once SSO integrated",
                    "allOf": [
                        {
                            "$ref": "#/definitions/auth_handler.SignUpDetails"
                        }
                    ]
                }
            }
        },
        "auth_handler.SignupOrLoginResponse": {
            "type": "object",
            "properties": {
                "authToken": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        },
        "user_handler.GetUserResponse": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
