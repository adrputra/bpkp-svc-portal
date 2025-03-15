# API Endpoints Documentation

## Project Overview
This project is a service portal for managing various functionalities related to attendance, institutions, parameters, roles, and users. It is built using the Go programming language and utilizes the Echo framework for handling HTTP requests.

## Installation Instructions
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/adrputra/bpkp-svc-portal
   cd bpkp-svc-portal
   ```

2. **Install Dependencies**:
   Ensure you have Go installed on your machine. Run the following command to install the necessary dependencies:
   ```bash
   go mod tidy
   ```

3. **Configuration**:
   Update the `config.yaml` file with your database and service configurations. Here are the key settings:
   - **Database**: Update the `database` section with your MySQL credentials.
   - **Redis**: Update the `redis` section if you are using Redis for caching.
   - **Jaeger**: Configure the Jaeger settings for tracing if needed.

4. **Run the Application**:
   Start the application by running:
   ```bash
   go run main.go
   ```

5. **Access the API**:
   The API will be available at `http://localhost:8002`.

## API Endpoints

### Attendance Endpoints
- **GET /attendance**: Retrieve today's attendances.
- **POST /attendance**: Retrieve attendances for a specific user.
- **POST /attendance/checkin**: Check in a user.
- **POST /attendance/checkout**: Check out a user.

### Institution Endpoints
- **GET /institution**: Retrieve all institutions.
- **GET /institution/:id**: Retrieve details of a specific institution by ID.
- **POST /institution**: Create a new institution.
- **PUT /institution**: Update an existing institution.
- **DELETE /institution/:id**: Delete an institution by ID.

### Parameter Endpoints
- **GET /param/:id**: Retrieve a parameter by its key.
- **GET /param**: Retrieve all parameters.
- **POST /param**: Insert a new parameter.
- **PUT /param**: Update an existing parameter.
- **DELETE /param/:id**: Delete a parameter by its key.

### Role Endpoints
- **GET /role**: Retrieve all roles.
- **GET /role/mapping**: Retrieve all role mappings.
- **POST /role/create**: Create a new role.
- **PUT /role/mapping**: Update a role mapping.
- **POST /role/mapping/create**: Create a new role mapping.
- **DELETE /role/mapping/:id**: Delete a role mapping by ID.
- **GET /role/menu**: Retrieve all menus.
- **PUT /role/menu**: Update a menu.
- **POST /role/menu/create**: Create a new menu.
- **DELETE /role/menu/:id**: Delete a menu by ID.

### User Endpoints
- **GET /user**: Retrieve all users.
- **GET /user/detail/:id**: Retrieve details of a specific user by ID.
- **PUT /user**: Update user information.
- **DELETE /user/:id**: Delete a user by ID.
- **GET /user/institutions**: Retrieve a list of institutions associated with users.
- **POST /user/profile-photo**: Upload a profile photo for a user.
- **POST /user/cover-photo**: Upload a cover photo for a user.
