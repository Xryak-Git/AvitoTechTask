-- Вставляем пользователей, если их нет
INSERT INTO employee (username, first_name, last_name)
VALUES
    ('test', 'First', 'Last'),
    ('anotherone', 'First', 'Last')
ON CONFLICT
    DO NOTHING;

-- Вставляем организации, если их нет
INSERT INTO organization (name, description, type)
VALUES
    ('Org1', 'Desc', 'IE'::organization_type),
    ('Org2', 'Desc', 'LLC'::organization_type),
    ('Org3', 'Desc', 'JSC'::organization_type),
    ('Org4', 'Desc', 'IE'::organization_type)
ON CONFLICT
    DO NOTHING;

-- Теперь создаем соответствие пользователей и организаций
WITH users AS (
    SELECT id, username
    FROM employee
    WHERE username IN ('test', 'anotherone')
),
     orgs AS (
         SELECT id, name
         FROM organization
         WHERE name IN ('Org1', 'Org2', 'Org3', 'Org4')
     )
INSERT INTO organization_responsible (organization_id, user_id)
VALUES
    ((SELECT id FROM orgs WHERE name = 'Org1'), (SELECT id FROM users WHERE username = 'test')),
    ((SELECT id FROM orgs WHERE name = 'Org2'), (SELECT id FROM users WHERE username = 'test')),
    ((SELECT id FROM orgs WHERE name = 'Org3'), (SELECT id FROM users WHERE username = 'anotherone')),
    ((SELECT id FROM orgs WHERE name = 'Org4'), (SELECT id FROM users WHERE username = 'anotherone'))
ON CONFLICT
    DO NOTHING;
