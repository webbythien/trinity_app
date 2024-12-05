# Trinity App

A subscription management system with campaign and voucher capabilities, built with Go and PostgreSQL.

## Overview

Trinity App is a robust subscription management platform that handles user subscriptions, promotional campaigns, and voucher distribution across multiple marketing platforms. The system supports different subscription tiers (Bronze, Silver, Gold) and provides flexible campaign management capabilities.

## Key Features

- **User Authentication**
  - Registration and login system
  - Role-based access control (Admin/User)

- **Subscription Management**
  - Multiple package tiers (Bronze, Silver, Gold)
  - Flexible subscription durations
  - Pricing management

- **Campaign Management**
  - Percentage or fixed amount discounts
  - Platform-specific campaign tracking (Facebook, TikTok, YouTube, etc.)
  - Voucher limit control per platform
  - Support for different entity types (packages, tours, merchandise)

- **Voucher System**
  - Unique voucher code generation
  - IP-based tracking for guest users
  - Usage tracking and validation
  - Time-based validity

## Technical Architecture

### Database Schema
- PostgreSQL with comprehensive schema design
- ENUM types for consistent data representation
- Foreign key constraints for data integrity
- Unique constraints to prevent duplicate subscriptions

### Security Features
- Password hashing
- Role-based access control
- IP tracking for voucher distribution

## Technical Decisions

1. **PostgreSQL as Primary Database**
   - Strong ACID compliance for financial transactions
   - Rich data types (ENUM, TIMESTAMP WITH TIME ZONE)
   - Complex query capabilities for reporting

2. **Docker-based Deployment**
   - Consistent development environment
   - Easy local setup
   - Production-ready configuration

3. **Modular Schema Design**
   - Separate tables for different entities (users, packages, campaigns)
   - Junction tables for many-to-many relationships
   - Extensible design for future features

## Assumptions

1. **Business Rules**
   - Users can only have one active subscription at a time
   - Vouchers are single-use only
   - Campaign discounts cannot exceed the maximum discount amount

2. **Technical Assumptions**
   - Single timezone (UTC) for all timestamps
   - Unique email addresses for users
   - Internet connectivity for platform tracking

## Local Setup Guide

1. **Prerequisites**
   ```bash
   # Required software
   - Docker and Docker Compose
   - Go 1.20 or higher
   - Git
   ```

2. **Clone and Setup**
   ```bash
   # Clone repository
   git clone https://github.com/webbythien/trinity_app.git
   cd trinity_app

   # Start database
   cd db
   docker compose up -d --build

   # Setup API
   cd ../core-api
   cp .env.example .env
   go mod tidy
   ```

3. **Run Application**
   ```bash
   go run main.go server
   ```

4. **Access Points**
   - API: http://localhost:5005/
   - Swagger: http://localhost:5005/swagger/index.html

## Planned Improvements

1. **Performance Optimizations**
   - Redis caching layer for frequently accessed data
   - Message queue (RabbitMQ) for subscription processing
   - Database indexing optimization

2. **Architecture Enhancements**
   - Microservices architecture for better scalability
   - Separate services for:
     - Authentication
     - Campaign management
     - Subscription processing
     - Voucher distribution

3. **Feature Additions**
   - Multi-tenant support
   - Advanced analytics dashboard
   - Automated campaign performance reporting
   - Integration with more marketing platforms

4. **Technical Debt**
   - Implement comprehensive test coverage
   - Add API rate limiting
   - Enhanced logging and monitoring
   - CI/CD pipeline setup


## License

[MIT License](LICENSE)