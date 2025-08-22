-- WohnFair Database Seed Script (idempotent)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tables
CREATE TABLE IF NOT EXISTS user_groups (
  id UUID PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  weight NUMERIC NOT NULL DEFAULT 1.0,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cities (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  country TEXT NOT NULL,
  region TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS property_types (
  id UUID PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE urgency_level AS ENUM ('LOW','MEDIUM','HIGH');
CREATE TYPE request_status AS ENUM ('PENDING','PROCESSING','ALLOCATED','CANCELLED');

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'urgency_level') THEN
    CREATE TYPE urgency_level AS ENUM ('LOW','MEDIUM','HIGH');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'request_status') THEN
    CREATE TYPE request_status AS ENUM ('PENDING','PROCESSING','ALLOCATED','CANCELLED');
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS housing_requests (
  id UUID PRIMARY KEY,
  user_id TEXT NOT NULL,
  city_id UUID REFERENCES cities(id),
  property_type_id UUID REFERENCES property_types(id),
  urgency_level urgency_level NOT NULL,
  status request_status NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS queue_entries (
  id UUID PRIMARY KEY,
  request_id UUID UNIQUE REFERENCES housing_requests(id) ON DELETE CASCADE,
  position INT NOT NULL,
  priority_score NUMERIC NOT NULL,
  estimated_wait_time INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fairness_metrics (
  id UUID PRIMARY KEY,
  metric_name TEXT NOT NULL,
  metric_value NUMERIC NOT NULL,
  group_id UUID REFERENCES user_groups(id),
  city_id UUID REFERENCES cities(id),
  calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Inserts (upserts)
INSERT INTO user_groups (id, name, weight, description)
VALUES
  ('550e8400-e29b-41d4-a716-446655440001','Refugee',1.5,'Refugees and asylum seekers'),
  ('550e8400-e29b-41d4-a716-446655440002','Disabled',1.3,'Persons with disabilities'),
  ('550e8400-e29b-41d4-a716-446655440003','Senior',1.2,'Elderly persons (65+)'),
  ('550e8400-e29b-41d4-a716-446655440004','Low-Income',1.1,'Low-income households'),
  ('550e8400-e29b-41d4-a716-446655440005','General',1.0,'General population')
ON CONFLICT (id) DO NOTHING;

INSERT INTO cities (id, name, country, region) VALUES
  ('550e8400-e29b-41d4-a716-446655440010','Berlin','Germany','Berlin'),
  ('550e8400-e29b-41d4-a716-446655440011','Munich','Germany','Bavaria'),
  ('550e8400-e29b-41d4-a716-446655440012','Hamburg','Germany','Hamburg'),
  ('550e8400-e29b-41d4-a716-446655440013','Cologne','Germany','North Rhine-Westphalia'),
  ('550e8400-e29b-41d4-a716-446655440014','Frankfurt','Germany','Hesse')
ON CONFLICT (id) DO NOTHING;

INSERT INTO property_types (id, name, description) VALUES
  ('550e8400-e29b-41d4-a716-446655440020','Apartment','Multi-family residential unit'),
  ('550e8400-e29b-41d4-a716-446655440021','Studio','Single-room apartment'),
  ('550e8400-e29b-41d4-a716-446655440022','House','Single-family detached house'),
  ('550e8400-e29b-41d4-a716-446655440023','Shared Room','Room in shared accommodation'),
  ('550e8400-e29b-41d4-a716-446655440024','Wheelchair Accessible','Accessible housing unit')
ON CONFLICT (id) DO NOTHING;

INSERT INTO housing_requests (id, user_id, city_id, property_type_id, urgency_level, status, created_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440030','user-001','550e8400-e29b-41d4-a716-446655440010','550e8400-e29b-41d4-a716-446655440020','HIGH','PENDING', NOW() - INTERVAL '2 hours'),
  ('550e8400-e29b-41d4-a716-446655440031','user-002','550e8400-e29b-41d4-a716-446655440011','550e8400-e29b-41d4-a716-446655440021','MEDIUM','PROCESSING', NOW() - INTERVAL '3 hours'),
  ('550e8400-e29b-41d4-a716-446655440032','user-003','550e8400-e29b-41d4-a716-446655440012','550e8400-e29b-41d4-a716-446655440022','LOW','PENDING', NOW() - INTERVAL '4 hours'),
  ('550e8400-e29b-41d4-a716-446655440033','user-004','550e8400-e29b-41d4-a716-446655440013','550e8400-e29b-41d4-a716-446655440023','HIGH','PENDING', NOW() - INTERVAL '5 hours'),
  ('550e8400-e29b-41d4-a716-446655440034','user-005','550e8400-e29b-41d4-a716-446655440014','550e8400-e29b-41d4-a716-446655440024','MEDIUM','PENDING', NOW() - INTERVAL '6 hours')
ON CONFLICT (id) DO NOTHING;

INSERT INTO queue_entries (id, request_id, position, priority_score, estimated_wait_time) VALUES
  ('550e8400-e29b-41d4-a716-446655440040','550e8400-e29b-41d4-a716-446655440030',1,95.5,2),
  ('550e8400-e29b-41d4-a716-446655440041','550e8400-e29b-41d4-a716-446655440031',2,87.3,5),
  ('550e8400-e29b-41d4-a716-446655440042','550e8400-e29b-41d4-a716-446655440032',3,76.8,8),
  ('550e8400-e29b-41d4-a716-446655440043','550e8400-e29b-41d4-a716-446655440033',4,92.1,3),
  ('550e8400-e29b-41d4-a716-446655440044','550e8400-e29b-41d4-a716-446655440034',5,83.7,6)
ON CONFLICT (id) DO NOTHING;

INSERT INTO fairness_metrics (id, metric_name, metric_value, group_id, city_id) VALUES
  ('550e8400-e29b-41d4-a716-446655440050','wait_time_ratio',0.85,'550e8400-e29b-41d4-a716-446655440001','550e8400-e29b-41d4-a716-446655440010'),
  ('550e8400-e29b-41d4-a716-446655440051','allocation_rate',0.78,'550e8400-e29b-41d4-a716-446655440002','550e8400-e29b-41d4-a716-446655440011'),
  ('550e8400-e29b-41d4-a716-446655440052','satisfaction_score',0.92,'550e8400-e29b-41d4-a716-446655440003','550e8400-e29b-41d4-a716-446655440012'),
  ('550e8400-e29b-41d4-a716-446655440053','fairness_index',0.88,'550e8400-e29b-41d4-a716-446655440004','550e8400-e29b-41d4-a716-446655440013'),
  ('550e8400-e29b-41d4-a716-446655440054','efficiency_score',0.91,'550e8400-e29b-41d4-a716-446655440005','550e8400-e29b-41d4-a716-446655440014')
ON CONFLICT (id) DO NOTHING;
