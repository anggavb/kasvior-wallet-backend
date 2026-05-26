TRUNCATE TABLE payment_methods RESTART IDENTITY CASCADE;

INSERT INTO payment_methods
(name, logo, method, tax)
VALUES
('BRI', 'logo/bri.png', 'bank', 5000),
('DANA', 'logo/dana.png', 'online', 5000),
('BCA', 'logo/bca.png', 'bank', 5000),
('GOPAY', 'logo/gopay.png', 'online', 5000),
('OVO', 'logo/ovo.png', 'online', 5000);