INSERT INTO keycloak_role (id, name, realm_id, client_role)
SELECT gen_random_uuid(), 'ROLE_ADMIN', r.id, false
FROM realm r
WHERE r.name = 'master'
  AND NOT EXISTS (
    SELECT 1 FROM keycloak_role WHERE name = 'ROLE_ADMIN'
);

INSERT INTO keycloak_role (id, name, realm_id, client_role)
SELECT gen_random_uuid(), 'ROLE_TEACHER', r.id, false
FROM realm r
WHERE r.name = 'master'
  AND NOT EXISTS (
    SELECT 1 FROM keycloak_role WHERE name = 'ROLE_TEACHER'
);

INSERT INTO keycloak_role (id, name, realm_id, client_role)
SELECT gen_random_uuid(), 'ROLE_USER', r.id, false
FROM realm r
WHERE r.name = 'master'
  AND NOT EXISTS (
    SELECT 1 FROM keycloak_role WHERE name = 'ROLE_USER'
);