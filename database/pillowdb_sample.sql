-- Sample Data for User Management System
-- This file contains sample roles, permissions, and relationships for testing the RBAC system

-- ===========================================
-- ROLES
-- ===========================================

-- Super Admin Role (highest level)
INSERT INTO "Roles" (id, name, description) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'super_admin', 'Super Administrator with full system access'),
('550e8400-e29b-41d4-a716-446655440001', 'admin', 'Administrator with management capabilities'),
('550e8400-e29b-41d4-a716-446655440002', 'manager', 'Manager with team management capabilities'),
('550e8400-e29b-41d4-a716-446655440003', 'user', 'Regular user with basic access'),
('550e8400-e29b-41d4-a716-446655440004', 'viewer', 'Read-only user with limited access');

-- ===========================================
-- PERMISSIONS
-- ===========================================

-- User Management Permissions
INSERT INTO "Permissions" (id, name, description, scope_level) VALUES
('660e8400-e29b-41d4-a716-446655440000', 'manage_users', 'Create, read, update, and delete users', 'system'),
('660e8400-e29b-41d4-a716-446655440001', 'view_users', 'View user information', 'organization'),
('660e8400-e29b-41d4-a716-446655440002', 'manage_own_profile', 'Manage own user profile', 'user'),

-- Role Management Permissions
('660e8400-e29b-41d4-a716-446655440003', 'manage_roles', 'Create, read, update, and delete roles', 'system'),
('660e8400-e29b-41d4-a716-446655440004', 'view_roles', 'View role information', 'organization'),
('660e8400-e29b-41d4-a716-446655440005', 'assign_roles', 'Assign roles to users', 'organization'),

-- Permission Management Permissions
('660e8400-e29b-41d4-a716-446655440006', 'manage_permissions', 'Create, read, update, and delete permissions', 'system'),
('660e8400-e29b-41d4-a716-446655440007', 'view_permissions', 'View permission information', 'organization'),

-- Organization Management Permissions
('660e8400-e29b-41d4-a716-446655440008', 'manage_organizations', 'Create, read, update, and delete organizations', 'system'),
('660e8400-e29b-41d4-a716-446655440009', 'view_organizations', 'View organization information', 'organization'),
('660e8400-e29b-41d4-a716-446655440010', 'manage_own_organization', 'Manage own organization', 'organization'),

-- Audit and Logging Permissions
('660e8400-e29b-41d4-a716-446655440011', 'view_audit_logs', 'View system audit logs', 'system'),
('660e8400-e29b-41d4-a716-446655440012', 'manage_audit_logs', 'Manage audit log settings', 'system'),

-- Content and Data Permissions
('660e8400-e29b-41d4-a716-446655440013', 'manage_content', 'Create, read, update, and delete content', 'organization'),
('660e8400-e29b-41d4-a716-446655440014', 'view_content', 'View content', 'organization'),
('660e8400-e29b-41d4-a716-446655440015', 'export_data', 'Export system data', 'organization'),

-- System Administration Permissions
('660e8400-e29b-41d4-a716-446655440016', 'system_admin', 'Full system administration access', 'system'),
('660e8400-e29b-41d4-a716-446655440017', 'view_system_info', 'View system information and statistics', 'system'),
('660e8400-e29b-41d4-a716-446655440018', 'manage_system_settings', 'Manage system-wide settings', 'system');

-- ===========================================
-- ROLE-PERMISSION RELATIONSHIPS
-- ===========================================

-- Super Admin - ALL permissions
INSERT INTO "Role_Permissions" (role_id, permission_id) VALUES
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
INSERT INTO "Users" (id, username, email, password_hash, is_active) VALUES
('43739cef-f82b-4b74-8e11-2c9c907300d1', 'superadmin', 'superadmin@pillow.com', '$2a$10$8K3VZ6Y8QX8QX8QX8QX8QeF8QX8QX8QX8QX8QX8QX8QX8QX8QX8Q', true);

-- Other sample users
INSERT INTO "Users" (id, username, email, password_hash, is_active) VALUES
('550e8400-e29b-41d4-a716-446655440100', 'admin_user', 'admin@pillow.com', '$2a$10$8K3VZ6Y8QX8QX8QX8QX8QeF8QX8QX8QX8QX8QX8QX8QX8QX8QX8Q', true),
('550e8400-e29b-41d4-a716-446655440101', 'manager_user', 'manager@pillow.com', '$2a$10$8K3VZ6Y8QX8QX8QX8QX8QeF8QX8QX8QX8QX8QX8QX8QX8QX8QX8Q', true),
('550e8400-e29b-41d4-a716-446655440102', 'regular_user', 'user@pillow.com', '$2a$10$8K3VZ6Y8QX8QX8QX8QX8QeF8QX8QX8QX8QX8QX8QX8QX8QX8QX8Q', true),
('550e8400-e29b-41d4-a716-446655440103', 'viewer_user', 'viewer@pillow.com', '$2a$10$8K3VZ6Y8QX8QX8QX8QX8QeF8QX8QX8QX8QX8QX8QX8QX8QX8QX8Q', true);

-- ===========================================
-- USER-ROLE ASSIGNMENTS
-- ===========================================

-- Assign Super Admin role to the specified user
INSERT INTO "User_Roles" (user_id, role_id, scope) VALUES
('43739cef-f82b-4b74-8e11-2c9c907300d1', '550e8400-e29b-41d4-a716-446655440000', 'system');

-- Assign other roles to sample users
INSERT INTO "User_Roles" (user_id, role_id, scope) VALUES
('550e8400-e29b-41d4-a716-446655440100', '550e8400-e29b-41d4-a716-446655440001', 'system'), -- admin
('550e8400-e29b-41d4-a716-446655440101', '550e8400-e29b-41d4-a716-446655440002', 'organization'), -- manager
('550e8400-e29b-41d4-a716-446655440102', '550e8400-e29b-41d4-a716-446655440003', 'organization'), -- user
('550e8400-e29b-41d4-a716-446655440103', '550e8400-e29b-41d4-a716-446655440004', 'organization'); -- viewer

-- ===========================================
-- SAMPLE ORGANIZATIONS
-- ===========================================

INSERT INTO "Organizations" (id, name, description, domain, managed_by) VALUES
('770e8400-e29b-41d4-a716-446655440000', 'Pillow Technologies', 'Main organization for Pillow Technologies', 'pillow.com', '43739cef-f82b-4b74-8e11-2c9c907300d1'),
('770e8400-e29b-41d4-a716-446655440001', 'Development Team', 'Software development department', 'dev.pillow.com', '550e8400-e29b-41d4-a716-446655440100'),
('770e8400-e29b-41d4-a716-446655440002', 'Marketing Team', 'Marketing and sales department', 'marketing.pillow.com', '550e8400-e29b-41d4-a716-446655440100');

-- ===========================================
-- USER-ORGANIZATION MEMBERSHIPS
-- ===========================================

INSERT INTO "User_Organizations" (id, user_id, org_id, role_id, invited_by) VALUES
('880e8400-e29b-41d4-a716-446655440000', '43739cef-f82b-4b74-8e11-2c9c907300d1', '770e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '43739cef-f82b-4b74-8e11-2c9c907300d1'),
('880e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440100', '770e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440001', '43739cef-f82b-4b74-8e11-2c9c907300d1'),
('880e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440101', '770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440100'),
('880e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440102', '770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440101'),
('880e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440103', '770e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440100');

-- ===========================================
-- SAMPLE AUDIT LOGS
-- ===========================================

INSERT INTO "Audit_Log" (id, user_id, action, details) VALUES
('990e8400-e29b-41d4-a716-446655440000', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'USER_CREATED', 'Super admin user created during system initialization'),
('990e8400-e29b-41d4-a716-446655440001', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'ROLE_CREATED', 'Super admin role created with full permissions'),
('990e8400-e29b-41d4-a716-446655440002', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'ORGANIZATION_CREATED', 'Main organization created'),
('990e8400-e29b-41d4-a716-446655440003', '43739cef-f82b-4b74-8e11-2c9c907300d1', 'USER_ROLE_ASSIGNED', 'Super admin role assigned to super admin user');

-- ===========================================
-- USEFUL QUERIES FOR TESTING
-- ===========================================

-- Check Super Admin User
-- SELECT u.id, u.username, u.email, r.name as role_name
-- FROM "Users" u
-- JOIN "User_Roles" ur ON u.id = ur.user_id
-- JOIN "Roles" r ON ur.role_id = r.id
-- WHERE u.id = '43739cef-f82b-4b74-8e11-2c9c907300d1';

-- Check User Permissions
-- SELECT u.username, r.name as role_name, p.name as permission_name, p.scope_level
-- FROM "Users" u
-- JOIN "User_Roles" ur ON u.id = ur.user_id
-- JOIN "Roles" r ON ur.role_id = r.id
-- JOIN "Role_Permissions" rp ON r.id = rp.role_id
-- JOIN "Permissions" p ON rp.permission_id = p.id
-- WHERE u.id = '43739cef-f82b-4b74-8e11-2c9c907300d1'
-- ORDER BY r.name, p.name;

-- Check Organization Structure
-- SELECT o.name as org_name, u.username, r.name as role_name
-- FROM "Organizations" o
-- JOIN "User_Organizations" uo ON o.id = uo.org_id
-- JOIN "Users" u ON uo.user_id = u.id
-- JOIN "Roles" r ON uo.role_id = r.id
-- ORDER BY o.name, u.username;