// Package docs provides Swagger documentation for the API
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Bahrul Ulum Fadhlur Rohman",
            "email": "bahrululumfr@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health": {
            "get": {
                "description": "Get API health status",
                "produces": ["application/json"],
                "tags": ["Health"],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "status": {"type": "string"},
                                "message": {"type": "string"}
                            }
                        }
                    }
                }
            }
        },
        "/v1/auth/register": {
            "post": {
                "description": "Register a new user account",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Register new user",
                "parameters": [{
                    "description": "User registration data",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/RegisterInput"}
                }],
                "responses": {
                    "201": {"description": "Registration successful", "schema": {"$ref": "#/definitions/TokenResponse"}},
                    "400": {"description": "Invalid input"},
                    "409": {"description": "User already exists"}
                }
            }
        },
        "/v1/auth/login": {
            "post": {
                "description": "Login with email and password",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "User login",
                "parameters": [{
                    "description": "Login credentials",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/LoginInput"}
                }],
                "responses": {
                    "200": {"description": "Login successful", "schema": {"$ref": "#/definitions/TokenResponse"}},
                    "401": {"description": "Invalid credentials"}
                }
            }
        },
        "/v1/auth/refresh": {
            "post": {
                "description": "Refresh access token using refresh token",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Refresh token",
                "parameters": [{
                    "description": "Refresh token",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/RefreshTokenInput"}
                }],
                "responses": {
                    "200": {"description": "Token refreshed", "schema": {"$ref": "#/definitions/TokenResponse"}},
                    "401": {"description": "Invalid refresh token"}
                }
            }
        },
        "/v1/auth/me": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get current authenticated user info",
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Get current user",
                "responses": {
                    "200": {"description": "User info", "schema": {"$ref": "#/definitions/User"}},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/v1/auth/logout": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "Logout from current session by deleting the refresh token from database",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Logout",
                "parameters": [{
                    "description": "Logout request with refresh token",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/LogoutInput"}
                }],
                "responses": {
                    "200": {"description": "Logout successful"},
                    "400": {"description": "Invalid refresh token"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/v1/auth/logout-all": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "Logout from all devices (invalidates all refresh tokens)",
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Logout from all devices",
                "responses": {
                    "200": {"description": "Logged out from all devices"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/v1/auth/profile": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update current user profile (name, email, password, photo)",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Auth"],
                "summary": "Update user profile",
                "parameters": [{
                    "description": "Profile update data",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/UpdateProfileInput"}
                }],
                "responses": {
                    "200": {"description": "Profile updated successfully", "schema": {"$ref": "#/definitions/User"}},
                    "400": {"description": "Invalid input or current password incorrect"},
                    "401": {"description": "Unauthorized"},
                    "409": {"description": "Email already in use"}
                }
            }
        },
        "/v1/public/projects": {
            "get": {
                "description": "Get list of published projects",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "List projects",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 10},
                    {"name": "search", "in": "query", "type": "string"},
                    {"name": "category", "in": "query", "type": "string"},
                    {"name": "tag", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {"description": "List of projects", "schema": {"$ref": "#/definitions/ProjectListResponse"}}
                }
            }
        },
        "/v1/public/projects/{slug}": {
            "get": {
                "description": "Get project details by slug",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "Get project by slug",
                "parameters": [{"name": "slug", "in": "path", "required": true, "type": "string"}],
                "responses": {
                    "200": {"description": "Project details", "schema": {"$ref": "#/definitions/Project"}},
                    "404": {"description": "Project not found"}
                }
            }
        },
        "/v1/public/categories": {
            "get": {
                "description": "Get all categories",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "List categories",
                "responses": {
                    "200": {"description": "List of categories", "schema": {"type": "array", "items": {"$ref": "#/definitions/Category"}}}
                }
            }
        },
        "/v1/public/tags": {
            "get": {
                "description": "Get all tags",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "List tags",
                "responses": {
                    "200": {"description": "List of tags", "schema": {"type": "array", "items": {"$ref": "#/definitions/Tag"}}}
                }
            }
        },
        "/v1/public/careers": {
            "get": {
                "description": "Get career/work experience list",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "List careers",
                "responses": {
                    "200": {"description": "List of careers", "schema": {"type": "array", "items": {"$ref": "#/definitions/Career"}}}
                }
            }
        },
        "/v1/public/educations": {
            "get": {
                "description": "Get education history",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "List educations",
                "responses": {
                    "200": {"description": "List of educations", "schema": {"type": "array", "items": {"$ref": "#/definitions/Education"}}}
                }
            }
        },
        "/v1/public/resume": {
            "get": {
                "description": "Get active resume/CV",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "Get active resume",
                "responses": {
                    "200": {"description": "Active resume", "schema": {"$ref": "#/definitions/Resume"}},
                    "404": {"description": "No active resume"}
                }
            }
        },
        "/v1/public/contact": {
            "post": {
                "description": "Submit contact form (rate limited)",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "Submit contact form",
                "parameters": [{
                    "description": "Contact form data",
                    "name": "body",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/CreateContactInput"}
                }],
                "responses": {
                    "201": {"description": "Message sent"},
                    "429": {"description": "Too many requests"}
                }
            }
        },
        "/v1/admin/projects": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all projects (including unpublished)",
                "produces": ["application/json"],
                "tags": ["Admin - Projects"],
                "summary": "List all projects",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 10}
                ],
                "responses": {
                    "200": {"description": "List of all projects"},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Forbidden - Admin only"}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new project",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Projects"],
                "summary": "Create project",
                "parameters": [{"description": "Project data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateProjectInput"}}],
                "responses": {
                    "201": {"description": "Project created"},
                    "401": {"description": "Unauthorized"}
                }
            }
        },
        "/v1/admin/projects/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get project by ID",
                "produces": ["application/json"],
                "tags": ["Admin - Projects"],
                "summary": "Get project by ID",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {
                    "200": {"description": "Project details", "schema": {"$ref": "#/definitions/Project"}},
                    "404": {"description": "Not found"}
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update a project",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Projects"],
                "summary": "Update project",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Project data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateProjectInput"}}
                ],
                "responses": {"200": {"description": "Project updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a project",
                "produces": ["application/json"],
                "tags": ["Admin - Projects"],
                "summary": "Delete project",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Project deleted"}}
            }
        },
        "/v1/admin/categories": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all categories",
                "produces": ["application/json"],
                "tags": ["Admin - Categories"],
                "summary": "List all categories",
                "responses": {
                    "200": {"description": "List of categories", "schema": {"type": "array", "items": {"$ref": "#/definitions/Category"}}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new category",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Categories"],
                "summary": "Create category",
                "parameters": [{"description": "Category data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateCategoryInput"}}],
                "responses": {"201": {"description": "Category created"}}
            }
        },
        "/v1/admin/categories/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update a category",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Categories"],
                "summary": "Update category",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Category data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateCategoryInput"}}
                ],
                "responses": {"200": {"description": "Category updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a category",
                "produces": ["application/json"],
                "tags": ["Admin - Categories"],
                "summary": "Delete category",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Category deleted"}}
            }
        },
        "/v1/admin/tags": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all tags",
                "produces": ["application/json"],
                "tags": ["Admin - Tags"],
                "summary": "List all tags",
                "responses": {
                    "200": {"description": "List of tags", "schema": {"type": "array", "items": {"$ref": "#/definitions/Tag"}}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new tag",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Tags"],
                "summary": "Create tag",
                "parameters": [{"description": "Tag data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateTagInput"}}],
                "responses": {"201": {"description": "Tag created"}}
            }
        },
        "/v1/admin/tags/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update a tag",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Tags"],
                "summary": "Update tag",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Tag data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateTagInput"}}
                ],
                "responses": {"200": {"description": "Tag updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a tag",
                "produces": ["application/json"],
                "tags": ["Admin - Tags"],
                "summary": "Delete tag",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Tag deleted"}}
            }
        },
        "/v1/admin/careers": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all careers",
                "produces": ["application/json"],
                "tags": ["Admin - Careers"],
                "summary": "List all careers",
                "responses": {
                    "200": {"description": "List of careers", "schema": {"type": "array", "items": {"$ref": "#/definitions/Career"}}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new career entry",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Careers"],
                "summary": "Create career",
                "parameters": [{"description": "Career data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateCareerInput"}}],
                "responses": {"201": {"description": "Career created"}}
            }
        },
        "/v1/admin/careers/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update a career entry",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Careers"],
                "summary": "Update career",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Career data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateCareerInput"}}
                ],
                "responses": {"200": {"description": "Career updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a career entry",
                "produces": ["application/json"],
                "tags": ["Admin - Careers"],
                "summary": "Delete career",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Career deleted"}}
            }
        },
        "/v1/admin/educations": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all educations",
                "produces": ["application/json"],
                "tags": ["Admin - Educations"],
                "summary": "List all educations",
                "responses": {
                    "200": {"description": "List of educations", "schema": {"type": "array", "items": {"$ref": "#/definitions/Education"}}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new education entry",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Educations"],
                "summary": "Create education",
                "parameters": [{"description": "Education data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateEducationInput"}}],
                "responses": {"201": {"description": "Education created"}}
            }
        },
        "/v1/admin/educations/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update an education entry",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Educations"],
                "summary": "Update education",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Education data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateEducationInput"}}
                ],
                "responses": {"200": {"description": "Education updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete an education entry",
                "produces": ["application/json"],
                "tags": ["Admin - Educations"],
                "summary": "Delete education",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Education deleted"}}
            }
        },
        "/v1/admin/resumes": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all resumes",
                "produces": ["application/json"],
                "tags": ["Admin - Resumes"],
                "summary": "List all resumes",
                "responses": {
                    "200": {"description": "List of resumes", "schema": {"type": "array", "items": {"$ref": "#/definitions/Resume"}}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new resume",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Resumes"],
                "summary": "Create resume",
                "parameters": [{"description": "Resume data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateResumeInput"}}],
                "responses": {"201": {"description": "Resume created"}}
            }
        },
        "/v1/admin/resumes/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update a resume",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Resumes"],
                "summary": "Update resume",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Resume data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateResumeInput"}}
                ],
                "responses": {"200": {"description": "Resume updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a resume",
                "produces": ["application/json"],
                "tags": ["Admin - Resumes"],
                "summary": "Delete resume",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Resume deleted"}}
            }
        },
        "/v1/admin/resumes/{id}/activate": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "Set a resume as active (deactivates others)",
                "produces": ["application/json"],
                "tags": ["Admin - Resumes"],
                "summary": "Activate resume",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Resume activated"}}
            }
        },
        "/v1/admin/contacts": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all contact submissions",
                "produces": ["application/json"],
                "tags": ["Admin - Contacts"],
                "summary": "List contacts",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 10},
                    {"name": "is_read", "in": "query", "type": "boolean"}
                ],
                "responses": {
                    "200": {"description": "List of contacts", "schema": {"$ref": "#/definitions/ContactListResponse"}}
                }
            }
        },
        "/v1/admin/contacts/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get contact by ID",
                "produces": ["application/json"],
                "tags": ["Admin - Contacts"],
                "summary": "Get contact",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {
                    "200": {"description": "Contact details", "schema": {"$ref": "#/definitions/Contact"}},
                    "404": {"description": "Not found"}
                }
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a contact",
                "produces": ["application/json"],
                "tags": ["Admin - Contacts"],
                "summary": "Delete contact",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Contact deleted"}}
            }
        },
        "/v1/admin/contacts/{id}/read": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Mark contact as read",
                "produces": ["application/json"],
                "tags": ["Admin - Contacts"],
                "summary": "Mark as read",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Contact marked as read"}}
            }
        },
        "/v1/admin/upload-url": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "Get presigned URL for file upload to R2",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Upload"],
                "summary": "Get upload URL",
                "parameters": [{"description": "Upload request", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UploadURLRequest"}}],
                "responses": {
                    "200": {"description": "Presigned URL", "schema": {"$ref": "#/definitions/PresignedURLResponse"}}
                }
            }
        },
        "/v1/admin/users": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all registered users (admin only)",
                "produces": ["application/json"],
                "tags": ["Admin - Users"],
                "summary": "List all users",
                "responses": {
                    "200": {"description": "List of users", "schema": {"type": "array", "items": {"$ref": "#/definitions/User"}}},
                    "401": {"description": "Unauthorized"},
                    "403": {"description": "Forbidden - Admin only"}
                }
            }
        },
        "/v1/public/about": {
            "get": {
                "description": "Get active about/profile information",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "Get about info",
                "responses": {
                    "200": {"description": "About info", "schema": {"$ref": "#/definitions/About"}},
                    "404": {"description": "No active profile found"}
                }
            }
        },
        "/v1/public/blogs": {
            "get": {
                "description": "Get list of published blogs",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "List blogs",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 10},
                    {"name": "tag_id", "in": "query", "type": "string"},
                    {"name": "search", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {"description": "List of blogs", "schema": {"$ref": "#/definitions/BlogListResponse"}}
                }
            }
        },
        "/v1/public/blogs/{slug}": {
            "get": {
                "description": "Get blog by slug",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "Get blog by slug",
                "parameters": [{"name": "slug", "in": "path", "required": true, "type": "string"}],
                "responses": {
                    "200": {"description": "Blog details", "schema": {"$ref": "#/definitions/Blog"}},
                    "404": {"description": "Blog not found"}
                }
            }
        },
        "/v1/public/certificates": {
            "get": {
                "description": "Get all certificates",
                "produces": ["application/json"],
                "tags": ["Public"],
                "summary": "List certificates",
                "responses": {
                    "200": {"description": "List of certificates", "schema": {"type": "array", "items": {"$ref": "#/definitions/Certificate"}}}
                }
            }
        },
        "/v1/admin/about": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all about entries",
                "produces": ["application/json"],
                "tags": ["Admin - About"],
                "summary": "List all about entries",
                "responses": {
                    "200": {"description": "List of about entries", "schema": {"type": "array", "items": {"$ref": "#/definitions/About"}}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new about entry",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - About"],
                "summary": "Create about",
                "parameters": [{"description": "About data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateAboutInput"}}],
                "responses": {"201": {"description": "About created"}}
            }
        },
        "/v1/admin/about/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update an about entry",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - About"],
                "summary": "Update about",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "About data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateAboutInput"}}
                ],
                "responses": {"200": {"description": "About updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete an about entry",
                "produces": ["application/json"],
                "tags": ["Admin - About"],
                "summary": "Delete about",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "About deleted"}}
            }
        },
        "/v1/admin/blogs": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all blogs (including unpublished)",
                "produces": ["application/json"],
                "tags": ["Admin - Blogs"],
                "summary": "List all blogs",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "limit", "in": "query", "type": "integer", "default": 10}
                ],
                "responses": {
                    "200": {"description": "List of blogs", "schema": {"$ref": "#/definitions/BlogListResponse"}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new blog",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Blogs"],
                "summary": "Create blog",
                "parameters": [{"description": "Blog data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateBlogInput"}}],
                "responses": {"201": {"description": "Blog created"}}
            }
        },
        "/v1/admin/blogs/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get blog by ID",
                "produces": ["application/json"],
                "tags": ["Admin - Blogs"],
                "summary": "Get blog by ID",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {
                    "200": {"description": "Blog details", "schema": {"$ref": "#/definitions/Blog"}},
                    "404": {"description": "Not found"}
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update a blog",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Blogs"],
                "summary": "Update blog",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Blog data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateBlogInput"}}
                ],
                "responses": {"200": {"description": "Blog updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a blog",
                "produces": ["application/json"],
                "tags": ["Admin - Blogs"],
                "summary": "Delete blog",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Blog deleted"}}
            }
        },
        "/v1/admin/certificates": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "Get all certificates",
                "produces": ["application/json"],
                "tags": ["Admin - Certificates"],
                "summary": "List all certificates",
                "responses": {
                    "200": {"description": "List of certificates", "schema": {"type": "array", "items": {"$ref": "#/definitions/Certificate"}}}
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "Create a new certificate",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Certificates"],
                "summary": "Create certificate",
                "parameters": [{"description": "Certificate data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/CreateCertificateInput"}}],
                "responses": {"201": {"description": "Certificate created"}}
            }
        },
        "/v1/admin/certificates/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "Update a certificate",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["Admin - Certificates"],
                "summary": "Update certificate",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string"},
                    {"description": "Certificate data", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/UpdateCertificateInput"}}
                ],
                "responses": {"200": {"description": "Certificate updated"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "Delete a certificate",
                "produces": ["application/json"],
                "tags": ["Admin - Certificates"],
                "summary": "Delete certificate",
                "parameters": [{"name": "id", "in": "path", "required": true, "type": "string"}],
                "responses": {"200": {"description": "Certificate deleted"}}
            }
        }
    },
    "definitions": {
        "LoginInput": {
            "type": "object",
            "required": ["email", "password"],
            "properties": {
                "email": {"type": "string", "example": "john@example.com"},
                "password": {"type": "string", "example": "password123"}
            }
        },
        "RegisterInput": {
            "type": "object",
            "required": ["name", "email", "password"],
            "properties": {
                "name": {"type": "string", "example": "John Doe"},
                "email": {"type": "string", "example": "john@example.com"},
                "password": {"type": "string", "example": "password123"}
            }
        },
        "RefreshTokenInput": {
            "type": "object",
            "required": ["refresh_token"],
            "properties": {
                "refresh_token": {"type": "string"}
            }
        },
        "LogoutInput": {
            "type": "object",
            "required": ["refresh_token"],
            "properties": {
                "refresh_token": {"type": "string", "description": "The refresh token to invalidate"}
            }
        },
        "UpdateProfileInput": {
            "type": "object",
            "properties": {
                "name": {"type": "string", "description": "New name (2-100 characters)"},
                "email": {"type": "string", "format": "email", "description": "New email address"},
                "current_password": {"type": "string", "description": "Current password (required when changing password)"},
                "new_password": {"type": "string", "description": "New password (6-100 characters)"},
                "image": {"type": "string", "description": "Profile image URL"}
            }
        },
        "TokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {"type": "string"},
                "refresh_token": {"type": "string"},
                "token_type": {"type": "string", "example": "Bearer"},
                "expires_in": {"type": "integer", "example": 3600}
            }
        },
        "User": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "email": {"type": "string"},
                "role": {"type": "string", "enum": ["USER", "ADMIN"]},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "Project": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "title": {"type": "string"},
                "slug": {"type": "string"},
                "description": {"type": "string"},
                "content": {"type": "string"},
                "thumbnail_url": {"type": "string"},
                "demo_url": {"type": "string"},
                "repo_url": {"type": "string"},
                "is_published": {"type": "boolean"},
                "is_featured": {"type": "boolean"},
                "sort_order": {"type": "integer"},
                "categories": {"type": "array", "items": {"$ref": "#/definitions/Category"}},
                "tags": {"type": "array", "items": {"$ref": "#/definitions/Tag"}},
                "images": {"type": "array", "items": {"$ref": "#/definitions/ProjectImage"}},
                "created_at": {"type": "string", "format": "date-time"},
                "updated_at": {"type": "string", "format": "date-time"}
            }
        },
        "ProjectImage": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "url": {"type": "string"},
                "alt": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "ProjectListResponse": {
            "type": "object",
            "properties": {
                "success": {"type": "boolean"},
                "data": {"type": "array", "items": {"$ref": "#/definitions/Project"}},
                "pagination": {"$ref": "#/definitions/Pagination"}
            }
        },
        "CreateProjectInput": {
            "type": "object",
            "required": ["title", "slug"],
            "properties": {
                "title": {"type": "string"},
                "slug": {"type": "string"},
                "description": {"type": "string"},
                "content": {"type": "string"},
                "thumbnail_url": {"type": "string"},
                "demo_url": {"type": "string"},
                "repo_url": {"type": "string"},
                "is_published": {"type": "boolean"},
                "is_featured": {"type": "boolean"},
                "sort_order": {"type": "integer"},
                "category_ids": {"type": "array", "items": {"type": "string"}},
                "tag_ids": {"type": "array", "items": {"type": "string"}}
            }
        },
        "UpdateProjectInput": {
            "type": "object",
            "properties": {
                "title": {"type": "string"},
                "slug": {"type": "string"},
                "description": {"type": "string"},
                "content": {"type": "string"},
                "thumbnail_url": {"type": "string"},
                "demo_url": {"type": "string"},
                "repo_url": {"type": "string"},
                "is_published": {"type": "boolean"},
                "is_featured": {"type": "boolean"},
                "sort_order": {"type": "integer"},
                "category_ids": {"type": "array", "items": {"type": "string"}},
                "tag_ids": {"type": "array", "items": {"type": "string"}}
            }
        },
        "Category": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "slug": {"type": "string"},
                "description": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateCategoryInput": {
            "type": "object",
            "required": ["name", "slug"],
            "properties": {
                "name": {"type": "string"},
                "slug": {"type": "string"},
                "description": {"type": "string"}
            }
        },
        "UpdateCategoryInput": {
            "type": "object",
            "properties": {
                "name": {"type": "string"},
                "slug": {"type": "string"},
                "description": {"type": "string"}
            }
        },
        "Tag": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "slug": {"type": "string"},
                "icon_url": {"type": "string", "description": "URL to tech stack icon (e.g., devicons)"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateTagInput": {
            "type": "object",
            "required": ["name", "slug"],
            "properties": {
                "name": {"type": "string"},
                "slug": {"type": "string"},
                "icon_url": {"type": "string"}
            }
        },
        "UpdateTagInput": {
            "type": "object",
            "properties": {
                "name": {"type": "string"},
                "slug": {"type": "string"},
                "icon_url": {"type": "string"}
            }
        },
        "Career": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "company": {"type": "string"},
                "position": {"type": "string"},
                "location": {"type": "string"},
                "description": {"type": "string"},
                "start_date": {"type": "string", "format": "date-time"},
                "end_date": {"type": "string", "format": "date-time"},
                "is_current": {"type": "boolean"},
                "logo_url": {"type": "string"},
                "company_url": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "CreateCareerInput": {
            "type": "object",
            "required": ["company", "position", "start_date"],
            "properties": {
                "company": {"type": "string"},
                "position": {"type": "string"},
                "location": {"type": "string"},
                "description": {"type": "string"},
                "start_date": {"type": "string", "format": "date-time"},
                "end_date": {"type": "string", "format": "date-time"},
                "is_current": {"type": "boolean"},
                "logo_url": {"type": "string"},
                "company_url": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "UpdateCareerInput": {
            "type": "object",
            "properties": {
                "company": {"type": "string"},
                "position": {"type": "string"},
                "location": {"type": "string"},
                "description": {"type": "string"},
                "start_date": {"type": "string", "format": "date-time"},
                "end_date": {"type": "string", "format": "date-time"},
                "is_current": {"type": "boolean"},
                "logo_url": {"type": "string"},
                "company_url": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "Education": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "school": {"type": "string"},
                "degree": {"type": "string"},
                "field": {"type": "string"},
                "location": {"type": "string"},
                "description": {"type": "string"},
                "start_date": {"type": "string", "format": "date-time"},
                "end_date": {"type": "string", "format": "date-time"},
                "gpa": {"type": "string"},
                "logo_url": {"type": "string"},
                "school_url": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "CreateEducationInput": {
            "type": "object",
            "required": ["school", "degree", "start_date"],
            "properties": {
                "school": {"type": "string"},
                "degree": {"type": "string"},
                "field": {"type": "string"},
                "location": {"type": "string"},
                "description": {"type": "string"},
                "start_date": {"type": "string", "format": "date-time"},
                "end_date": {"type": "string", "format": "date-time"},
                "gpa": {"type": "string"},
                "logo_url": {"type": "string"},
                "school_url": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "UpdateEducationInput": {
            "type": "object",
            "properties": {
                "school": {"type": "string"},
                "degree": {"type": "string"},
                "field": {"type": "string"},
                "location": {"type": "string"},
                "description": {"type": "string"},
                "start_date": {"type": "string", "format": "date-time"},
                "end_date": {"type": "string", "format": "date-time"},
                "gpa": {"type": "string"},
                "logo_url": {"type": "string"},
                "school_url": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "Resume": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "file_url": {"type": "string"},
                "file_name": {"type": "string"},
                "file_size": {"type": "integer"},
                "version": {"type": "string"},
                "is_active": {"type": "boolean"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateResumeInput": {
            "type": "object",
            "required": ["file_url", "file_name"],
            "properties": {
                "file_url": {"type": "string"},
                "file_name": {"type": "string"},
                "file_size": {"type": "integer"},
                "version": {"type": "string"},
                "is_active": {"type": "boolean"}
            }
        },
        "UpdateResumeInput": {
            "type": "object",
            "properties": {
                "file_url": {"type": "string"},
                "file_name": {"type": "string"},
                "file_size": {"type": "integer"},
                "version": {"type": "string"},
                "is_active": {"type": "boolean"}
            }
        },
        "Contact": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "email": {"type": "string"},
                "subject": {"type": "string"},
                "message": {"type": "string"},
                "is_read": {"type": "boolean"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateContactInput": {
            "type": "object",
            "required": ["name", "email", "message"],
            "properties": {
                "name": {"type": "string"},
                "email": {"type": "string"},
                "subject": {"type": "string"},
                "message": {"type": "string"}
            }
        },
        "ContactListResponse": {
            "type": "object",
            "properties": {
                "success": {"type": "boolean"},
                "data": {"type": "array", "items": {"$ref": "#/definitions/Contact"}},
                "pagination": {"$ref": "#/definitions/Pagination"}
            }
        },
        "UploadURLRequest": {
            "type": "object",
            "required": ["file_name", "content_type"],
            "properties": {
                "file_name": {"type": "string", "example": "image.jpg"},
                "content_type": {"type": "string", "example": "image/jpeg"},
                "folder": {"type": "string", "example": "projects"}
            }
        },
        "PresignedURLResponse": {
            "type": "object",
            "properties": {
                "upload_url": {"type": "string"},
                "file_url": {"type": "string"},
                "key": {"type": "string"},
                "expires_in": {"type": "integer"}
            }
        },
        "Pagination": {
            "type": "object",
            "properties": {
                "page": {"type": "integer"},
                "limit": {"type": "integer"},
                "total": {"type": "integer"},
                "total_pages": {"type": "integer"},
                "has_next": {"type": "boolean"},
                "has_prev": {"type": "boolean"}
            }
        },
        "About": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "full_name": {"type": "string"},
                "nickname": {"type": "string"},
                "role": {"type": "string"},
                "bio": {"type": "string"},
                "avatar_url": {"type": "string"},
                "cover_url": {"type": "string"},
                "location": {"type": "string"},
                "email": {"type": "string"},
                "phone": {"type": "string"},
                "is_active": {"type": "boolean"},
                "created_at": {"type": "string", "format": "date-time"},
                "updated_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateAboutInput": {
            "type": "object",
            "required": ["full_name", "role"],
            "properties": {
                "full_name": {"type": "string"},
                "nickname": {"type": "string"},
                "role": {"type": "string"},
                "bio": {"type": "string"},
                "avatar_url": {"type": "string"},
                "cover_url": {"type": "string"},
                "location": {"type": "string"},
                "email": {"type": "string"},
                "phone": {"type": "string"},
                "is_active": {"type": "boolean", "default": true}
            }
        },
        "UpdateAboutInput": {
            "type": "object",
            "properties": {
                "full_name": {"type": "string"},
                "nickname": {"type": "string"},
                "role": {"type": "string"},
                "bio": {"type": "string"},
                "avatar_url": {"type": "string"},
                "cover_url": {"type": "string"},
                "location": {"type": "string"},
                "email": {"type": "string"},
                "phone": {"type": "string"},
                "is_active": {"type": "boolean"}
            }
        },
        "Blog": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "title": {"type": "string"},
                "slug": {"type": "string"},
                "excerpt": {"type": "string"},
                "content": {"type": "string"},
                "cover_image": {"type": "string"},
                "is_published": {"type": "boolean"},
                "is_featured": {"type": "boolean"},
                "published_at": {"type": "string", "format": "date-time"},
                "sort_order": {"type": "integer"},
                "tags": {"type": "array", "items": {"$ref": "#/definitions/Tag"}},
                "created_at": {"type": "string", "format": "date-time"},
                "updated_at": {"type": "string", "format": "date-time"}
            }
        },
        "BlogListResponse": {
            "type": "object",
            "properties": {
                "success": {"type": "boolean"},
                "data": {"type": "array", "items": {"$ref": "#/definitions/Blog"}},
                "pagination": {"$ref": "#/definitions/Pagination"}
            }
        },
        "CreateBlogInput": {
            "type": "object",
            "required": ["title", "slug"],
            "properties": {
                "title": {"type": "string"},
                "slug": {"type": "string"},
                "excerpt": {"type": "string"},
                "content": {"type": "string"},
                "cover_image": {"type": "string"},
                "is_published": {"type": "boolean"},
                "is_featured": {"type": "boolean"},
                "published_at": {"type": "string", "format": "date-time"},
                "sort_order": {"type": "integer"},
                "tag_ids": {"type": "array", "items": {"type": "string"}}
            }
        },
        "UpdateBlogInput": {
            "type": "object",
            "properties": {
                "title": {"type": "string"},
                "slug": {"type": "string"},
                "excerpt": {"type": "string"},
                "content": {"type": "string"},
                "cover_image": {"type": "string"},
                "is_published": {"type": "boolean"},
                "is_featured": {"type": "boolean"},
                "published_at": {"type": "string", "format": "date-time"},
                "sort_order": {"type": "integer"},
                "tag_ids": {"type": "array", "items": {"type": "string"}}
            }
        },
        "Certificate": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "issuer": {"type": "string"},
                "issue_date": {"type": "string", "format": "date-time"},
                "expiry_date": {"type": "string", "format": "date-time"},
                "credential_id": {"type": "string"},
                "credential_url": {"type": "string"},
                "image_url": {"type": "string"},
                "description": {"type": "string"},
                "sort_order": {"type": "integer"},
                "created_at": {"type": "string", "format": "date-time"},
                "updated_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateCertificateInput": {
            "type": "object",
            "required": ["name", "issuer", "issue_date"],
            "properties": {
                "name": {"type": "string"},
                "issuer": {"type": "string"},
                "issue_date": {"type": "string", "format": "date-time"},
                "expiry_date": {"type": "string", "format": "date-time"},
                "credential_id": {"type": "string"},
                "credential_url": {"type": "string"},
                "image_url": {"type": "string"},
                "description": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        },
        "UpdateCertificateInput": {
            "type": "object",
            "properties": {
                "name": {"type": "string"},
                "issuer": {"type": "string"},
                "issue_date": {"type": "string", "format": "date-time"},
                "expiry_date": {"type": "string", "format": "date-time"},
                "credential_id": {"type": "string"},
                "credential_url": {"type": "string"},
                "image_url": {"type": "string"},
                "description": {"type": "string"},
                "sort_order": {"type": "integer"}
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "JWT Authorization header using Bearer scheme. Example: \"Bearer {token}\""
        }
    }
}`

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api.ulumfr.my.id",
	BasePath:         "/",
	Schemes:          []string{"https"},
	Title:            "UlumFR Portfolio API",
	Description:      "Backend API for Portfolio CMS with JWT Authentication",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
