# Database Design Documentation

## Overview
The database is designed to support a subscription-based system with promotional campaigns and voucher management capabilities. It uses PostgreSQL and leverages ENUMs, foreign keys, and constraints to maintain data integrity.

## Core Data Types (ENUMs)

```sql
campaign_user_type: 'internal' | 'external' | 'both'
package_type: 'bronze' | 'silver' | 'gold'
voucher_status: 'active' | 'processing' | 'used' | 'expired'
subscription_status: 'active' | 'expired' | 'cancelled' | 'pending'
payment_status: 'pending' | 'completed' | 'failed' | 'refunded'
discount_type: 'percentage' | 'fixed'
discount_entity_type: 'package' | 'tour' | 'merchandise'
```

## Database Schema

### Core Tables

1. **Users & Authentication**
   - `roles`: Defines user roles (admin, user)
   - `users`: Stores user information and authentication details

2. **Subscription Management**
   - `packages`: Defines available subscription packages
   - `subscriptions`: Tracks user subscriptions and their status

3. **Campaign System**
   - `promotional_campaigns`: Main campaign configuration
   - `campaign_entities`: Links campaigns to specific entities (packages, tours, etc.)
   - `campaign_platform_limits`: Platform-specific campaign settings
   - `campaign_tracking`: Tracks campaign performance and user interactions

4. **Voucher System**
   - `vouchers`: Manages voucher generation and usage
   - Links to users, campaigns, and tracking information

### Key Features

1. **Entity Type Management**
   - Flexible entity system allowing campaigns to target different types of items
   - Support for packages, tours, and merchandise
   - Extensible for future entity types

2. **Platform Integration**
   - Support for multiple marketing platforms (Facebook, Google, Twitter, etc.)
   - Platform-specific voucher limits
   - Tracking capabilities per platform

3. **Campaign Tracking**
   - IP-based tracking for guest users
   - User agent and referrer tracking
   - Conversion tracking

4. **Subscription Management**
   - Multiple package tiers (Bronze, Silver, Gold)
   - Flexible duration settings
   - Price and discount management

### Key Constraints & Features

1. **Data Integrity**
   ```sql
   -- Ensures valid voucher counts
   CHECK (remaining_vouchers >= 0 AND max_vouchers > 0)
   
   -- Prevents platform limit overflow
   CHECK (used_count <= COALESCE(voucher_limit, used_count))
   
   -- Prevents multiple active subscriptions
   UNIQUE INDEX on subscriptions (user_id) WHERE status IN ('active', 'pending')
   ```

2. **Tracking & Analytics**
   - Full timestamp tracking on all relevant tables
   - Platform-specific performance metrics
   - User interaction tracking

3. **Flexibility**
   - Support for both percentage and fixed-amount discounts
   - Configurable voucher limits per platform
   - Extensible entity system for future expansion

### Security Features

1. **User Management**
   - Password hashing
   - Role-based access control
   - Status tracking for account management

2. **Campaign Security**
   - IP tracking for voucher distribution
   - Unique voucher codes
   - Platform-specific hashing for tracking

### Future Extensibility

The database is designed to be easily extended for:
1. New entity types (beyond packages, tours, merchandise)
2. Additional marketing platforms
3. New discount types
4. Enhanced tracking capabilities

### Performance Considerations

1. **Indexes**
   - Primary keys on all tables
   - Foreign key relationships for referential integrity
   - Unique constraints where necessary

2. **Constraints**
   - Check constraints for business rules
   - Unique constraints for data integrity
   - Foreign key constraints for relationships