# Swagger Documentation Setup

## ✅ Swagger Works Out of the Box!

The Swagger documentation is **already generated and included** in the repository. When you clone the code, Swagger will work immediately without running `swag init`.

##  Swagger UI Access

Once services are running, access Swagger UI at:

1. **Auth Service**: http://localhost:8002/swagger/index.html
2. **Contact Service**: http://localhost:8001/swagger/index.html
3. **Inventory Service**: http://localhost:8003/swagger/index.html
4. **Sales Service**: http://localhost:8004/swagger/index.html
5. **Purchase Service**: http://localhost:8005/swagger/index.html

## 🔄 Regenerating Swagger Docs (Optional)

If you modify Swagger annotations in the code and want to regenerate the docs:

### Option 1: Use the Script (Recommended)

```bash
./generate_swagger.sh
```

This script will regenerate Swagger docs for all services.

### Option 2: Manual Regeneration

For each service:

```bash
# Auth Service
cd services/auth
swag init -g cmd/main.go --parseDependency --parseInternal

# Contact Service
cd services/contact
swag init -g cmd/main.go --parseDependency --parseInternal

# Inventory Service
cd services/inventory
swag init -g cmd/main.go --parseDependency --parseInternal

# Sales Service
cd services/sales
swag init -g cmd/main.go --parseDependency --parseInternal

# Purchase Service
cd services/purchase
swag init -g cmd/main.go --parseDependency --parseInternal
```

## 📝 Swagger Annotations

Swagger annotations are added to:

1. **Main files** (`cmd/main.go`): API metadata
   ```go
   // @title           Service API
   // @version         1.0
   // @description     Service description
   // @host            localhost:8000
   // @BasePath        /
   ```

2. **Handler functions** (`httphandler/httphandler.go`): Endpoint documentation
   ```go
   // @Summary      Endpoint summary
   // @Description  Detailed description
   // @Tags         tag-name
   // @Accept       json
   // @Produce      json
   // @Param        request body Model true "Request description"
   // @Success      200 {object} response.SuccessResponse{data=Model}
   // @Failure      400 {object} response.ValidationErrorResponse
   // @Router       /endpoint [post]
   ```

3. **Model structs** (`model/*.go`): Field examples
   ```go
   type Model struct {
       Field string `json:"field" example:"example-value"`
   }
   ```

## 🎯 Features

1. **Complete API Documentation**: All endpoints documented
2. **Request/Response Examples**: Examples in DTOs
3. **Try It Out**: Test endpoints directly from Swagger UI
4. **Authentication Support**: Bearer token authentication
5. **Schema Definitions**: Complete data models

## 🔐 Authentication in Swagger

To use authenticated endpoints in Swagger UI:

1. Click the **"Authorize"** button at the top
2. Enter your JWT token: `Bearer <your-token>`
3. Click **"Authorize"**
4. All protected endpoints will now use this token

## 📦 Files Included

Each service has a `docs/` folder with:

1. `docs.go` - Generated Go code
2. `swagger.json` - JSON format documentation
3. `swagger.yaml` - YAML format documentation

These files are **committed to the repository**, so Swagger works immediately after cloning.

## 🚀 Quick Start

1. **Clone the repository**
2. **Start services**: `docker compose up -d`
3. **Access Swagger UI**: http://localhost:8002/swagger/index.html
4. **No additional setup required!**

---

**Note**: If you modify Swagger annotations, regenerate docs using `./generate_swagger.sh` and commit the updated `docs/` folders.

