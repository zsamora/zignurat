INSERT INTO users (names, surnames, date_birth, created_at, updated_at)
VALUES ('Padre Pedro', 'Gutiérrez Rojas', '1978-03-14', now(), now());

INSERT INTO organizations (org_type, name, diocese, commune, address, admin_id, created_at, updated_at)
SELECT 1, 'Santo Espíritu Sagrado', 'Santiago', 'Puente Alto', 'Calle Larga 123', id, now(), now()
FROM users
WHERE names = 'Padre Pedro' AND surnames = 'Gutiérrez Rojas';
