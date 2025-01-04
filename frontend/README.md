# Clean up
cd ~/projects/FleetFlow
rm -rf frontend

# Create a new directory
mkdir frontend
cd frontend

# Initialize package.json
npm init -y

# Install React 18 and core dependencies
npm install react@18.2.0 react-dom@18.2.0 typescript@4.9.5 @types/react@18.2.0 @types/react-dom@18.2.0

# Install development dependencies
npm install --save-dev @typescript-eslint/eslint-plugin @typescript-eslint/parser eslint eslint-plugin-react

# Install our app dependencies
npm install @reduxjs/toolkit@1.9.7 react-redux@8.1.3 axios@1.6.2 react-router-dom@6.20.1 @types/react-router-dom@5.3.3 jwt-decode@4.0.0 @types/jwt-decode@3.1.0 @emotion/react@11.11.1 @emotion/styled@11.11.0

Dashboards to create:
Driver Dashboard
- Vehicle Statistics
- Latest Orders/Deliveries
- Tracking Information
- Performance Metrics
- Support Chat*

Fleet Manager Dashboard
- Fleet Overview
- Vehicle Status Summary
- Driver Performance Analytics
- Delivery Statistics
- Route Optimization Insights

Loader Dashboard
- Warehouse Capacity
- Package Processing Queue
- Loading Schedule
- Performance Metrics
- Inventory Status

Customer Dashboard
- Order Tracking
- Delivery History
- Package Status
- Feedback History
- Support Access

Admin Dashboard
- System Overview
- User Management
- Fleet Analytics
- Warehouse Status
- Performance Reports

Manager Dashboard
- Operations Overview
- Performance Analytics
- Resource Utilization
- Cost Analysis
- Team Performance