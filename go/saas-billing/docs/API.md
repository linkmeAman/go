# SaaS Billing System API Documentation

## Base URL
- Development: `http://localhost:8080`
- Production: `https://api.example.com`

## Authentication
Most endpoints require JWT authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Response Format
All API responses follow this standard format:
```json
{
  "success": true,
  "data": {},
  "error": null,
  "metadata": {
    "timestamp": "2025-09-07T10:00:00Z",
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_pages": 5,
      "total_records": 48,
      "has_next": true,
      "has_previous": false
    }
  }
}
```

## Error Format
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": "Additional error details",
    "timestamp": "2025-09-07T10:00:00Z",
    "request_id": "req_123abc"
  }
}
```

## API Endpoints

### Authentication

#### Register User
- **POST** `/api/v1/auth/register`
- **Description**: Register a new user
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123"
  }
  ```
- **Response (201)**:
  ```json
  {
    "success": true,
    "data": {
      "user_id": "uuid",
      "email": "user@example.com"
    }
  }
  ```

#### Login
- **POST** `/api/v1/auth/login`
- **Description**: Authenticate user and get JWT token
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123"
  }
  ```
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": {
      "token": "jwt_token_here",
      "expires_at": "2025-09-07T11:00:00Z"
    }
  }
  ```

### Organizations

#### Create Organization
- **POST** `/api/v1/organizations`
- **Auth**: Required
- **Description**: Create a new organization
- **Request Body**:
  ```json
  {
    "name": "My Company",
    "description": "Optional description"
  }
  ```
- **Response (201)**:
  ```json
  {
    "success": true,
    "data": {
      "id": "org_uuid",
      "name": "My Company",
      "description": "Optional description",
      "created_at": "2025-09-07T10:00:00Z"
    }
  }
  ```

#### List Organizations
- **GET** `/api/v1/organizations`
- **Auth**: Required
- **Description**: List organizations user has access to
- **Query Parameters**:
  - `page` (int, default: 1)
  - `page_size` (int, default: 10)
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": [
      {
        "id": "org_uuid",
        "name": "My Company",
        "role": "admin",
        "member_count": 5
      }
    ],
    "metadata": {
      "pagination": {
        "current_page": 1,
        "page_size": 10,
        "total_pages": 1,
        "total_records": 1
      }
    }
  }
  ```

#### Add Organization Member
- **POST** `/api/v1/organizations/:orgID/members`
- **Auth**: Required (admin only)
- **Description**: Add a member to organization
- **Request Body**:
  ```json
  {
    "user_id": "user_uuid",
    "role": "member"
  }
  ```
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": {
      "user_id": "user_uuid",
      "organization_id": "org_uuid",
      "role": "member"
    }
  }
  ```

### Billing

#### Get Plans
- **GET** `/api/v1/organizations/:orgID/billing/plans`
- **Auth**: Required
- **Description**: List available subscription plans
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": [
      {
        "id": "plan_uuid",
        "name": "Pro",
        "description": "Professional plan",
        "price": 49.99,
        "billing_interval": "monthly",
        "features": {
          "users": 20,
          "storage": "10GB",
          "api_calls": 10000
        }
      }
    ]
  }
  ```

#### Subscribe to Plan
- **POST** `/api/v1/organizations/:orgID/billing/subscribe/:planID`
- **Auth**: Required (admin only)
- **Description**: Subscribe organization to a plan
- **Request Body**:
  ```json
  {
    "payment_method_id": "pm_123",
    "billing_interval": "monthly"
  }
  ```
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": {
      "subscription_id": "sub_uuid",
      "plan_id": "plan_uuid",
      "status": "active",
      "current_period_end": "2025-10-07T10:00:00Z"
    }
  }
  ```

#### Get Current Subscription
- **GET** `/api/v1/organizations/:orgID/billing/subscription`
- **Auth**: Required
- **Description**: Get organization's current subscription
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": {
      "subscription_id": "sub_uuid",
      "plan": {
        "id": "plan_uuid",
        "name": "Pro"
      },
      "status": "active",
      "current_period_start": "2025-09-07T10:00:00Z",
      "current_period_end": "2025-10-07T10:00:00Z",
      "cancel_at_period_end": false
    }
  }
  ```

#### Get Invoices
- **GET** `/api/v1/organizations/:orgID/billing/invoices`
- **Auth**: Required
- **Description**: List organization's invoices
- **Query Parameters**:
  - `page` (int, default: 1)
  - `page_size` (int, default: 10)
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": [
      {
        "id": "inv_uuid",
        "amount": 49.99,
        "status": "paid",
        "created_at": "2025-09-07T10:00:00Z",
        "paid_at": "2025-09-07T10:00:00Z"
      }
    ],
    "metadata": {
      "pagination": {
        "current_page": 1,
        "page_size": 10,
        "total_pages": 1,
        "total_records": 1
      }
    }
  }
  ```

### Usage Tracking

#### Record Usage
- **POST** `/api/v1/organizations/:orgID/usage`
- **Auth**: Required
- **Description**: Record usage for an organization
- **Request Body**:
  ```json
  {
    "metric": "api_calls",
    "quantity": 1,
    "timestamp": "2025-09-07T10:00:00Z"
  }
  ```
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": {
      "usage_id": "usage_uuid",
      "metric": "api_calls",
      "quantity": 1,
      "recorded_at": "2025-09-07T10:00:00Z"
    }
  }
  ```

#### Get Usage Report
- **GET** `/api/v1/organizations/:orgID/usage`
- **Auth**: Required
- **Description**: Get organization's usage report
- **Query Parameters**:
  - `start_date` (ISO date)
  - `end_date` (ISO date)
  - `metric` (string, optional)
- **Response (200)**:
  ```json
  {
    "success": true,
    "data": {
      "period": {
        "start": "2025-09-01T00:00:00Z",
        "end": "2025-09-07T23:59:59Z"
      },
      "metrics": {
        "api_calls": {
          "total": 5000,
          "limit": 10000,
          "usage_percentage": 50
        },
        "storage": {
          "total": 5368709120,
          "limit": 10737418240,
          "usage_percentage": 50
        }
      }
    }
  }
  ```

## Rate Limits
- 100 requests per minute per IP address
- 1000 requests per minute per authenticated user
- Endpoints return `429 Too Many Requests` when limit is exceeded

## Error Codes
- `INVALID_INPUT`: Request validation failed
- `UNAUTHORIZED`: Authentication required
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `RATE_LIMITED`: Too many requests
- `INTERNAL_SERVER_ERROR`: Unexpected server error

## Webhooks
Webhooks are available for the following events:
- `subscription.created`
- `subscription.updated`
- `subscription.cancelled`
- `invoice.created`
- `invoice.paid`
- `usage.threshold_reached`

## Best Practices
1. Always include proper authentication headers
2. Implement proper error handling
3. Use pagination for list endpoints
4. Monitor rate limits
5. Store and log request IDs for debugging
