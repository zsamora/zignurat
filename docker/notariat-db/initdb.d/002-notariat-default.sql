INSERT INTO organizations (org_type, name, diocese, commune, address, created_at, updated_at)
VALUES (1, 'Santo Espíritu Sagrado', 'Santiago', 'Puente Alto', 'Calle Larga 123', now(), now());

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Soporte IT', '', 0, id, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Padre Pedro', 'Gutiérrez Rojas', 1, id, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Ana Belén', 'Rodríguez Fuentes', 2, id, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';

INSERT INTO organizations (org_type, name, diocese, commune, address, created_at, updated_at)
VALUES (2, 'Capilla San José', 'Santiago', 'Providencia', 'Av. Providencia 1000', now(), now());

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Hermana Rosa', 'Martínez Silva', 1, id, now(), now()
FROM organizations WHERE name = 'Capilla San José';

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Hermano Francisco', 'Díaz Contreras', 0, id, now(), now()
FROM organizations WHERE name = 'Capilla San José';

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Carla Jiménez', 'Paredes Rojas', 2, id, now(), now()
FROM organizations WHERE name = 'Capilla San José';
