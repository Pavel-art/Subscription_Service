package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterSwagger регистрирует эндпоинты для Swagger UI и OpenAPI спецификации.
func RegisterSwagger(r *gin.Engine) {
	r.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(openAPISpecJSON))
	})

	r.GET("/swagger", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})
}

// Минимальная OpenAPI 3.0 спецификация для текущих эндпоинтов.
const openAPISpecJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Subscription Service API",
    "version": "1.0.0",
    "description": "API для управления подписками и расчета их стоимости."
  },
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {
                      "type": "string",
                      "example": "ok"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/subscriptions": {
      "post": {
        "summary": "Create subscription",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/CreateSubscriptionRequest"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Created",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Subscription"
                }
              }
            }
          },
          "400": {
            "description": "Invalid request"
          },
          "500": {
            "description": "Internal error"
          }
        }
      },
      "get": {
        "summary": "Get all subscriptions",
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "schema": {
              "type": "integer",
              "format": "int64",
              "default": 1
            }
          },
          {
            "name": "page_size",
            "in": "query",
            "schema": {
              "type": "integer",
              "format": "int64",
              "default": 20,
              "maximum": 100
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/GetAllResponse"
                }
              }
            }
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/api/v1/subscriptions/{id}": {
      "get": {
        "summary": "Get subscription by ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Subscription"
                }
              }
            }
          },
          "400": {
            "description": "Bad request"
          },
          "404": {
            "description": "Not found"
          }
        }
      },
      "put": {
        "summary": "Update subscription",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UpdateSubscriptionRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Subscription"
                }
              }
            }
          },
          "400": {
            "description": "Bad request"
          },
          "404": {
            "description": "Not found"
          }
        }
      },
      "delete": {
        "summary": "Delete subscription",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "No Content"
          },
          "400": {
            "description": "Bad request"
          },
          "404": {
            "description": "Not found"
          }
        }
      }
    },
    "/api/v1/subscriptions/cost": {
      "get": {
        "summary": "Calculate total cost of subscriptions",
        "parameters": [
          {
            "name": "user_id",
            "in": "query",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "name": "service_name",
            "in": "query",
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "from",
            "in": "query",
            "schema": {
              "type": "string",
              "format": "date-time"
            }
          },
          {
            "name": "to",
            "in": "query",
            "schema": {
              "type": "string",
              "format": "date-time"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "total_cost": {
                      "type": "integer",
                      "format": "int64"
                    },
                    "currency": {
                      "type": "string",
                      "example": "RUB"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "Bad request"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Subscription": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid"
          },
          "service_name": {
            "type": "string"
          },
          "price": {
            "type": "integer",
            "format": "int64"
          },
          "user_id": {
            "type": "string",
            "format": "uuid"
          },
          "start_date": {
            "type": "string",
            "format": "date-time"
          },
          "end_date": {
            "type": "string",
            "format": "date-time",
            "nullable": true
          },
          "created_at": {
            "type": "string",
            "format": "date-time"
          },
          "updated_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "id",
          "service_name",
          "price",
          "user_id",
          "start_date",
          "created_at",
          "updated_at"
        ]
      },
      "CreateSubscriptionRequest": {
        "type": "object",
        "properties": {
          "service_name": {
            "type": "string"
          },
          "price": {
            "type": "integer",
            "format": "int64"
          },
          "user_id": {
            "type": "string",
            "format": "uuid"
          },
          "start_date": {
            "type": "string",
            "format": "date-time"
          },
          "end_date": {
            "type": "string",
            "format": "date-time",
            "nullable": true
          }
        },
        "required": [
          "service_name",
          "price",
          "user_id",
          "start_date"
        ]
      },
      "UpdateSubscriptionRequest": {
        "type": "object",
        "properties": {
          "service_name": {
            "type": "string"
          },
          "price": {
            "type": "integer",
            "format": "int64"
          },
          "start_date": {
            "type": "string",
            "format": "date-time"
          },
          "end_date": {
            "type": "string",
            "format": "date-time",
            "nullable": true
          }
        }
      },
      "GetAllResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/Subscription"
            }
          },
          "pagination": {
            "$ref": "#/components/schemas/PaginationInfo"
          }
        }
      },
      "PaginationInfo": {
        "type": "object",
        "properties": {
          "page": {
            "type": "integer",
            "format": "int64"
          },
          "page_size": {
            "type": "integer",
            "format": "int64"
          },
          "total_count": {
            "type": "integer",
            "format": "int64"
          },
          "total_pages": {
            "type": "integer",
            "format": "int64"
          }
        }
      }
    }
  }
}
`

// Простейшая HTML-страница со Swagger UI, использующая CDN.
const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Subscription Service API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; padding: 0; }
    #swagger-ui { width: 100%; height: 100vh; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: "#swagger-ui"
      });
    };
  </script>
</body>
</html>
`
