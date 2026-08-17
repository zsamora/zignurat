INSERT INTO organizations (org_type, name, diocese, commune, address, uuid, created_at, updated_at)
VALUES (1, 'Santo Espíritu Sagrado', 'Santiago', 'Puente Alto', 'Calle Larga 123',
        '67f484f7-ffe9-46ee-b6e6-0b47066afc5a', now(), now());

INSERT INTO members (names, surnames, role, org_id, created_at, updated_at)
SELECT 'Soporte IT', '', 0, id, now(), now()
FROM organizations WHERE name = 'Santo Espíritu Sagrado';

INSERT INTO members (names, surnames, role, org_id, uuid, created_at, updated_at)
SELECT 'Padre Pedro', 'Gutiérrez Rojas', 1, id, 'f4e3b845-942b-4082-ab9e-aa045f15ee26', now(), now()
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

INSERT INTO users (names, surnames, date_birth, uuid, created_at, updated_at)
VALUES ('María Elena', 'Vásquez Bravo', '1990-04-12', '830169b7-e3bd-44e8-8f3a-ce1e5b2a5ff8', now(), now());
