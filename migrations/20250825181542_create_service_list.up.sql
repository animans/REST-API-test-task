CREATE TABLE service_list (
	service_id SERIAL PRIMARY KEY,
	service_price INTEGER NOT NULL,
	service_name VARCHAR(128) NOT NULL,
	service_uuid VARCHAR(36) NOT NULL,
	service_created_at DATE NOT NULL
);