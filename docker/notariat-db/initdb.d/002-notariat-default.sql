INSERT INTO organizations (org_type, name, diocese, commune, address, created_at, updated_at)
VALUES (1, 'Santo Espíritu Sagrado', 'Santiago', 'Puente Alto', 'Calle Larga 123', now(), now());

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Padre Pedro', 'Gutiérrez Rojas', 0, id, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Ana Belén', 'Rodríguez Fuentes', 1, id, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';
