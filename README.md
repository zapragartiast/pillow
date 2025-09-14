# User Management System

A comprehensive user management system built with Go backend and PostgreSQL database, designed to handle user authentication, role-based access control (RBAC), organization management, and audit logging.

## Project Overview

This project implements a full-featured User Management System with the following key components:

- **Backend**: Go application using Gorilla Mux for routing and PostgreSQL for data storage
- **Database**: Comprehensive schema with users, roles, permissions, organizations, and audit logging
- **Frontend**: To be implemented (React/Vue/Angular framework)
- **Authentication**: JWT-based authentication system
- **Authorization**: Role-based access control with granular permissions
- **Organizations**: Multi-tenant organization support with hierarchical structure

## Database Schema

The system uses PostgreSQL with the following main tables:

- **Users**: User accounts with authentication details
- **Roles**: System roles (e.g., admin, user, manager)
- **Permissions**: Granular permissions (e.g., read_user, create_user)
- **Organizations**: Multi-tenant organization support
- **Role_Permissions**: Many-to-many relationship between roles and permissions
- **User_Roles**: User-role assignments with scope support
- **User_Organizations**: User-organization memberships
- **Audit_Log**: Comprehensive audit trail for all user actions

All tables use UUID for primary keys and include proper foreign key relationships and indexes.

## Current Implementation Status

### Backend (Go)
- ✅ Basic database connection setup
- ✅ User model and basic CRUD operations (GET, POST)
- ✅ PostgreSQL integration with lib/pq driver
- ❌ Authentication system (JWT, login/logout)
- ❌ Role and permission management
- ❌ Organization management
- ❌ Audit logging
- ❌ Input validation and error handling
- ❌ Security middleware

### Frontend
- ❌ Not implemented yet
- ❌ Authentication UI
- ❌ User management interface
- ❌ Role/permission management UI
- ❌ Organization management UI

### Database
- ✅ Complete schema design
- ✅ Foreign key relationships
- ✅ Indexes for performance
- ❌ Data seeding/migration scripts

## Roadmap

### Phase 1: Backend Foundation
1. **Analyze current codebase and database schema for gaps and inconsistencies**
   - Review existing code for issues
   - Identify data type mismatches and architectural problems

2. **Fix data type mismatches**
   - Update User model from int ID to UUID
   - Ensure all models align with database schema
   - Update handlers to use UUID types

3. **Implement JWT-based authentication system**
   - Add login endpoint with password verification
   - Implement JWT token generation and validation
   - Add logout functionality
   - Create authentication middleware

4. **Add password hashing and validation utilities**
   - Implement secure password hashing (bcrypt)
   - Add password strength validation
   - Create password reset functionality

5. **Create models for Roles, Permissions, Organizations**
   - Define Go structs for all database entities
   - Implement relationships and associations
   - Add JSON serialization tags

### Phase 2: Core Features
6. **Implement role-based access control (RBAC) middleware**
   - Create authorization middleware
   - Implement permission checking
   - Add role hierarchy support

7. **Develop API endpoints for user management**
   - Expand CRUD operations (GET, POST, PUT, DELETE)
   - Add role/permission assignment to users
   - Implement user search and filtering

8. **Develop API endpoints for role and permission management**
   - CRUD operations for roles
   - CRUD operations for permissions
   - Role-permission assignment endpoints

9. **Develop API endpoints for organization management**
   - Organization CRUD operations
   - User-organization relationship management
   - Hierarchical organization support

10. **Implement audit logging for all user actions**
    - Log all CRUD operations
    - Track user authentication events
    - Implement audit trail queries

11. **Add input validation and error handling**
    - Implement request validation middleware
    - Add comprehensive error responses
    - Handle edge cases and invalid inputs

### Phase 3: Frontend Development
12. **Create frontend framework**
    - Choose and set up React/Vue/Angular
    - Configure build tools and project structure
    - Set up routing and state management

13. **Build user authentication UI**
    - Login and registration forms
    - Password reset functionality
    - User profile management

14. **Build user management UI**
    - User list with search and filtering
    - User creation and editing forms
    - Bulk user operations

15. **Build role and permission management UI**
    - Role creation and management
    - Permission assignment interface
    - Role hierarchy visualization

16. **Build organization management UI**
    - Organization structure management
    - User-organization assignment
    - Organization settings

17. **Integrate frontend with backend APIs**
    - Implement API client
    - Handle authentication tokens
    - Add error handling and loading states

### Phase 4: Testing and Quality Assurance
18. **Add unit and integration tests for backend**
    - Unit tests for handlers and models
    - Integration tests for API endpoints
    - Database testing utilities

19. **Add end-to-end tests for frontend**
    - Test complete user workflows
    - Authentication flow testing
    - Cross-browser compatibility

### Phase 5: Security and Production
20. **Implement security best practices**
    - CORS configuration
    - Rate limiting
    - Input sanitization
    - Security headers

21. **Add logging and monitoring**
    - Application logging
    - Performance monitoring
    - Error tracking

22. **Create deployment configuration**
    - Docker containerization
    - CI/CD pipeline setup
    - Environment configuration

23. **Write documentation**
    - API documentation (Swagger/OpenAPI)
    - User guides and deployment instructions
    - Code documentation

## Technology Stack

### Backend
- **Language**: Go 1.25.1
- **Framework**: Gorilla Mux
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Password Hashing**: bcrypt

### Frontend (Planned)
- **Framework**: React/Vue/Angular (TBD)
- **State Management**: Redux/Context API
- **Styling**: CSS-in-JS/Tailwind CSS
- **Build Tool**: Vite/Webpack

### DevOps
- **Containerization**: Docker
- **CI/CD**: GitHub Actions/Jenkins
- **Monitoring**: Prometheus/Grafana

## Getting Started

### Prerequisites
- Go 1.25.1 or later
- PostgreSQL 12+
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd pillow
   ```

2. Set up PostgreSQL database:
   ```bash
   createdb pillowdb
   psql -d pillowdb -f database/pillowdb.sql
   ```

3. Configure environment variables:
   ```bash
   cd backend
   cp .env.example .env
   # Edit .env file with your configuration
   ```

   Required environment variables:
   - `DATABASE_URL`: PostgreSQL connection string
   - `JWT_SECRET`: Secret key for JWT token signing (minimum 32 characters)
   - `SERVER_PORT`: Port for the server to listen on (default: 8080)

4. Install dependencies and run the backend:
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

## API Endpoints (Current)

### API Endpoints (/api group)
- `POST /api/register` - Register new user
- `POST /api/login` - User login
- `GET /api/users` - Get all active users
- `PUT /api/users/{id}` - Update user (planned)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.