CREATE database assignment1_db;
USE assignment1_db;

CREATE TABLE Membership (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tier ENUM('Basic', 'Premium', 'VIP') NOT NULL,
    hourly_rate DECIMAL(5, 2),
    priority_access BOOLEAN,
    max_booking_limit INT
);

CREATE TABLE Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    membership_id INT DEFAULT 1, -- Defaults to 'Basic'
    registration_status ENUM('Pending', 'Verified') DEFAULT 'Pending',
    FOREIGN KEY (membership_id) REFERENCES Membership(id) ON DELETE SET NULL
);

CREATE TABLE Vehicles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    model VARCHAR(255) NOT NULL,
    license_plate VARCHAR(255) UNIQUE NOT NULL,
    location VARCHAR(255) NOT NULL,
    charge_level INT NOT NULL,
    cleanliness_status ENUM('Clean', 'Moderate', 'Dirty')
);

CREATE TABLE Reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status ENUM('Pending', 'Confirmed', 'Cancelled', 'Completed', 'Paid') DEFAULT 'Pending',
    estimated_cost DECIMAL(10, 2) DEFAULT 0.00,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (vehicle_id) REFERENCES Vehicles(id) ON DELETE CASCADE
);

CREATE TABLE Billing (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT NOT NULL,
    user_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_status ENUM('Pending', 'Paid', 'Refunded') DEFAULT 'Paid',
    FOREIGN KEY (reservation_id) REFERENCES Reservations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE TABLE Promotions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    discount_percentage DECIMAL(5, 2) NOT NULL,
    valid_from DATE NOT NULL,
    valid_until DATE NOT NULL,
);

INSERT INTO Membership (tier, hourly_rate, priority_access, max_booking_limit)
VALUES 
('Basic', 30.00, FALSE, 2),
('Premium', 20.00, TRUE, 4),
('VIP', 10.00, TRUE, 6);

INSERT INTO Vehicles (model, license_plate, location, charge_level, cleanliness_status)
VALUES 
('BMW iX', 'SGT1234X', 'Jurong East Interchange', 90, 'Clean'),
('Tesla Model S', 'SKP5009Z', 'Bukit Batok MRT Station', 70, 'Clean'),
('Hyundai Ioniq 5', 'SNK2022G', 'Orchard Road Central', 85, 'Moderate'),
('Audi Q4 E-tron', 'SKS3609Z', 'Harbourfront', 45, 'Dirty');

INSERT INTO Promotions (code, discount_percentage, valid_from, valid_until)
VALUES
('MerryChristmas2024', 10, '2024-12-01', '2024-12-31'),
('CODE20', 20, '2024-12-01', '2024-12-05'); 

select * from Users;
select * from Reservations where user_id = 2;
select * from Billing where user_id = 2;



