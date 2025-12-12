# Logging Best Practices and Conventions for POS Cafe System

## Overview
This document outlines the logging best practices and conventions for the POS Cafe system. The system uses structured logging with Logrus to provide consistent, searchable, and actionable logs in both development and production environments.

## Logger Configuration

### Environment-Based Configuration
The logger is configured differently based on the environment:

- **Development**: Uses text formatter with colors for better readability
- **Production**: Uses JSON formatter for better machine parsing and log aggregation

### Log Levels
- `debug`: Detailed information for debugging purposes
- `info`: General information about application flow
- `warn`: Warning about potential issues
- `error`: Error events that don't prevent application flow
- `fatal`: Critical errors that cause application termination
- `panic`: Errors that cause panic

## Structured Logging

All log messages include structured fields to facilitate searching and analysis:

```go
utils.LogInfo("User login event", map[string]any{
    "user_id": "12345",
    "username": "john_doe",
    "ip_address": "192.168.1.1",
    "timestamp": time.Now().Format(time.RFC3339),
})
```

### Standard Fields
When applicable, logs should include these standard fields:
- `timestamp`: RFC3339 formatted timestamp
- `request_id`: Unique identifier for the request
- `user_id`: User identifier
- `username`: Username
- `ip_address`: Client IP address
- `session_id`: Session identifier
- `action`: The action being performed
- `resource`: The resource being accessed

## Logging Conventions

### Log Levels Usage

#### Info Level
Use for:
- Successful operations
- State changes
- Significant business events
- Application startup/shutdown

```go
utils.LogInfo("Starting POS Cafe server", map[string]any{
    "environment": "production",
    "port":        "8080",
    "version":     "3.0.0",
})
```

#### Error Level
Use for:
- Failed operations
- Unexpected conditions
- External service failures

```go
utils.LogError("Database connection failed", map[string]any{
    "error": err.Error(),
    "host": config.DB.Host,
    "port": config.DB.Port,
})
```

#### Warn Level
Use for:
- Potential issues that don't stop execution
- Deprecated feature usage
- Security-related events that don't block operation

```go
utils.LogWarn("User login failed", map[string]any{
    "username": "invalid_user",
    "ip_address": "192.168.1.100",
    "reason": "invalid_credentials",
})
```

### Sensitive Information
Never log sensitive information such as:
- Passwords
- Credit card numbers
- JWT tokens
- Social security numbers
- Personal identification numbers

## Audit Logging

The system includes a specialized audit logger for tracking sensitive operations:

### Events Tracked
- User login/logout
- Data modifications
- Sensitive operations (order creation, payments)
- Financial operations

### Audit Event Structure
```go
type AuditEvent struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time              `json:"timestamp"`
    UserID      string                 `json:"user_id,omitempty"`
    Username    string                 `json:"username,omitempty"`
    Action      string                 `json:"action"`
    Resource    string                 `json:"resource"`
    Details     map[string]interface{} `json:"details,omitempty"`
    IPAddress   string                 `json:"ip_address,omitempty"`
    UserAgent   string                 `json:"user_agent,omitempty"`
    Status      string                 `json:"status"` // success, failed, warning
    SessionID   string                 `json:"session_id,omitempty"`
    RequestID   string                 `json:"request_id,omitempty"`
    Environment string                 `json:"environment"`
}
```

## Middleware Logging

The request logging middleware automatically logs:
- Incoming requests with method, URI, remote address, user agent
- Completed requests with status, response size, duration
- Request IDs for tracking request flows

## Performance Considerations

- Be mindful of logging frequency in tight loops
- Avoid logging large data structures unnecessarily
- Use log level checks in performance-critical paths:
  ```go
  if utils.Logger.GetLevel() >= logrus.DebugLevel {
      // Expensive operation to prepare log data
      utils.LogDebug("Detailed info", expensiveData)
  }
  ```

## Environment-Specific Practices

### Development
- Use text format for easier reading
- Include more verbose debug information
- Use colors to distinguish log levels

### Production
- Use JSON format for log aggregation systems
- Avoid logging sensitive information
- Ensure logs include sufficient context for debugging
- Consider log volume and retention policies

## Integration with External Systems

The JSON logging format is compatible with:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Fluentd/Fluent Bit
- AWS CloudWatch
- Google Cloud Logging
- Datadog
- Other centralized logging solutions