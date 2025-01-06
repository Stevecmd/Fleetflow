-- Connect to postgres database first
\c postgres;

-- Drop the database if it exists
SELECT 'DROP DATABASE "fleetflow"' \gexec

-- Create the database if it doesn't exist
SELECT 'CREATE DATABASE "fleetflow"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'fleetflow')\gexec

-- Create the postgres user if it doesn't exist
DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres_user') THEN
      CREATE USER postgres_user WITH PASSWORD 'postgres_password';
   END IF;
END
$do$;

-- Connect to the application database
\c "fleetflow";

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert available roles
INSERT INTO roles (id, name, description) VALUES
    (1, 'admin', 'System administrator with full access to all features'),
    (2, 'fleet_manager', 'Manages fleet operations and vehicle assignments'),
    (3, 'driver', 'Operates vehicles and manages deliveries'),
    (4, 'loader', 'Handles warehouse operations and cargo loading'),
    (5, 'customer', 'End user who requests deliveries'),
    (6, 'manager', 'General manager with administrative privileges')
ON CONFLICT (id) DO NOTHING;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role_id INTEGER NOT NULL REFERENCES roles(id),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(50),
    date_of_birth DATE DEFAULT '1970-01-01',
    gender VARCHAR(50) DEFAULT 'unknown',
    nationality VARCHAR(100) DEFAULT 'unknown',
    preferred_language VARCHAR(50) DEFAULT 'unknown',
    status VARCHAR(50) DEFAULT 'active',
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    profile_image_url TEXT DEFAULT ''
);
-- Create addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    street1 VARCHAR(255) NOT NULL,
    street2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    zip VARCHAR(20),
    country VARCHAR(100) NOT NULL,
    address_type VARCHAR(50),
    is_default BOOLEAN DEFAULT false,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create emergency_contacts table
CREATE TABLE IF NOT EXISTS emergency_contacts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    relationship VARCHAR(100) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create warehouses table
CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    street1 VARCHAR(255) NOT NULL,
    street2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    zip VARCHAR(20),
    country VARCHAR(100) NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    capacity INTEGER,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create vehicle_statuses table
CREATE TABLE IF NOT EXISTS vehicle_statuses (
    id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create driver_statuses table
CREATE TABLE IF NOT EXISTS driver_statuses (
    id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create delivery_statuses table
CREATE TABLE IF NOT EXISTS delivery_statuses (
    id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create vehicles table with comprehensive tracking fields
CREATE TABLE IF NOT EXISTS vehicles (
    id SERIAL PRIMARY KEY,
    plate_number VARCHAR(20) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    capacity DECIMAL(10,2) NOT NULL,  -- in kg or cubic meters
    fuel_type VARCHAR(50) NOT NULL,
    status_id INTEGER NOT NULL REFERENCES vehicle_statuses(id),
    gps_unit_id VARCHAR(100) UNIQUE,
    last_maintenance TIMESTAMP WITH TIME ZONE,
    next_maintenance TIMESTAMP WITH TIME ZONE,
    mileage DECIMAL(10,2),
    insurance_expiry DATE,
    current_location_latitude DECIMAL(10,8),
    current_location_longitude DECIMAL(11,8),
    current_location_updated_at TIMESTAMP WITH TIME ZONE,
    fuel_efficiency_rating DECIMAL(5,2),
    total_fuel_consumption DECIMAL(10,2),
    total_maintenance_cost DECIMAL(10,2),
    vehicle_images JSONB, -- Array of image URLs for different angles
    registration_document_image_url TEXT,
    insurance_document_image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create driver_profiles table with essential fields
CREATE TABLE IF NOT EXISTS driver_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    license_number VARCHAR(50) NOT NULL UNIQUE,
    license_type VARCHAR(50),
    license_expiry DATE,
    vehicle_type VARCHAR(100),
    years_experience VARCHAR(50),
    certification TEXT[] DEFAULT ARRAY[]::TEXT[],
    status_id INTEGER REFERENCES driver_statuses(id),
    current_vehicle_id INTEGER REFERENCES vehicles(id),
    rating DECIMAL(3,2) CHECK (rating >= 0 AND rating <= 5),
    total_trips INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create maintenance_records table for vehicle service history
CREATE TABLE IF NOT EXISTS maintenance_records (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER REFERENCES vehicles(id),
    type VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    service_date TIMESTAMP WITH TIME ZONE NOT NULL,
    cost DECIMAL(10,2) NOT NULL,
    odometer_reading DECIMAL(10,2),
    performed_by VARCHAR(255) NOT NULL,
    next_service_date TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create storage_locations table
CREATE TABLE IF NOT EXISTS storage_locations (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id),
    section VARCHAR(50),
    aisle VARCHAR(50),
    shelf VARCHAR(50),
    capacity DECIMAL(10,2),
    is_occupied BOOLEAN DEFAULT FALSE,
    current_package_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create deliveries table for tracking shipments
CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    tracking_number VARCHAR(50) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES driver_profiles(id),
    vehicle_id INTEGER REFERENCES vehicles(id),
    status_id INTEGER NOT NULL REFERENCES delivery_statuses(id),
    pickup_time TIMESTAMP WITH TIME ZONE,
    delivery_time TIMESTAMP WITH TIME ZONE,
    estimated_delivery_time TIMESTAMP WITH TIME ZONE,
    actual_delivery_time TIMESTAMP WITH TIME ZONE,
    cargo_type VARCHAR(100) NOT NULL,
    cargo_weight DECIMAL(10,2),
    special_instructions TEXT,
    customer_id INTEGER REFERENCES users(id),
    customer_signature TEXT,
    pickup_warehouse_id INTEGER REFERENCES warehouses(id),
    delivery_warehouse_id INTEGER REFERENCES warehouses(id),
    pickup_latitude DECIMAL(10,8) NOT NULL,
    pickup_longitude DECIMAL(11,8) NOT NULL,
    delivery_latitude DECIMAL(10,8) NOT NULL,
    delivery_longitude DECIMAL(11,8) NOT NULL,
    payment_status VARCHAR(50) DEFAULT 'pending',
    estimated_fuel_consumption DECIMAL(10,2),
    actual_fuel_consumption DECIMAL(10,2),
    route_efficiency_score DECIMAL(5,2),
    weather_conditions TEXT,
    proof_of_delivery_image_url TEXT,
    package_condition_images JSONB,
    from_location VARCHAR(255),
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create packages table
CREATE TABLE IF NOT EXISTS packages (
    id SERIAL PRIMARY KEY,
    delivery_id INTEGER REFERENCES deliveries(id),
    tracking_number VARCHAR(50) UNIQUE NOT NULL,
    weight DECIMAL(10,2) NOT NULL,
    dimensions JSONB,
    storage_location_id INTEGER REFERENCES storage_locations(id),
    status VARCHAR(50),
    special_handling TEXT,
    package_images JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id SERIAL PRIMARY KEY,
    delivery_id INTEGER REFERENCES deliveries(id),
    customer_id INTEGER REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    due_date DATE NOT NULL,
    payment_date DATE,
    payment_method VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create location_history table for GPS tracking
CREATE TABLE IF NOT EXISTS location_history (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER REFERENCES vehicles(id),
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    altitude DECIMAL(10,2),
    speed DECIMAL(10,2),
    heading DECIMAL(5,2),
    accuracy DECIMAL(10,2),
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_latitude CHECK (latitude BETWEEN -90 AND 90),
    CONSTRAINT valid_longitude CHECK (longitude BETWEEN -180 AND 180),
    CONSTRAINT valid_heading CHECK (heading BETWEEN 0 AND 360)
);

-- Create route_history table for tracking vehicle routes
CREATE TABLE IF NOT EXISTS route_history (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER REFERENCES vehicles(id),
    user_id INTEGER REFERENCES users(id),
    delivery_id INTEGER REFERENCES deliveries(id),
    start_location_latitude DECIMAL(10,8) NOT NULL,
    start_location_longitude DECIMAL(11,8) NOT NULL,
    end_location_latitude DECIMAL(10,8) NOT NULL,
    end_location_longitude DECIMAL(11,8) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    distance_covered DECIMAL(10,2),
    fuel_consumption DECIMAL(10,2),
    route_data JSONB,  -- Stores waypoints and route details
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create driver_documents table
CREATE TABLE IF NOT EXISTS driver_documents (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(100),
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    document_url TEXT,
    verification_status VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create driver_performance_metrics table
CREATE TABLE IF NOT EXISTS driver_performance_metrics (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    metric_date DATE NOT NULL,
    deliveries_completed INTEGER DEFAULT 0,
    on_time_delivery_rate DECIMAL(5,2),
    customer_rating_avg DECIMAL(3,2),
    fuel_efficiency DECIMAL(10,2),
    safety_score DECIMAL(5,2),
    violations_count INTEGER DEFAULT 0,
    total_distance_covered DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create driver_schedules table
CREATE TABLE IF NOT EXISTS driver_schedules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    schedule_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(50),
    break_duration INTEGER, -- in minutes
    overtime_duration INTEGER, -- in minutes
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create fleet_analytics table
CREATE TABLE IF NOT EXISTS fleet_analytics (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER REFERENCES vehicles(id),
    analysis_date DATE NOT NULL,
    total_distance DECIMAL(10,2),
    fuel_consumption DECIMAL(10,2),
    fuel_cost DECIMAL(10,2),
    maintenance_cost DECIMAL(10,2),
    idle_time INTEGER, -- in minutes
    carbon_emissions DECIMAL(10,2),
    operating_cost DECIMAL(10,2),
    efficiency_score DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create delivery_feedback table
CREATE TABLE IF NOT EXISTS delivery_feedback (
    id SERIAL PRIMARY KEY,
    delivery_id INTEGER REFERENCES deliveries(id),
    customer_id INTEGER REFERENCES users(id),
    rating DECIMAL(3,2),
    feedback_text TEXT,
    timeliness_rating DECIMAL(3,2),
    driver_rating DECIMAL(3,2),
    package_condition_rating DECIMAL(3,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create warehouse_inventory table
CREATE TABLE IF NOT EXISTS warehouse_inventory (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id),
    item_name VARCHAR(255) NOT NULL,
    item_category VARCHAR(100),
    quantity INTEGER NOT NULL,
    unit VARCHAR(50),
    minimum_threshold INTEGER,
    maximum_capacity INTEGER,
    last_restocked TIMESTAMP WITH TIME ZONE,
    next_restock_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create warehouse_equipment table
CREATE TABLE IF NOT EXISTS warehouse_equipment (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id),
    equipment_type VARCHAR(100) NOT NULL,
    equipment_id VARCHAR(50) UNIQUE,
    status VARCHAR(50),
    last_maintenance DATE,
    next_maintenance DATE,
    condition_rating DECIMAL(3,2),
    notes TEXT,
    equipment_images JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loading_schedules table
CREATE TABLE IF NOT EXISTS loading_schedules (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id),
    vehicle_id INTEGER REFERENCES vehicles(id),
    scheduled_time TIMESTAMP WITH TIME ZONE NOT NULL,
    estimated_duration INTEGER, -- in minutes
    actual_duration INTEGER, -- in minutes
    status VARCHAR(50),
    dock_number VARCHAR(50),
    loader_id INTEGER REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create inventory_items table
CREATE TABLE IF NOT EXISTS inventory_items (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id),
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    min_quantity INTEGER NOT NULL,
    max_quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    last_restocked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create real_time_tracking table
CREATE TABLE IF NOT EXISTS real_time_tracking (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER REFERENCES vehicles(id),
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    speed DECIMAL(5,2), -- in km/h
    heading INTEGER, -- in degrees
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active',
    metadata JSONB
);

-- Create optimized_routes table
CREATE TABLE IF NOT EXISTS optimized_routes (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER REFERENCES vehicles(id),
    start_location_id INTEGER REFERENCES addresses(id),
    end_location_id INTEGER REFERENCES addresses(id),
    waypoints JSONB,
    estimated_duration INTEGER, -- in minutes
    estimated_distance DECIMAL(10,2), -- in kilometers
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    related_entity_type VARCHAR(50),
    related_entity_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    from_location VARCHAR(255) NOT NULL,
    to_location VARCHAR(255) NOT NULL,
    weight DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample users with role references
INSERT INTO users (username, password, email, role_id, first_name, last_name, phone, date_of_birth, gender, nationality, preferred_language, status, profile_image_url)
VALUES
    ('customer1', '$2a$10$i6dRcu0/6INsYZ170mo0COGpN/CxOVSfM9hZr53DWShQBTYIXPjO.', 'customer1@example.com', 5, 'John', 'Doe', '+1234567890', '1990-01-01', 'male', 'USA', 'en', 'active', 'https://fleetflow-assets.s3.amazonaws.com/profiles/user1.jpg'),
    ('fleet_manager1', '$2a$10$i6dRcu0/6INsYZ170mo0COGpN/CxOVSfM9hZr53DWShQBTYIXPjO.', 'fleet_manager1@example.com', 2, 'Jane', 'Smith', '+1234567891', '1985-05-15', 'female', 'Canada', 'en', 'active', 'https://fleetflow-assets.s3.amazonaws.com/profiles/user2.jpg'),
    ('driver1', '$2a$10$i6dRcu0/6INsYZ170mo0COGpN/CxOVSfM9hZr53DWShQBTYIXPjO.', 'driver1@example.com', 3, 'Alice', 'Johnson', '+1234567892', '1988-08-20', 'female', 'UK', 'en', 'active', 'https://fleetflow-assets.s3.amazonaws.com/profiles/user3.jpg'),
    ('loader1', '$2a$10$i6dRcu0/6INsYZ170mo0COGpN/CxOVSfM9hZr53DWShQBTYIXPjO.', 'loader1@example.com', 4, 'Bob', 'Brown', '+1234567893', '1992-03-10', 'male', 'Australia', 'en', 'active', 'https://fleetflow-assets.s3.amazonaws.com/profiles/user4.jpg'),
    ('admin1', '$2a$10$i6dRcu0/6INsYZ170mo0COGpN/CxOVSfM9hZr53DWShQBTYIXPjO.', 'admin1@example.com', 1, 'Eve', 'Wilson', '+1234567895', '1987-11-25', 'female', 'Germany', 'de', 'active', 'https://fleetflow-assets.s3.amazonaws.com/profiles/user5.jpg'),
    ('manager1', '$2a$10$i6dRcu0/6INsYZ170mo0COGpN/CxOVSfM9hZr53DWShQBTYIXPjO.', 'manager1@example.com', 6, 'Gloria', 'Estate', '+1234567895', '1987-11-25', 'female', 'norway', 'du', 'active', 'https://fleetflow-assets.s3.amazonaws.com/profiles/user6.jpg'),
    ('driver2', '$2a$10$i6dRcu0/6INsYZ170mo0COGpN/CxOVSfM9hZr53DWShQBTYIXPjO.', 'driver2@example.com', 3, 'Alex', 'Jackson', '+1234567892', '1988-08-20', 'male', 'UK', 'en', 'active', 'https://fleetflow-assets.s3.amazonaws.com/profiles/user3.jpg');

-- Insert sample emergency contacts
INSERT INTO emergency_contacts (user_id, name, relationship, phone, email)
VALUES
    (1, 'Jane Smith', 'Spouse', '+1-206-555-0101', 'jane.smith@example.com'),
    (2, 'John Williams', 'Brother', '+1-604-555-0102', 'john.williams@example.com'),
    (3, 'Mary Johnson', 'Sister', '+44-20-7123-4567', 'mary.johnson@example.com'),
    (4, 'Robert Brown', 'Father', '+61-2-8765-4321', 'robert.brown@example.com');
    -- (5, 'Anna Wilson', 'Mother', '+49-30-1234-5678', 'anna.wilson@example.com');

-- Insert status values
INSERT INTO vehicle_statuses (id, name, description) VALUES
    (1, 'available', 'Vehicle is available for assignments'),
    (2, 'in_maintenance', 'Vehicle is undergoing maintenance'),
    (3, 'on_route', 'Vehicle is currently on a delivery route'),
    (4, 'out_of_service', 'Vehicle is not operational');

INSERT INTO driver_statuses (id, name, description) VALUES
    (1, 'active', 'Driver is available for assignments'),
    (2, 'inactive', 'Driver is not currently working'),
    (3, 'on_leave', 'Driver is on approved leave'),
    (4, 'suspended', 'Driver privileges are temporarily suspended')
ON CONFLICT (id) DO NOTHING;

INSERT INTO delivery_statuses (id, name, description) VALUES
    (1, 'pending', 'Delivery is awaiting pickup'),
    (2, 'in_transit', 'Delivery is in progress'),
    (3, 'delivered', 'Delivery has been completed'),
    (4, 'cancelled', 'Delivery has been cancelled');

INSERT INTO vehicles (plate_number, type, make, model, year, capacity, fuel_type, status_id, gps_unit_id, last_maintenance, next_maintenance, mileage, insurance_expiry, current_location_latitude, current_location_longitude, current_location_updated_at, fuel_efficiency_rating, total_fuel_consumption, total_maintenance_cost, vehicle_images, registration_document_image_url, insurance_document_image_url)
VALUES 
('KAA123A', 'Truck', 'Isuzu', 'FVZ', 2022, 10000, 'Diesel', 1, 'GPS123', NULL, NULL, 15000, '2025-12-31', 40.7128, -74.0060, NULL, 15.5, 200.0, 500.0, '[]', NULL, NULL),
('XYZ456', 'Van', 'Toyota', 'Sienna', 2021, 800, 'Gasoline', 1, 'GPS456', NULL, NULL, 12000, '2025-12-31', 34.0522, -118.2437, NULL, 18.0, 150.0, 300.0, '[]', NULL, NULL), 
('DEF123', 'Truck', 'Ford', 'F-150', 2022, 1000, 'Gasoline', 1, 'GPS789', NULL, NULL, 15000, '2025-12-31', 40.7128, -74.0060, NULL, 15.5, 200.0, 500.0, '[]', NULL, NULL),
('JKL456', 'Van', 'Toyota', 'Sienna', 2021, 800, 'Gasoline', 1, 'GPS012', NULL, NULL, 12000, '2025-12-31', 34.0522, -118.2437, NULL, 18.0, 150.0, 300.0, '[]', NULL, NULL),
('ABC456', 'Truck', 'Ford', 'F-150', 2022, 1000, 'Gasoline', 1, 'GPS013', NULL, NULL, 12000, '2025-12-31', 34.0522, -118.2437, NULL, 18.0, 150.0, 300.0, '[]', NULL, NULL);

INSERT INTO driver_profiles (
    user_id,
    license_number, 
    license_type, 
    license_expiry, 
    vehicle_type, 
    years_experience, 
    certification, 
    status_id, 
    current_vehicle_id, 
    rating, 
    total_trips, 
    created_at, 
    updated_at
) VALUES 
(1, 'ABC12345', 'Class B', '2025-12-31', 'Truck', '5', ARRAY['CDL', 'Hazmat'], 1, 1, 4.5, 10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'XYZ98765', 'Class A', '2024-11-30', 'Van', '3', ARRAY['CDL'], 1, NULL, 4.0, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'DL123456', 'Class A', '2025-12-31', 'Truck', '5', ARRAY['CDL'], 1, 1, 4.5, 10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'QYU98765', 'Class A', '2024-11-30', 'Van', '3', ARRAY['CDL'], 1, NULL, 4.0, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(5, 'RTE56765', 'Class A', '2024-11-30', 'Truck', '3', ARRAY['CDL'], 1, NULL, 4.0, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- -- Insert a vehicle for driver 1
-- INSERT INTO vehicles (plate_number, type, make, model, year, capacity, fuel_type, status_id, current_location_latitude, current_location_longitude)
-- VALUES ('XYZ123', 'Truck', 'Ford', 'F-150', 2020, 1000, 'Gasoline', 1, 47.6062, -122.3321);

-- -- Insert an order for driver 1
-- INSERT INTO orders (user_id, from_location, to_location, status, created_at, updated_at)
-- VALUES (1, 'Seattle', 'Portland', 'Pending', NOW(), NOW());

-- Insert sample warehouse data
INSERT INTO warehouses (name, street1, city, state, zip, country, latitude, longitude, capacity)
VALUES
    ('Seattle Hub', '789 Port Way', 'Seattle', 'WA', '98101', 'USA', 47.6062, -122.3321, 10000),
    ('Vancouver DC', '456 Harbor St', 'Vancouver', 'BC', 'V6B 1A1', 'Canada', 49.2827, -123.1207, 15000),
    ('London Center', '123 Dock Lane', 'London', NULL, 'E1 6AN', 'UK', 51.5074, -0.1278, 12000);

-- Insert sample storage locations
INSERT INTO storage_locations (warehouse_id, section, aisle, shelf, capacity, is_occupied) VALUES
    (1, 'A', '1', '1A', 500.00, false),
    (1, 'A', '1', '1B', 500.00, true),
    (1, 'B', '2', '2A', 750.00, false),
    (2, 'A', '1', '1A', 1000.00, true),
    (2, 'B', '1', '1B', 1000.00, false),
    (3, 'A', '1', '1A', 1500.00, true);

-- Insert sample orders
INSERT INTO orders (user_id, from_location, to_location, weight, status, created_at, updated_at)
VALUES
    (1, 'Seattle', 'Portland', 1500.00, 'Pending', '2024-12-10 16:24:52+03', '2024-12-10 16:24:52+03'),
    (2, 'Vancouver', 'San Francisco', 2000.00, 'In Progress', '2024-12-10 16:24:52+03', '2024-12-10 16:24:52+03'),
    (3, 'London', 'Birmingham', 1200.50, 'Completed', '2024-12-10 16:24:52+03', '2024-12-10 16:24:52+03'),
    (4, 'Berlin', 'Munich', 1800.75, 'Pending', '2024-12-10 16:24:52+03', '2024-12-10 16:24:52+03'),
    (5, 'Toronto', 'Montreal', 1300.25, 'In Progress', '2024-12-10 16:24:52+03', '2024-12-10 16:24:52+03');

-- -- Insert sample deliveries
INSERT INTO deliveries (tracking_number, user_id, vehicle_id, status_id, pickup_time, delivery_time, estimated_delivery_time, actual_delivery_time, cargo_type, cargo_weight, special_instructions, customer_id, customer_signature, pickup_warehouse_id, delivery_warehouse_id, pickup_latitude, pickup_longitude, delivery_latitude, delivery_longitude, payment_status, estimated_fuel_consumption, actual_fuel_consumption, route_efficiency_score, weather_conditions, proof_of_delivery_image_url, package_condition_images, from_location, completed_at, created_at, updated_at)
VALUES
    ('DEL001', 1, 3, 2, '2024-12-10 14:00:00+03', '2024-12-10 16:00:00+03', '2024-12-10 15:30:00+03', NULL, 'Electronics', 500.00, 'Handle with care', 1, 'John Doe', 1, 2, 47.6062, -122.3321, 45.5155, -122.6789, 'paid', 5.00, 4.50, 0.95, 'Clear', 'https://img.freepik.com/free-photo/close-up-delivery-person-offering-parcel-client_23-2149095936.jpg', '{"condition": "good", "images": ["https://img.freepik.com/free-photo/close-up-delivery-person-offering-parcel-client_23-2149095936.jpg"]}', 'Seattle', NULL, '2024-12-10 16:33:49+03', '2024-12-10 16:33:49+03'),
    ('DEL002', 2, 2, 1, '2024-12-10 14:30:00+03', '2024-12-10 16:30:00+03', '2024-12-10 16:00:00+03', NULL, 'Furniture', 1500.50, 'Fragile', 2, 'Jane Smith', 1, 3, 49.2827, -123.1207, 51.5074, -0.1278, 'pending', 10.00, 9.00, 0.90, 'Rainy', 'https://img.freepik.com/free-photo/close-up-delivery-person-with-parcels_23-2149095943.jpg', '{"condition": "damaged", "images": ["https://img.freepik.com/free-photo/close-up-delivery-person-with-parcels_23-2149095943.jpg"]}', 'Vancouver', NULL, '2024-12-10 16:33:49+03', '2024-12-10 16:33:49+03'),
    ('DEL003', 3, 1, 3, '2024-12-10 15:00:00+03', '2024-12-10 17:00:00+03', '2024-12-10 16:30:00+03', NULL, 'Groceries', 300.25, NULL, 3, 'Alice Johnson', 2, 1, 51.5074, -0.1278, 37.8044, -122.2711, 'paid', 3.00, 2.50, 0.85, 'Sunny', 'https://img.freepik.com/free-photo/portrait-happy-black-courier-delivering-packages-looking-camera_637285-2084.jpg', '{"condition": "good", "images": ["https://img.freepik.com/free-photo/portrait-happy-black-courier-delivering-packages-looking-camera_637285-2084.jpg"]}', 'London', NULL, '2024-12-10 16:33:49+03', '2024-12-10 16:33:49+03');

-- Insert sample driver documents using user_id
INSERT INTO driver_documents (
    user_id,
    document_type,
    document_number,
    issue_date,
    expiry_date,
    document_url,
    verification_status
) VALUES 
    -- For driver1 (user_id = 3)
    (3, 'CDL_LICENSE', 'DL123456', '2023-01-01', '2025-12-31', 
     'https://fleetflow-assets.s3.amazonaws.com/documents/drivers/3/cdl_license.pdf', 
     'verified'),
    (3, 'MEDICAL_CERTIFICATE', 'MC789012', '2023-06-01', '2024-06-01',
     'https://fleetflow-assets.s3.amazonaws.com/documents/drivers/3/medical_cert.pdf', 
     'verified'),
    
    -- For driver2 (user_id = 7)
    (7, 'CDL_LICENSE', 'DL789012', '2023-03-15', '2024-11-30',
     'https://fleetflow-assets.s3.amazonaws.com/documents/drivers/7/cdl_license.pdf', 
     'verified'),
    (7, 'MEDICAL_CERTIFICATE', 'MC345678', '2023-09-01', '2024-09-01',
     'https://fleetflow-assets.s3.amazonaws.com/documents/drivers/7/medical_cert.pdf', 
     'pending');

-- Insert sample addresses for each user
INSERT INTO addresses (
    user_id,
    street1,
    street2,
    city,
    state,
    zip,
    country,
    address_type,
    is_default,
    latitude,
    longitude
) VALUES
    -- Customer1 (user_id = 1) addresses
    (1, '123 Pine Street', 'Apt 4B', 'Seattle', 'WA', '98101', 'USA', 'residential', true, 47.6062, -122.3321),
    (1, '456 Oak Avenue', 'Suite 200', 'Seattle', 'WA', '98102', 'USA', 'business', false, 47.6152, -122.3301),

    -- Fleet Manager1 (user_id = 2) addresses
    (2, '789 Maple Drive', NULL, 'Vancouver', 'BC', 'V6B 1A1', 'Canada', 'residential', true, 49.2827, -123.1207),
    (2, '321 Fleet Street', 'Floor 3', 'Vancouver', 'BC', 'V6B 2B2', 'Canada', 'business', false, 49.2762, -123.1187),

    -- Driver1 (user_id = 3) addresses
    (3, '567 Thames Road', 'Unit 12', 'London', NULL, 'E1 6AN', 'UK', 'residential', true, 51.5074, -0.1278),
    (3, '890 Victoria Street', NULL, 'London', NULL, 'E1 7BB', 'UK', 'business', false, 51.5124, -0.1258),

    -- Loader1 (user_id = 4) addresses
    (4, '234 Harbor Road', NULL, 'Sydney', 'NSW', '2000', 'Australia', 'residential', true, -33.8688, 151.2093),
    (4, '567 Dock Street', 'Unit 5', 'Sydney', 'NSW', '2001', 'Australia', 'business', false, -33.8712, 151.2033),

    -- Admin1 (user_id = 5) addresses
    (5, '789 Berlin Way', 'Apt 15', 'Berlin', NULL, '10115', 'Germany', 'residential', true, 52.5200, 13.4050),
    (5, '432 Admin Strasse', NULL, 'Berlin', NULL, '10117', 'Germany', 'business', false, 52.5166, 13.3890),

    -- Manager1 (user_id = 6) addresses
    (6, '123 Oslo Street', NULL, 'Oslo', NULL, '0155', 'Norway', 'residential', true, 59.9139, 10.7522),
    (6, '456 Manager Road', 'Floor 2', 'Oslo', NULL, '0157', 'Norway', 'business', false, 59.9127, 10.7461),

    -- Driver2 (user_id = 7) addresses
    (7, '901 Leeds Road', 'Flat 3C', 'London', NULL, 'E2 8AN', 'UK', 'residential', true, 51.5225, -0.1389),
    (7, '432 Driver Lane', NULL, 'London', NULL, 'E2 9BB', 'UK', 'business', false, 51.5214, -0.1352);

-- Insert sample inventory items
INSERT INTO inventory_items (
    warehouse_id,
    name,
    sku,
    quantity,
    min_quantity,
    max_quantity,
    unit_price,
    last_restocked_at
) VALUES
    -- Warehouse 1 (Seattle Hub) inventory
    (1, 'Moving Boxes Large', 'BOX-L-001', 500, 100, 1000, 2.99, '2024-01-15 10:00:00+00'),
    (1, 'Packing Tape', 'TAPE-001', 200, 50, 400, 3.50, '2024-01-15 10:00:00+00'),
    (1, 'Bubble Wrap Roll', 'BUBBLE-001', 150, 30, 300, 15.99, '2024-01-15 10:00:00+00'),
    (1, 'Furniture Blankets', 'BLANKET-001', 100, 20, 200, 12.99, '2024-01-15 10:00:00+00'),

    -- Warehouse 2 (Vancouver DC) inventory
    (2, 'Moving Boxes Medium', 'BOX-M-001', 750, 150, 1500, 1.99, '2024-01-16 11:00:00+00'),
    (2, 'Strapping Tape', 'TAPE-002', 300, 75, 600, 4.50, '2024-01-16 11:00:00+00'),
    (2, 'Packing Peanuts', 'PEANUTS-001', 1000, 200, 2000, 8.99, '2024-01-16 11:00:00+00'),
    (2, 'Hand Truck', 'TRUCK-001', 20, 5, 40, 89.99, '2024-01-16 11:00:00+00'),

    -- Warehouse 3 (London Center) inventory
    (3, 'Moving Boxes Small', 'BOX-S-001', 1000, 200, 2000, 0.99, '2024-01-17 12:00:00+00'),
    (3, 'Packaging Labels', 'LABEL-001', 5000, 1000, 10000, 0.05, '2024-01-17 12:00:00+00'),
    (3, 'Stretch Wrap', 'WRAP-001', 200, 40, 400, 12.99, '2024-01-17 12:00:00+00'),
    (3, 'Box Cutter', 'CUTTER-001', 50, 10, 100, 4.99, '2024-01-17 12:00:00+00');

-- Insert sample delivery feedback (only for existing deliveries 1, 2, and 3)
INSERT INTO delivery_feedback (
    delivery_id,
    customer_id,
    rating,
    feedback_text,
    timeliness_rating,
    driver_rating,
    package_condition_rating
) VALUES
    -- Feedback for DEL001
    (1, 1, 4.5, 
     'Excellent service, very professional driver. Package arrived in perfect condition.', 
     4.8, 4.7, 5.0),

    -- Feedback for DEL002
    (2, 2, 3.8, 
     'Good delivery, but slightly delayed. Driver was courteous and kept me informed.', 
     3.5, 4.5, 4.0),

    -- Feedback for DEL003
    (3, 3, 5.0, 
     'Perfect delivery timing and handling. Driver went above and beyond expectations.', 
     5.0, 5.0, 4.9);

-- Insert sample driver schedules
INSERT INTO driver_schedules (
    user_id,
    schedule_date,
    start_time,
    end_time,
    status,
    break_duration,
    overtime_duration,
    notes
) VALUES
    -- Driver 1 schedules
    (1, '2024-03-01', '08:00', '16:00', 'scheduled', 60, 0, 'Regular shift - Seattle route'),
    (1, '2024-03-02', '09:00', '17:00', 'completed', 45, 30, 'Covered extra delivery in downtown'),
    (1, '2024-03-03', '07:00', '15:00', 'scheduled', 60, 0, 'Early morning deliveries'),

    -- Driver 2 schedules
    (2, '2024-03-01', '10:00', '18:00', 'completed', 60, 0, 'Vancouver downtown route'),
    (2, '2024-03-02', '08:00', '16:00', 'scheduled', 45, 0, 'Standard shift - suburban route'),
    (2, '2024-03-03', '12:00', '20:00', 'scheduled', 60, 0, 'Evening delivery shift');

-- Insert sample driver performance metrics
INSERT INTO driver_performance_metrics (
    user_id,
    metric_date,
    deliveries_completed,
    on_time_delivery_rate,
    customer_rating_avg,
    fuel_efficiency,
    safety_score,
    violations_count,
    total_distance_covered
) VALUES
    -- Driver 1 performance metrics
    (1, '2024-03-01', 10, 95.0, 4.8, 15.0, 98.0, 0, 120.5),
    (1, '2024-03-02', 8, 90.0, 4.7, 14.5, 97.5, 1, 110.3),
    (1, '2024-03-03', 12, 100.0, 4.9, 15.2, 99.0, 0, 130.7),

    -- Driver 2 performance metrics
    (2, '2024-03-01', 9, 92.0, 4.6, 14.8, 96.0, 0, 115.4),
    (2, '2024-03-02', 7, 88.0, 4.5, 14.2, 95.5, 1, 105.2),
    (2, '2024-03-03', 11, 98.0, 4.8, 15.1, 98.5, 0, 125.6),

    -- Driver 3 performance metrics
    (3, '2024-03-01', 8, 90.0, 4.7, 14.7, 97.0, 0, 112.3),
    (3, '2024-03-02', 6, 85.0, 4.4, 13.9, 94.5, 2, 100.1),
    (3, '2024-03-03', 10, 95.0, 4.9, 15.0, 98.0, 0, 120.9);

-- Insert sample fleet analytics data
INSERT INTO fleet_analytics (
    vehicle_id,
    analysis_date,
    total_distance,
    fuel_consumption,
    fuel_cost,
    maintenance_cost,
    idle_time,
    carbon_emissions
) VALUES
    -- Vehicle 1 analytics
    (1, '2024-03-01', 150.0, 20.0, 60.0, 100.0, 30, 50.0),
    (1, '2024-03-02', 140.0, 18.5, 55.5, 80.0, 25, 48.0),
    (1, '2024-03-03', 160.0, 21.0, 63.0, 90.0, 35, 52.0),

    -- Vehicle 2 analytics
    (2, '2024-03-01', 130.0, 19.0, 57.0, 110.0, 20, 49.0),
    (2, '2024-03-02', 120.0, 17.5, 52.5, 85.0, 15, 46.0),
    (2, '2024-03-03', 140.0, 20.5, 61.5, 95.0, 25, 51.0),

    -- Vehicle 3 analytics
    (3, '2024-03-01', 160.0, 22.0, 66.0, 120.0, 40, 54.0),
    (3, '2024-03-02', 150.0, 20.0, 60.0, 100.0, 30, 50.0),
    (3, '2024-03-03', 170.0, 23.0, 69.0, 110.0, 45, 56.0);

-- Insert sample invoices
INSERT INTO invoices (
    delivery_id,
    customer_id,
    amount,
    status,
    due_date,
    payment_date,
    payment_method,
    created_at,
    updated_at
) VALUES
    -- Invoice for delivery DEL001
    (1, 1, 150.00, 'paid', '2024-03-10', '2024-03-05', 'credit_card', '2024-03-01 10:00:00+00', '2024-03-05 10:00:00+00'),
    
    -- Invoice for delivery DEL002
    (2, 2, 200.00, 'pending', '2024-03-15', NULL, 'bank_transfer', '2024-03-01 11:00:00+00', '2024-03-01 11:00:00+00'),
    
    -- Invoice for delivery DEL003
    (3, 3, 120.50, 'paid', '2024-03-12', '2024-03-07', 'paypal', '2024-03-01 12:00:00+00', '2024-03-07 12:00:00+00');

-- Insert sample loading schedules
INSERT INTO loading_schedules (
    warehouse_id,
    vehicle_id,
    scheduled_time,
    estimated_duration,
    actual_duration,
    status,
    dock_number,
    loader_id,
    notes
) VALUES
    -- Loading schedule for warehouse 1, vehicle 1
    (1, 1, '2024-03-01 08:00:00+00', 120, 115, 'completed', 'Dock 1', 4, 'Loaded electronics for delivery to Portland'),
    (1, 2, '2024-03-01 10:00:00+00', 90, 95, 'completed', 'Dock 2', 4, 'Loaded furniture for delivery to San Francisco'),
    (1, 3, '2024-03-01 12:00:00+00', 60, 60, 'completed', 'Dock 3', 4, 'Loaded groceries for delivery to Birmingham'),

    -- Loading schedule for warehouse 2, vehicle 2
    (2, 2, '2024-03-02 09:00:00+00', 120, 110, 'completed', 'Dock 1', 4, 'Loaded electronics for delivery to Vancouver'),
    (2, 1, '2024-03-02 11:00:00+00', 90, 85, 'completed', 'Dock 2', 4, 'Loaded furniture for delivery to London'),
    (2, 3, '2024-03-02 13:00:00+00', 60, 65, 'completed', 'Dock 3', 4, 'Loaded groceries for delivery to Munich'),

    -- Loading schedule for warehouse 3, vehicle 3
    (3, 3, '2024-03-03 08:00:00+00', 120, 115, 'completed', 'Dock 1', 4, 'Loaded electronics for delivery to Toronto'),
    (3, 1, '2024-03-03 10:00:00+00', 90, 95, 'completed', 'Dock 2', 4, 'Loaded furniture for delivery to Montreal'),
    (3, 2, '2024-03-03 12:00:00+00', 60, 60, 'completed', 'Dock 3', 4, 'Loaded groceries for delivery to Seattle');

-- Insert sample location history data
INSERT INTO location_history (
    vehicle_id,
    latitude,
    longitude,
    altitude,
    speed,
    heading,
    accuracy,
    recorded_at,
    created_at
) VALUES
    -- Location history for vehicle 1
    (1, 47.6062, -122.3321, 30.0, 60.0, 90.0, 5.0, '2024-03-01 08:00:00+00', '2024-03-01 08:00:00+00'),
    (1, 47.6097, -122.3331, 32.0, 62.0, 92.0, 5.0, '2024-03-01 08:05:00+00', '2024-03-01 08:05:00+00'),
    (1, 47.6123, -122.3342, 35.0, 65.0, 95.0, 5.0, '2024-03-01 08:10:00+00', '2024-03-01 08:10:00+00'),

    -- Location history for vehicle 2
    (2, 49.2827, -123.1207, 40.0, 55.0, 85.0, 5.0, '2024-03-02 09:00:00+00', '2024-03-02 09:00:00+00'),
    (2, 49.2850, -123.1215, 42.0, 57.0, 87.0, 5.0, '2024-03-02 09:05:00+00', '2024-03-02 09:05:00+00'),
    (2, 49.2873, -123.1223, 45.0, 60.0, 90.0, 5.0, '2024-03-02 09:10:00+00', '2024-03-02 09:10:00+00'),

    -- Location history for vehicle 3
    (3, 51.5074, -0.1278, 50.0, 50.0, 80.0, 5.0, '2024-03-03 10:00:00+00', '2024-03-03 10:00:00+00'),
    (3, 51.5090, -0.1285, 52.0, 52.0, 82.0, 5.0, '2024-03-03 10:05:00+00', '2024-03-03 10:05:00+00'),
    (3, 51.5106, -0.1292, 55.0, 55.0, 85.0, 5.0, '2024-03-03 10:10:00+00', '2024-03-03 10:10:00+00');

-- Insert sample maintenance records
INSERT INTO maintenance_records (
    vehicle_id,
    type,
    description,
    service_date,
    cost,
    odometer_reading,
    performed_by,
    next_service_date,
    notes
) VALUES
    -- Maintenance records for vehicle 1
    (1, 'Oil Change', 'Changed engine oil and filter', '2024-01-15', 75.00, 15000, 'Quick Lube', '2024-07-15', 'Recommended next oil change in 6 months'),
    (1, 'Brake Inspection', 'Inspected and replaced brake pads', '2024-02-20', 150.00, 16000, 'Brake Masters', '2024-08-20', 'Brake pads replaced, rotors in good condition'),
    (1, 'Tire Rotation', 'Rotated all four tires', '2024-03-10', 50.00, 17000, 'Tire Shop', '2024-09-10', 'Tires rotated, tread depth checked'),

    -- Maintenance records for vehicle 2
    (2, 'Oil Change', 'Changed engine oil and filter', '2024-01-18', 80.00, 14000, 'Quick Lube', '2024-07-18', 'Recommended next oil change in 6 months'),
    (2, 'Battery Replacement', 'Replaced vehicle battery', '2024-02-25', 120.00, 15000, 'Auto Electric', '2025-02-25', 'New battery installed, electrical system checked'),
    (2, 'Transmission Service', 'Serviced transmission fluid and filter', '2024-03-15', 200.00, 16000, 'Transmission Experts', '2025-03-15', 'Transmission fluid and filter replaced'),

    -- Maintenance records for vehicle 3
    (3, 'Oil Change', 'Changed engine oil and filter', '2024-01-20', 70.00, 13000, 'Quick Lube', '2024-07-20', 'Recommended next oil change in 6 months'),
    (3, 'Coolant Flush', 'Flushed and replaced coolant', '2024-02-28', 100.00, 14000, 'Radiator Shop', '2025-02-28', 'Coolant system flushed, no leaks detected'),
    (3, 'Brake Fluid Change', 'Changed brake fluid', '2024-03-20', 60.00, 15000, 'Brake Masters', '2025-03-20', 'Brake fluid replaced, brakes bled');

-- Insert sample notifications
INSERT INTO notifications (
    user_id,
    type,
    title,
    message,
    read,
    related_entity_type,
    related_entity_id,
    created_at,
    updated_at
) VALUES
    (1, 'delivery', 'Delivery Completed', 'Your delivery DEL001 has been completed successfully.', false, 'delivery', 1, '2024-03-01 10:00:00+00', '2024-03-01 10:00:00+00'),
    (2, 'maintenance', 'Maintenance Scheduled', 'Your vehicle has a scheduled maintenance on 2024-03-15.', false, 'vehicle', 2, '2024-03-01 11:00:00+00', '2024-03-01 11:00:00+00'),
    (3, 'order', 'Order Dispatched', 'Your order ORD003 has been dispatched.', true, 'order', 3, '2024-03-01 12:00:00+00', '2024-03-01 12:00:00+00');

-- Insert sample packages
INSERT INTO packages (
    delivery_id,
    tracking_number,
    weight,
    dimensions,
    storage_location_id,
    status,
    created_at,
    updated_at
) VALUES
    (1, 'PKG001', 10.5, '{"length": 20, "width": 15, "height": 10}', 1, 'in_transit', '2024-03-01 10:00:00+00', '2024-03-01 10:00:00+00'),
    (2, 'PKG002', 5.0, '{"length": 15, "width": 10, "height": 8}', 2, 'delivered', '2024-03-01 11:00:00+00', '2024-03-01 11:00:00+00'),
    (3, 'PKG003', 7.5, '{"length": 18, "width": 12, "height": 9}', 3, 'pending', '2024-03-01 12:00:00+00', '2024-03-01 12:00:00+00');

-- Insert sample warehouse equipment
INSERT INTO warehouse_equipment (
    warehouse_id,
    equipment_type,
    equipment_id,
    status,
    last_maintenance,
    next_maintenance,
    condition_rating,
    notes
) VALUES
    (1, 'Forklift', 'EQ001', 'operational', '2024-01-15', '2024-07-15', 4.5, 'Routine maintenance completed'),
    (2, 'Pallet Jack', 'EQ002', 'operational', '2024-02-20', '2024-08-20', 4.0, 'Hydraulic system checked'),
    (3, 'Conveyor Belt', 'EQ003', 'under_maintenance', '2024-03-10', '2024-09-10', 3.5, 'Belt replacement scheduled');

-- Insert sample warehouse inventory
INSERT INTO warehouse_inventory (
    warehouse_id,
    item_name,
    item_category,
    quantity,
    unit,
    minimum_threshold,
    maximum_capacity,
    last_restocked,
    next_restock_date,
    created_at,
    updated_at
) VALUES
    (1, 'Moving Boxes Large', 'Packaging', 500, 'units', 100, 1000, '2024-01-15', '2024-07-15', '2024-01-15 10:00:00+00', '2024-01-15 10:00:00+00'),
    (2, 'Packing Tape', 'Packaging', 200, 'rolls', 50, 400, '2024-02-20', '2024-08-20', '2024-02-20 11:00:00+00', '2024-02-20 11:00:00+00'),
    (3, 'Bubble Wrap Roll', 'Packaging', 150, 'rolls', 30, 300, '2024-03-10', '2024-09-10', '2024-03-10 12:00:00+00', '2024-03-10 12:00:00+00');
