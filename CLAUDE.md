# Claude Code Review Guidelines

This file defines the code review standards and guidelines for the ficcc-backend project.

## Review Focus Areas

### Code Quality
- Ensure code is readable, maintainable, and follows established patterns
- Check for proper error handling and edge cases
- Verify that functions and methods have single responsibilities
- Look for code duplication and suggest refactoring opportunities

### Security
- Check for common vulnerabilities (SQL injection, XSS, authentication issues, etc.)
- Verify proper input validation and sanitization
- Ensure sensitive data is properly protected
- Check for secure API endpoint implementations
- Verify proper authentication and authorization checks

### Performance
- Identify potential performance bottlenecks
- Check for inefficient database queries (N+1 queries, missing indexes)
- Look for unnecessary loops or redundant operations
- Verify proper use of caching where appropriate

### Testing
- Ensure new features have appropriate test coverage
- Check that tests are meaningful and test the right things
- Verify edge cases are covered in tests
- Look for brittle tests that might break easily

### Documentation
- Check that complex logic is properly commented
- Verify API endpoints are documented
- Ensure function/method signatures are clear or documented
- Look for outdated comments that don't match the code

### Best Practices
- Follow consistent coding style throughout the project
- Use meaningful variable and function names
- Keep functions and files at a reasonable size
- Ensure proper dependency management
- Check for proper logging and monitoring hooks

## Review Tone
- Be constructive and helpful in feedback
- Explain the "why" behind suggestions
- Offer alternative approaches when suggesting changes
- Acknowledge good practices and improvements
