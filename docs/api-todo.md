# API Flaws To Do list

## Architectural Improvements (High Priority)

### 1. **Reduce Service Coupling and Optimize Data Fetching**

* [ ] Refactor `CourseService` to avoid cross-service calls for data enrichment.

    * [ ] Utilize repository joins (`Joins("Faculty").Joins("Professor")`) correctly.
    * [ ] Preload related entities to mitigate N+1 query issues.
    * [ ] Validate and benchmark queries post-refactor to ensure improved performance.

### 2. **Improve Transaction Boundaries in Repositories**

* [ ] Correct logic in batch creation endpoints:

    * [ ] Only fetch newly created courses post-transaction instead of the entire dataset.
    * [ ] Write tests covering the batch creation scenario comprehensively.

## Error Handling & Validation (High Priority)

### 3. **Standardize and Improve Error Handling**

* [ ] Consistently use custom error types (`ErrNotFound`, `ErrConflict`, `ErrInvalid`) across all services.

    * [ ] Refactor `AuthService.Register` to return `errors.NewConflictError` instead of raw errors.
    * [ ] Replace string-comparison based error handling with type-based (`errors.Is`) checks.

### 4. **Centralized Validation for Foreign Keys**

* [ ] Implement validation for existence checks:

    * [ ] During user registration (`AuthService.Register`), verify university and faculty IDs explicitly through service-level checks or repositories, rather than relying on database constraints.

## Security Enhancements (High Priority)

### 5. **Implement Rate Limiting and Protection Against Abuse**

* [ ] Add middleware or gateway-level rate limiting, particularly for:

    * [ ] Login endpoint (`AuthService.Login`).
    * [ ] Password reset and email verification endpoints.

### 6. **Enhance Refresh Token Handling**

* [ ] Fully implement refresh token rotation and revocation:

    * [ ] Utilize the `revoked` flag on refresh tokens consistently.
    * [ ] Ensure logout endpoint invalidates tokens appropriately.

## Business Logic & Service Responsibilities (Medium Priority)

### 7. **Separate Concerns in AuthService**

* [ ] Decouple side effects like email sending from synchronous request-response cycles:

    * [ ] Move email sending logic to asynchronous tasks or background workers.

### 8. **Consistent Transaction Handling**

* [ ] Ensure robust database transaction handling:

    * [ ] Handle edge cases such as race conditions (use unique constraints in DB and catch corresponding errors).

## REST API & Endpoint Design (Medium Priority)

### 9. **Ensure RESTful Endpoint Consistency**

* [ ] Adjust inconsistent endpoint definitions:

    * [x] Correct route from `/v1/admin/faculties/courses/:facultyID` to `/v1/admin/faculties/{facultyID}/courses`.
    * [ ] Clarify access control for `/v1/courses` endpoint. Decide if it should remain admin-only or open to regular users.

### 10. **Improve Endpoint Naming Conventions**

* [ ] Refactor user course selection endpoints for RESTful design:

    * [ ] Change `POST /v1/courses/select` to `POST /v1/user/courses`.
    * [ ] Change `DELETE /v1/courses/select/{courseId}` to `DELETE /v1/user/courses/{courseId}`.

## Input Validation & Error Detail Management (Medium Priority)

### 11. **Standardize Input Validation Messages**

* [ ] Ensure all input validation errors are user-friendly and generic:

    * [ ] Abstract and standardize JSON decoding errors from `ShouldBindJSON`.

### 12. **Uniform Error Responses**

* [ ] Adjust error handling in handlers to avoid generic 500 Internal Server Errors:

    * [ ] Explicitly handle validation errors, such as invalid UUID filters, and return 400 status codes.

## Documentation & Developer Experience (Medium Priority)

### 13. **Align Documentation with Implementation**

* [ ] Audit API documentation and implementation thoroughly for inconsistencies:

    * [ ] Fix mismatches between documented endpoints and actual code implementation.

## Logging and Observability (Low Priority)

### 14. **Enhance Logging Context**

* [ ] Add request or correlation IDs in logs to trace errors and exceptions easily without exposing internal details to users.

### 15. **Reduce Redundant Code**

* [ ] Abstract repetitive UUID parsing and error-mapping logic:

    * [ ] Create middleware for parsing and validating common request parameters.


## Code Maintenance & Stylistic Consistency (Low Priority)

### 16. **Ensure Endpoint Naming Conventions**

* [ ] Standardize singular vs. plural endpoint names (e.g., faculty vs. faculties):

    * [ ] Consider using consistently plural forms across the API.


# API Flaws Full Analysis

Open [api-architectural-and-design-flaws.pdf](api-architectural-and-design-flaws.pdf)