-- Sample Data for User Management System
-- This file contains sample roles, permissions, and relationships for testing the RBAC system

-- ===========================================
-- ROLES
-- ===========================================

-- Super Admin Role (highest level)
INSERT INTO "roles" (id, name, description) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'super_admin', 'Super Administrator with full system access'),
('550e8400-e29b-41d4-a716-446655440001', 'admin', 'Administrator with management capabilities'),
('550e8400-e29b-41d4-a716-446655440002', 'manager', 'Manager with team management capabilities'),
('550e8400-e29b-41d4-a716-446655440003', 'user', 'Regular user with basic access'),
('550e8400-e29b-41d4-a716-446655440004', 'viewer', 'Read-only user with limited access');

-- ===========================================
-- PERMISSIONS
-- ===========================================

-- User Management permissions
INSERT INTO "permissions" (id, name, description, scope_level) VALUES
('660e8400-e29b-41d4-a716-446655440000', 'manage_users', 'Create, read, update, and delete users', 'system'),
('660e8400-e29b-41d4-a716-446655440001', 'view_users', 'View user information', 'organization'),
('660e8400-e29b-41d4-a716-446655440002', 'manage_own_profile', 'Manage own user profile', 'user'),

-- Role Management permissions
('660e8400-e29b-41d4-a716-446655440003', 'manage_roles', 'Create, read, update, and delete roles', 'system'),
('660e8400-e29b-41d4-a716-446655440004', 'view_roles', 'View role information', 'organization'),
('660e8400-e29b-41d4-a716-446655440005', 'assign_roles', 'Assign roles to users', 'organization'),

-- Permission Management permissions
('660e8400-e29b-41d4-a716-446655440006', 'manage_permissions', 'Create, read, update, and delete permissions', 'system'),
('660e8400-e29b-41d4-a716-446655440007', 'view_permissions', 'View permission information', 'organization'),

-- Organization Management permissions
('660e8400-e29b-41d4-a716-446655440008', 'manage_organizations', 'Create, read, update, and delete organizations', 'system'),
('660e8400-e29b-41d4-a716-446655440009', 'view_organizations', 'View organization information', 'organization'),
('660e8400-e29b-41d4-a716-446655440010', 'manage_own_organization', 'Manage own organization', 'organization'),

-- Audit and Logging permissions
('660e8400-e29b-41d4-a716-446655440011', 'view_audit_logs', 'View system audit logs', 'system'),
('660e8400-e29b-41d4-a716-446655440012', 'manage_audit_logs', 'Manage audit log settings', 'system'),

-- Content and Data permissions
('660e8400-e29b-41d4-a716-446655440013', 'manage_content', 'Create, read, update, and delete content', 'organization'),
('660e8400-e29b-41d4-a716-446655440014', 'view_content', 'View content', 'organization'),
('660e8400-e29b-41d4-a716-446655440015', 'export_data', 'Export system data', 'organization'),

-- System Administration permissions
('660e8400-e29b-41d4-a716-446655440016', 'system_admin', 'Full system administration access', 'system'),
('660e8400-e29b-41d4-a716-446655440017', 'view_system_info', 'View system information and statistics', 'system'),
('660e8400-e29b-41d4-a716-446655440018', 'manage_system_settings', 'Manage system-wide settings', 'system'),
('660e8400-e29b-41d4-a716-446655440019', 'manage_custom_fields', 'Manage global custom fields', 'system');

-- ===========================================
-- ROLE-PERMISSION RELATIONSHIPS
-- ===========================================

-- Super Admin - ALL permissions
INSERT INTO "role_permissions" (role_id, permission_id) VALUES
-- Super Admin gets ALL permissions
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440000'), -- manage_users
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440001'), -- view_users
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440002'), -- manage_own_profile
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440003'), -- manage_roles
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440004'), -- view_roles
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440005'), -- assign_roles
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440006'), -- manage_permissions
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440007'), -- view_permissions
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440008'), -- manage_organizations
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440009'), -- view_organizations
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440010'), -- manage_own_organization
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440011'), -- view_audit_logs
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440012'), -- manage_audit_logs
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440013'), -- manage_content
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440014'), -- view_content
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440015'), -- export_data
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440016'), -- system_admin
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440017'), -- view_system_info
('550e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440018'), -- manage_system_settings

-- Admin - Most permissions except super admin specific ones
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440000'), -- manage_users
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001'), -- view_users
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440002'), -- manage_own_profile
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440003'), -- manage_roles
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440004'), -- view_roles
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440005'), -- assign_roles
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440006'), -- manage_permissions
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440007'), -- view_permissions
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440008'), -- manage_organizations
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440009'), -- view_organizations
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440010'), -- manage_own_organization
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440011'), -- view_audit_logs
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440013'), -- manage_content
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440014'), -- view_content
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440015'), -- export_data
('550e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440017'), -- view_system_info

-- Manager - Team management permissions
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001'), -- view_users
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440002'), -- manage_own_profile
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440004'), -- view_roles
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440005'), -- assign_roles
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440009'), -- view_organizations
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440010'), -- manage_own_organization
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440013'), -- manage_content
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440014'), -- view_content
('550e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440015'), -- export_data

-- User - Basic permissions
('550e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440002'), -- manage_own_profile
('550e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440014'), -- view_content

-- Viewer - Read-only permissions
('550e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440001'), -- view_users
('550e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440004'), -- view_roles
('550e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440007'), -- view_permissions
('550e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440009'), -- view_organizations
('550e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440014'); -- view_content

-- ===========================================
-- SAMPLE USERS
-- ===========================================

-- Super Admin User (the one you specified)
-- Password is "password123"
INSERT INTO "users" (id, username, email, password_hash, is_active) VALUES
('43739cef-f82b-4b74-8e11-2c9c907300d1', 'superadmin', 'superadmin@pillow.com', '$2a$10$aGHfPOCjVnSBHmg/WE5y7eKLCnH.epcLtVdMiN8B4lRR1FBkb1J7K', true);

-- Other sample users
INSERT INTO "users" (id, username, email, password_hash, is_active) VALUES
('550e8400-e29b-41d4-a716-446655440100', 'admin_user', 'admin@pillow.com', '$2a$10$aGHfPOCjVnSBHmg/WE5y7eKLCnH.epcLtVdMiN8B4lRR1FBkb1J7K', true),
('550e8400-e29b-41d4-a716-446655440101', 'manager_user', 'manager@pillow.com', '$2a$10$aGHfPOCjVnSBHmg/WE5y7eKLCnH.epcLtVdMiN8B4lRR1FBkb1J7K', true),
('550e8400-e29b-41d4-a716-446655440102', 'regular_user', 'user@pillow.com', '$2a$10$aGHfPOCjVnSBHmg/WE5y7eKLCnH.epcLtVdMiN8B4lRR1FBkb1J7K', true),
('550e8400-e29b-41d4-a716-446655440103', 'viewer_user', 'viewer@pillow.com', '$2a$10$aGHfPOCjVnSBHmg/WE5y7eKLCnH.epcLtVdMiN8B4lRR1FBkb1J7K', true);

-- ===========================================
-- USER-ROLE ASSIGNMENTS
-- ===========================================

-- Assign Super Admin role to the specified user
INSERT INTO "user_roles" (user_id, role_id, scope) VALUES
('43739cef-f82b-4b74-8e11-2c9c907300d1', '550e8400-e29b-41d4-a716-446655440000', 'system');

-- Assign other roles to sample users
INSERT INTO "user_roles" (user_id, role_id, scope) VALUES
('550e8400-e29b-41d4-a716-446655440100', '550e8400-e29b-41d4-a716-446655440001', 'system'), -- admin
('550e8400-e29b-41d4-a716-446655440101', '550e8400-e29b-41d4-a716-446655440002', 'organization'), -- manager
('550e8400-e29b-41d4-a716-446655440102', '550e8400-e29b-41d4-a716-446655440003', 'organization'), -- user
('550e8400-e29b-41d4-a716-446655440103', '550e8400-e29b-41d4-a716-446655440004', 'organization'); -- viewer

-- ===========================================
-- SAMPLE ORGANIZATIONS
-- ===========================================

INSERT INTO "organizations" (id, name, description, domain, managed_by) VALUES
('770e8400-e29b-41d4-a716-446655440000', 'Pillow Technologies', 'Main organization for Pillow Technologies', 'pillow.com', '43739cef-f82b-4b74-8e11-2c9c907300d1'),
('770e8400-e29b-41d4-a716-446655440001', 'Development Team', 'Software development department', 'dev.pillow.com', '550e8400-e29b-41d4-a716-446655440100'),
('770e8400-e29b-41d4-a716-446655440002', 'Marketing Team', 'Marketing and sales department', 'marketing.pillow.com', '550e8400-e29b-41d4-a716-446655440100');

-- ===========================================
-- USER-ORGANIZATION MEMBERSHIPS
-- ===========================================

INSERT INTO "user_organizations" (id, user_id, org_id, role_id, invited_by) VALUES
('880e8400-e29b-41d4-a716-446655440000', '43739cef-f82b-4b74-8e11-2c9c907300d1', '770e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '43739cef-f82b-4b74-8e11-2c9c907300d1'),
('880e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440100', '770e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440001', '43739cef-f82b-4b74-8e11-2c9c907300d1'),
('880e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440101', '770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440100'),
('880e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440102', '770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440101'),
('880e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440103', '770e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440100');

-- ===========================================
-- SAMPLE AUDIT LOGS
-- ===========================================

INSERT INTO "audit_log" (id, user_id, action, details) VALUES
('990e8400-e29b-41d4-a716-446655440000', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'USER_CREATED', 'Super admin user created during system initialization'),
('990e8400-e29b-41d4-a716-446655440001', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'ROLE_CREATED', 'Super admin role created with full permissions'),
('990e8400-e29b-41d4-a716-446655440002', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'ORGANIZATION_CREATED', 'Main organization created'),
('990e8400-e29b-41d4-a716-446655440003', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'USER_ROLE_ASSIGNED', 'Super admin role assigned to super admin user');

-- ===========================================
-- USEFUL QUERIES FOR TESTING
-- ===========================================

-- Check Super Admin User
-- SELECT u.id, u.username, u.email, r.name as role_name
-- FROM "users" u
-- JOIN "user_roles" ur ON u.id = ur.user_id
-- JOIN "roles" r ON ur.role_id = r.id
-- WHERE u.id = '43739cef-f82b-4b74-8e11-2c9c907300d1';

-- Check User permissions
-- SELECT u.username, r.name as role_name, p.name as permission_name, p.scope_level
-- FROM "users" u
-- JOIN "user_roles" ur ON u.id = ur.user_id
-- JOIN "roles" r ON ur.role_id = r.id
-- JOIN "role_permissions" rp ON r.id = rp.role_id
-- JOIN "permissions" p ON rp.permission_id = p.id
-- WHERE u.id = '43739cef-f82b-4b74-8e11-2c9c907300d1'
-- ORDER BY r.name, p.name;

-- Check Organization Structure
-- SELECT o.name as org_name, u.username, r.name as role_name
-- FROM "organizations" o
-- JOIN "user_organizations" uo ON o.id = uo.org_id
-- JOIN "users" u ON uo.user_id = u.id
-- JOIN "roles" r ON uo.role_id = r.id
-- ORDER BY o.name, u.username;

-- ===========================================
-- SAMPLE CUSTOM FIELDS
-- ===========================================

-- Global custom field definitions
INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440000', 'full_name', 'Full Name', 'text', true, null, '{"min_length": 2, "max_length": 100}', 1, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440001', 'phone_number', 'Phone Number', 'phone', false, null, '{"pattern": "^\\+?[1-9]\\d{1,14}$"}', 2, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440002', 'date_of_birth', 'Date of Birth', 'date', false, null, '{"min": "1900-01-01", "max": "2010-12-31"}', 3, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440003', 'department', 'Department', 'select', true, '["Engineering", "Marketing", "Sales", "HR", "Finance", "Operations"]', null, 4, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440004', 'job_title', 'Job Title', 'text', false, null, '{"min_length": 2, "max_length": 50}', 5, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440005', 'bio', 'Bio', 'textarea', false, null, '{"max_length": 500}', 6, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440006', 'years_experience', 'Years of Experience', 'number', false, null, '{"min": 0, "max": 50}', 7, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440007', 'skills', 'Skills', 'multiselect', false, '["JavaScript", "Python", "Java", "C++", "React", "Node.js", "SQL", "AWS", "Docker", "Kubernetes"]', null, 8, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440008', 'is_remote', 'Remote Work', 'boolean', false, null, null, 9, true);

INSERT INTO "custom_fields" (id, name, label, type, required, options, validation, "order", is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440009', 'emergency_contact', 'Emergency Contact', 'text', false, null, '{"min_length": 2, "max_length": 100}', 10, true);

-- ===========================================
-- SAMPLE USER CUSTOM FIELD VALUES
-- ===========================================

-- Super Admin User custom field values
INSERT INTO "user_custom_field_values" (id, user_id, field_id, value) VALUES
('bb0e8400-e29b-41d4-a716-446655440000', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440000', 'Super Administrator'),
('bb0e8400-e29b-41d4-a716-446655440001', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440001', '+1-555-0100'),
('bb0e8400-e29b-41d4-a716-446655440002', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440003', 'Engineering'),
('bb0e8400-e29b-41d4-a716-446655440003', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440004', 'Chief Technology Officer'),
('bb0e8400-e29b-41d4-a716-446655440004', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440005', 'Experienced technology leader with 15+ years in software development and system architecture.'),
('bb0e8400-e29b-41d4-a716-446655440005', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440006', '15'),
('bb0e8400-e29b-41d4-a716-446655440006', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440007', '["JavaScript", "Python", "AWS", "Docker", "Kubernetes"]'),
('bb0e8400-e29b-41d4-a716-446655440007', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440008', 'true'),
('bb0e8400-e29b-41d4-a716-446655440008', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'aa0e8400-e29b-41d4-a716-446655440009', 'Jane Administrator +1-555-0101');

-- Admin User custom field values
INSERT INTO "user_custom_field_values" (id, user_id, field_id, value) VALUES
('bb0e8400-e29b-41d4-a716-446655440010', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440000', 'Admin User'),
('bb0e8400-e29b-41d4-a716-446655440011', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440001', '+1-555-0102'),
('bb0e8400-e29b-41d4-a716-446655440012', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440002', '1985-03-15'),
('bb0e8400-e29b-41d4-a716-446655440013', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440003', 'Engineering'),
('bb0e8400-e29b-41d4-a716-446655440014', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440004', 'System Administrator'),
('bb0e8400-e29b-41d4-a716-446655440015', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440006', '8'),
('bb0e8400-e29b-41d4-a716-446655440016', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440007', '["JavaScript", "Python", "SQL", "AWS"]'),
('bb0e8400-e29b-41d4-a716-446655440017', '550e8400-e29b-41d4-a716-446655440100', 'aa0e8400-e29b-41d4-a716-446655440008', 'false');

-- Manager User custom field values
INSERT INTO "user_custom_field_values" (id, user_id, field_id, value) VALUES
('bb0e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440000', 'Manager User'),
('bb0e8400-e29b-41d4-a716-446655440021', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440001', '+1-555-0103'),
('bb0e8400-e29b-41d4-a716-446655440022', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440003', 'Marketing'),
('bb0e8400-e29b-41d4-a716-446655440023', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440004', 'Marketing Manager'),
('bb0e8400-e29b-41d4-a716-446655440024', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440005', 'Experienced marketing professional focused on digital campaigns and team leadership.'),
('bb0e8400-e29b-41d4-a716-446655440025', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440006', '10'),
('bb0e8400-e29b-41d4-a716-446655440026', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440007', '["JavaScript", "React", "Node.js"]'),
('bb0e8400-e29b-41d4-a716-446655440027', '550e8400-e29b-41d4-a716-446655440101', 'aa0e8400-e29b-41d4-a716-446655440008', 'true');

-- Regular User custom field values
INSERT INTO "user_custom_field_values" (id, user_id, field_id, value) VALUES
('bb0e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440000', 'Regular User'),
('bb0e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440001', '+1-555-0104'),
('bb0e8400-e29b-41d4-a716-446655440032', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440002', '1990-07-22'),
('bb0e8400-e29b-41d4-a716-446655440033', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440003', 'Sales'),
('bb0e8400-e29b-41d4-a716-446655440034', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440004', 'Sales Representative'),
('bb0e8400-e29b-41d4-a716-446655440035', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440006', '3'),
('bb0e8400-e29b-41d4-a716-446655440036', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440007', '["JavaScript", "React"]'),
('bb0e8400-e29b-41d4-a716-446655440037', '550e8400-e29b-41d4-a716-446655440102', 'aa0e8400-e29b-41d4-a716-446655440008', 'false');

-- Viewer User custom field values
INSERT INTO "user_custom_field_values" (id, user_id, field_id, value) VALUES
('bb0e8400-e29b-41d4-a716-446655440040', '550e8400-e29b-41d4-a716-446655440103', 'aa0e8400-e29b-41d4-a716-446655440000', 'Viewer User'),
('bb0e8400-e29b-41d4-a716-446655440041', '550e8400-e29b-41d4-a716-446655440103', 'aa0e8400-e29b-41d4-a716-446655440001', '+1-555-0105'),
('bb0e8400-e29b-41d4-a716-446655440042', '550e8400-e29b-41d4-a716-446655440103', 'aa0e8400-e29b-41d4-a716-446655440003', 'HR'),
('bb0e8400-e29b-41d4-a716-446655440043', '550e8400-e29b-41d4-a716-446655440103', 'aa0e8400-e29b-41d4-a716-446655440004', 'HR Assistant'),
('bb0e8400-e29b-41d4-a716-446655440044', '550e8400-e29b-41d4-a716-446655440103', 'aa0e8400-e29b-41d4-a716-446655440006', '2'),
('bb0e8400-e29b-41d4-a716-446655440045', '550e8400-e29b-41d4-a716-446655440103', 'aa0e8400-e29b-41d4-a716-446655440008', 'false');

-- ===========================================
-- USEFUL QUERIES FOR CUSTOM FIELDS TESTING
-- ===========================================

-- Check all custom fields
-- SELECT id, name, label, type, required, "order" FROM "custom_fields" WHERE is_active = true ORDER BY "order";

-- Check user custom field values
-- SELECT u.username, cf.label, ucfv.value
-- FROM "user_custom_field_values" ucfv
-- JOIN "users" u ON ucfv.user_id = u.id
-- JOIN "custom_fields" cf ON ucfv.field_id = cf.id
-- ORDER BY u.username, cf."order";

-- Check custom fields with their options
-- SELECT name, label, type, options FROM "custom_fields"
-- WHERE type IN ('select', 'multiselect') AND options IS NOT NULL;

-- Check validation rules
-- SELECT name, label, type, validation FROM "custom_fields"
-- WHERE validation IS NOT NULL;