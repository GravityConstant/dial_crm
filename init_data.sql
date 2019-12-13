insert into crm_user_company (name, name_abbr, created) 
    values ('和企时代', 'heqi', '2019-08-20 10:50:50');

insert into crm_backend_user(real_name, user_name, user_pwd, is_super, status, avatar, user_company_id) 
    values ('admin', 'admin', 'e10adc3949ba59abbe56e057f20f883e', 'true', 1, '/static/upload/1.jpg', 1);

insert  into crm_role(name,seq) values ('超级管理员',1);

insert  into crm_role_backenduser_rel(role_id,backend_user_id,created) values (1,1,'2019-08-20 11:00:00');

-- insert  into crm_resource(rtype,name,parent_id,seq,icon,url_for) values 
--     (7,1,'权限管理',8,100,'fa fa-balance-scale',''),
--     (8,0,'系统菜单',NULL,200,'',''),
--     (9,1,'资源管理',7,100,'','ResourceController.Index'),
--     (12,1,'角色管理',7,100,'','RoleController.Index'),
--     (13,1,'用户管理',7,100,'','BackendUserController.Index'),
--     (14,1,'系统管理',8,90,'fa fa-gears',''),
--     (21,0,'业务菜单',NULL,170,'',''),
--     (22,1,'课程资源(空)',21,100,'fa fa-book',''),
--     (23,1,'日志管理(空)',14,100,'',''),
--     (25,2,'编辑',9,100,'fa fa-pencil','ResourceController.Edit'),
--     (26,2,'编辑',13,100,'fa fa-pencil','BackendUserController.Edit'),
--     (27,2,'删除',9,100,'fa fa-trash','ResourceController.Delete'),
--     (29,2,'删除',13,100,'fa fa-trash','BackendUserController.Delete'),
--     (30,2,'编辑',12,100,'fa fa-pencil','RoleController.Edit'),
--     (31,2,'删除',12,100,'fa fa-trash','RoleController.Delete'),
--     (32,2,'分配资源',12,100,'fa fa-th','RoleController.Allocate'),
--     (35,1,' 首页',NULL,100,'fa fa-dashboard','HomeController.Index');


insert into crm_resource(rtype,name,parent_id,seq,icon,url_for) values 
    (1,' 首页',NULL,100,'fa fa-dashboard','HomeController.Index'),
    (0,'系统菜单',NULL,200,'',''),
    (0,'业务菜单',NULL,170,'','');

insert into crm_resource(rtype,name,parent_id,seq,icon,url_for) values
    (1,'权限管理',2,100,'fa fa-balance-scale',''),
    (1,'系统管理',2,90,'fa fa-gears',''); 


insert into crm_resource(rtype,name,parent_id,seq,icon,url_for) values
    (1,'资源管理',4,100,'','ResourceController.Index'),
    (1,'角色管理',4,100,'','RoleController.Index'),
    (1,'用户管理',4,100,'','BackendUserController.Index');

insert into crm_resource(rtype,name,parent_id,seq,icon,url_for) values
    (2,'编辑',6,100,'fa fa-pencil','ResourceController.Edit'),
    (2,'删除',6,100,'fa fa-trash','ResourceController.Delete');

insert into crm_resource(rtype,name,parent_id,seq,icon,url_for) values
    (2,'编辑',7,100,'fa fa-pencil','RoleController.Edit'),
    (2,'删除',7,100,'fa fa-trash','RoleController.Delete'),
    (2,'分配资源',7,100,'fa fa-th','RoleController.Allocate');

insert into crm_resource(rtype,name,parent_id,seq,icon,url_for) values
    (2,'编辑',8,100,'fa fa-pencil','BackendUserController.Edit'),
    (2,'删除',8,100,'fa fa-trash','BackendUserController.Delete');


-------------------------------------------------------------------------------
---- 2019-08-21 补充
-------------------------------------------------------------------------------

-- 字段类型。0：单行文本，1：多行文本，2：单选，3：多选，4：日期，5：数字
insert into crm_user_client_field(column_name, field_name, field_species, field_type, field_type_value, list_show, add_show, query_show, required, field_color) values 
    ('name', '客户姓名', '1', '0', '', 't', 't', 't', 't', '#191315'),
    ('mobile_phone', '手机', '1', '0', '', 't', 't', 't', 't', '#131719'),
    ('contact_phone', '联系电话', '1', '0', '', 't', 't', 't', 'f', '#2b3337'),
    ('backend_user_id', '归属人', '1', '0', '', 't', 't', 't', 't', '#131719'),
    ('created', '添加时间', '1', '4', '', 't', 't', 't', '', '#131719'),
    ('comment', '备注', '1', '1', '', 't', 't', 't', 'f', '#131719'),
    ('address', '地址', '1', '0', '', 't', 't', 't', 'f', '#131719'),
    ('state', '客户状态', '1', '2', '无意向,普通,有意向', 't', 't', 't', 't', '#131719'),
    ('dial_state', '接通状态', '1', '0', '', 't', 't', 't', '1', '#131719'),
    ('feature', '客户特征', '1', '0', '', 't', 't', 't', 'f', '#a19313'),
    ('complaint', '是否投诉', '1', '0', '', 't', 't', 't', 'f', '#1a1f21'),
    ('latest_communicated', '最近沟通时间', '1', '4', '', 't', 't', 't', 'f', '#191315'),
    ('clue_from', '线索来源', '1', '0', '', 't', 't', 't', 'f', '#676a6c'),
    ('email', '邮箱', '1', '0', '', 't', 't', 't', 'f', '#676a6c');