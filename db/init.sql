
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(120) NOT NULL,
    owner_team VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO services (service_name, owner_team) VALUES
('Auth-Service-Container', 'Mahasiswa: Gede Troy Wiswama Bhareswara (42430052)'),
('Inventory-API-Gateway', 'Mahasiswa: Gede Troy Wiswama Bhareswara (42430052)');


