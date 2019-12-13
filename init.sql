-- 客户公司表
CREATE TABLE crm_user_company (
    id serial primary key,
    name text DEFAULT '' NOT null,
    name_abbr varchar(32) DEFAULT '' NOT NULL, -- 公司名称的英文缩略，可能将其作为url的第一个部分/dchj/xxx
    created timestamp(6) without time zone NOT NULL,
    limit_dial smallint default 0,             -- 限制的外呼同一个号码的次数
    limit_caller smallint default 0            -- 限制主叫一天能呼几通电话
);

COMMENT ON COLUMN crm_user_company.name_abbr IS '公司名称的英文缩略';

-- origin
-- user表
CREATE TABLE crm_backend_user (
    id serial primary key,
    real_name character varying(32) DEFAULT ''::character varying NOT NULL,
    user_name character varying(24) DEFAULT ''::character varying NOT NULL,
    user_pwd text DEFAULT ''::text NOT NULL,
    is_super boolean DEFAULT false NOT NULL,
    status integer DEFAULT 0 NOT NULL,
    mobile character varying(16) DEFAULT ''::character varying NOT NULL,
    email character varying(256) DEFAULT ''::character varying NOT NULL,
    avatar character varying(256) DEFAULT ''::character varying NOT NULL,
    leader_id integer DEFAULT 0,
    user_company_id integer DEFAULT 0
);

COMMENT ON COLUMN crm_backend_user.user_company_id IS '公司表对应的id';

-- 角色表
CREATE TABLE crm_role (
    id serial primary key,
    name text DEFAULT ''::text NOT NULL,
    seq integer DEFAULT 0 NOT NULL,
    user_company_id integer DEFAULT 0
);
-- 角色用户关系表
CREATE TABLE crm_role_backenduser_rel (
    id serial primary key,
    role_id integer NOT NULL,
    backend_user_id integer NOT NULL,
    created timestamp(6) without time zone NOT NULL
);
-- 管理controller action
CREATE TABLE crm_resource (
    id serial primary key,
    name character varying(64) DEFAULT ''::character varying NOT NULL,
    parent_id integer,
    rtype integer DEFAULT 0 NOT NULL,
    seq integer DEFAULT 0 NOT NULL,
    icon character varying(32) DEFAULT ''::character varying NOT NULL,
    url_for character varying(256) DEFAULT ''::character varying NOT NULL
);
-- 资源表，管理controller，model
CREATE TABLE crm_role_resource_rel (
    id serial primary key,
    role_id integer NOT NULL,
    resource_id integer NOT NULL,
    created timestamp(6) with time zone NOT NULL
);


-------------------------------------------------------------------------------


-- 座席表
CREATE TABLE crm_agent (
    id serial primary key,
    ext_no varchar(25) DEFAULT '' NOT null UNIQUE,                  -- 分机号码
    ext_pwd text DEFAULT '' NOT null,                               -- 分机密码
    gateway_phone_number varchar(25) DEFAULT '',    -- 中继号
    origination_caller_id_number varchar(25) DEFAULT '' NOT Null,   -- 透传号码
    backend_user_id int NOT NULL REFERENCES crm_backend_user(id),
    call_way smallint DEFAULT 0, -- 呼叫方式：0.拨座席，再拨客户，1.先拨客户，放音，再拨座席，2.先拨客户，ivr导航，再根据按键接通座席。默认0 
    gateway_id int DEFAULT 0, -- 网关id
    param text DEFAULT '', -- 呼叫方式对应的extension，转接的没值，放音的写入彩铃，ivr导航的填入ivr extension
    default_trunk int default 0 -- 跟call_gateway里的gateway_type取值一样的，0代表电话线中继，1代表插手机卡中继
);

COMMENT ON COLUMN crm_agent.ext_no IS '分机号码';
COMMENT ON COLUMN crm_agent.ext_pwd IS '分机密码';
COMMENT ON COLUMN crm_agent.gateway_phone_number IS '中继号';
COMMENT ON COLUMN crm_agent.origination_caller_id_number IS '透传号码';
COMMENT ON COLUMN crm_agent.call_way IS '呼叫方式：0.拨座席，再拨客户，1.先拨客户，放音，再拨座席，2.先拨客户，ivr导航，再根据按键接通座席。默认0';
COMMENT ON COLUMN crm_agent.param IS '呼叫方式对应的extension，转接的没值，放音的写入彩铃，ivr导航的填入ivr extension';


-- 网关表
-- CREATE TABLE call_gateway (
--     id serial primary key,
--     name varchar(255) DEFAULT '',
--     sip_name varchar(255) DEFAULT '' UNIQUE,
--     ip varchar(255) DEFAULT ''
-- );


-- 作了动态客户表，不再显示
-- 客户表
-- CREATE TABLE crm_user_client (
--     id serial primary key,
--     company_name text DEFAULT '',   -- 公司名称
--     address text DEFAULT '',        -- 地址
--     contact varchar(255) DEFAULT '', -- 联系方式
--     is_answer boolean DEFAULT false NOT NULL, -- 是否接听：bool
--     name varchar(60) DEFAULT '' NOT NULL, -- 客户姓名
--     mobile_phone varchar(25) DEFAULT '' NOT null, -- 手机
--     trend smallint DEFAULT 0, -- 客户意向:0:无意向；1：中等；2：有意向
--     email varchar(255) DEFAULT '', -- 邮箱
--     backend_user_id int NOT NULL REFERENCES crm_backend_user(id), -- 归属人
--     client_source varchar(255) DEFAULT '', -- 线索来源:qq,wx,微博，网站
--     comment text DEFAULT '', -- 备注
--     created timestamp(6) without time zone NOT NULL, -- 添加时间
--     lastest_communicated timestamp(6) without time zone -- 最近沟通时间
-- );

COMMENT ON COLUMN crm_user_client.company_name IS '客户的客户的公司名称';
COMMENT ON COLUMN crm_user_client.address IS '客户的客户的地址';
COMMENT ON COLUMN crm_user_client.is_answer IS '是否接听';
COMMENT ON COLUMN crm_user_client.name IS '客户的客户的姓名';
COMMENT ON COLUMN crm_user_client.mobile_phone IS '客户的客户的手机';
COMMENT ON COLUMN crm_user_client.trend IS '客户的客户的客户意向:0:无意向；1：中等；2：有意向';
COMMENT ON COLUMN crm_user_client.email IS '客户的客户的邮箱';
COMMENT ON COLUMN crm_user_client.client_source IS '线索来源:qq,wx,微博，网站';
COMMENT ON COLUMN crm_user_client.comment IS '备注';
COMMENT ON COLUMN crm_user_client.created IS '添加时间';
COMMENT ON COLUMN crm_user_client.lastest_communicated IS '最近沟通时间';


-- 问卷调查表
CREATE TABLE crm_asq (
    id serial primary key,
    name varchar(255) DEFAULT '' NOT NULL, -- 问卷名称
    description text DEFAULT '', -- 问卷说明
    backend_user_id int NOT NULL REFERENCES crm_backend_user(id), -- 归属人
    created timestamp(6) without time zone NOT NULL, -- 添加时间
    updated timestamp(6) without time zone -- 更新时间
);

COMMENT ON TABLE crm_asq IS '问卷调查表';


-- 问卷调查表详细问卷内容
CREATE TABLE crm_asq_detail (
    id bigserial primary key,
    asq_id int DEFAULT 0, 
    qst text DEFAULT '',            -- 题目
    qst_no int DEFAULT 0 UNIQUE,    -- 题号，unique
    qst_type smallint DEFAULT 0,    -- 题目类型:0:文本，1：单选，2：多选
    qst_type_value text DEFAULT '', -- 只有单选和多选才有值，pg中我想用array来做
    created timestamp(6) without time zone NOT NULL, -- 创建时间
    updated timestamp(6) without time zone           -- 修改时间
);

COMMENT ON TABLE crm_asq_detail IS '问卷调查表详细问卷内容';
COMMENT ON COLUMN crm_asq_detail.qst IS '题目';
COMMENT ON COLUMN crm_asq_detail.qst_no IS '题号，unique';
COMMENT ON COLUMN crm_asq_detail.qst_type IS '题目类型:0:文本，1：单选，2：多选';
COMMENT ON COLUMN crm_asq_detail.qst_type_value IS '只有单选和多选才有值，pg中我想用array来做';
COMMENT ON COLUMN crm_asq_detail.created IS '创建时间';
COMMENT ON COLUMN crm_asq_detail.updated IS '修改时间';


-- 问卷调查表详细内容的答案
CREATE TABLE crm_asq_detail_answer (
    id bigserial primary key,
    asq_detail_id bigint NOT NULL REFERENCES crm_asq_detail(id),
    backend_user_id int NOT NULL REFERENCES crm_backend_user(id),
    -- uuid char(50) NOT NULL DEFAULT '',
    answer_content text DEFAULT ''  -- 问卷题目的答案
    survey_time timestamp(6) without time zone NOT NULL, -- 调查时间
    survey_phone varchar(25) NOT NULL DEFAULT ''         -- 调查号码
);

COMMENT ON TABLE crm_asq_detail_answer IS '问卷调查表详细问卷内容的答案';
COMMENT ON COLUMN crm_asq_detail_answer.answer_content IS '问卷题目的答案';
COMMENT ON COLUMN crm_asq_detail_answer.survey_time IS '调查时间';
COMMENT ON COLUMN crm_asq_detail_answer.survey_phone IS '调查号码';


-- 短信账号设置
CREATE TABLE crm_sms_set (
    id serial primary key,
    account varchar(32) DEFAULT '', 
    password varchar(32) DEFAULT '',
    signature varchar(50) DEFAULT '',
    backend_user_id int  NOT NULL REFERENCES crm_backend_user(id),
    created timestamp(6) without time zone NOT NULL,
    updated timestamp(6) without time zone,
    user_company_id int default 0;
);

COMMENT ON TABLE crm_sms_set IS '短信账号设置';
COMMENT ON COLUMN crm_sms_set.account IS '与短信平台对接的账号';
COMMENT ON COLUMN crm_sms_set.password IS '与短信平台对接的密码';
COMMENT ON COLUMN crm_sms_set.signature IS '公司签名：如：【德诚黄金】';


-- 短信模板设置
CREATE TABLE crm_sms_template (
    id serial primary key,
    title varchar(255) DEFAULT '',
    content text DEFAULT '',
    classify smallint DEFAULT 0, -- 模板类型：0：业务短信、1：营销短信。
    backend_user_id int NOT NULL REFERENCES crm_backend_user(id),
    created timestamp(6) without time zone NOT NULL,
    updated timestamp(6) without time zone,
    state int default 0, -- 0:新建，1：审核通过
    user_company_id int default 0
);

COMMENT ON TABLE crm_sms_template IS '短信模板设置';
COMMENT ON COLUMN crm_sms_template.title IS '标题';
COMMENT ON COLUMN crm_sms_template.content IS '内容';
COMMENT ON COLUMN crm_sms_template.classify IS '模板类型：0：业务短信、1：营销短信。';
COMMENT ON COLUMN crm_sms_template.state IS '状态:0:新建，1：审核通过';


-- 短信记录
CREATE TABLE crm_sms_record (
    id bigserial primary key,
    mobile char(11) NOT NULL, -- 发送号码
    content text NOT NULL, -- 发送内容
    result smallint NOT NULL, -- 发送结果: 0: 失败， 1：成功。
    send_time timestamp(6) without time zone NOT NULL, -- 发送时间
    classify smallint DEFAULT 0, -- 类型：0：业务短信、1：营销短信。
    backend_user_id int NOT NULL REFERENCES crm_backend_user(id) -- 发送人
);

COMMENT ON TABLE crm_sms_record IS '短信记录';
COMMENT ON COLUMN crm_sms_record.mobile IS '发送号码';
COMMENT ON COLUMN crm_sms_record.content IS '发送内容';
COMMENT ON COLUMN crm_sms_record.result IS '发送结果: 0: 失败， 1：成功。';
COMMENT ON COLUMN crm_sms_record.send_time IS '发送时间';
COMMENT ON COLUMN crm_sms_record.classify IS '类型：0：业务短信、1：营销短信。';
COMMENT ON COLUMN crm_sms_record.backend_user_id IS '发送人';


-------------------------------------------------------------------------------
---- 2019-08-21 补充
-------------------------------------------------------------------------------


-- 客户字段表（自定义客户字段）
CREATE TABLE crm_user_client_field (
    id serial primary key,
    column_name varchar(60) DEFAULT '',     -- 列的英文名
    field_name varchar(60) DEFAULT '',      -- 列的中文名
    field_species smallint DEFAULT 0,       -- 字段的所属。0：自定义，1：系统默认
    field_type smallint DEFAULT 0,          -- 字段类型。0：单行文本，1：多行文本，2：单选，3：多选，4：日期，5：数字
    -- field_type_value VARCHAR(255)[] DEFAULT '',  -- 单选，多选的值，pg中用数组来做
    field_type_value text DEFAULT '',  -- 单选，多选的值，pg中用数组来做
    list_show boolean DEFAULT false,        -- 列表显示。true：显示，false：不显示
    add_show boolean DEFAULT false,         -- 添加显示。true：显示，false：不显示
    query_show boolean DEFAULT false,       -- 查询显示。true：显示，false：不显示
    -- required boolean DEFAULT false,         -- 是否必填
    required smallint DEFAULT 0,         -- 是否必填
    field_color char(7) DEFAULT '#000000',  -- 列字体颜色
    user_company_id int DEFAULT 0           -- 用户所属公司id。
);

COMMENT ON TABLE crm_user_client_field IS '定制客户字段表';
COMMENT ON COLUMN crm_user_client_field.column_name IS '列的英文名';
COMMENT ON COLUMN crm_user_client_field.field_name IS '列的中文名';
COMMENT ON COLUMN crm_user_client_field.field_species IS '字段的所属。0：自定义，1：系统默认';
COMMENT ON COLUMN crm_user_client_field.field_type IS '字段类型。0：单行文本，1：多行文本，2：单选，3：多选，4：日期，5：数字';
COMMENT ON COLUMN crm_user_client_field.field_type_value IS '单选，多选的值，pg中用数组来做';
COMMENT ON COLUMN crm_user_client_field.list_show IS '列表显示。true：显示，false：不显示';
COMMENT ON COLUMN crm_user_client_field.add_show IS '添加显示。true：显示，false：不显示';
COMMENT ON COLUMN crm_user_client_field.query_show IS '查询显示。true：显示，false：不显示';
COMMENT ON COLUMN crm_user_client_field.required IS '是否必填';
COMMENT ON COLUMN crm_user_client_field.field_color IS '列字体颜色';
COMMENT ON COLUMN crm_user_client_field.user_company_id IS '用户所属公司id。';

-- 将bool转为smallint类型。
alter table crm_user_client_field alter column required type smallint using case required when 'false' then 0 else 1 end;
alter table crm_user_client_field alter column required set default 0;
-- 联合唯一约束
alter table crm_user_client_field add constraint unique_column_name_user_company_id unique(column_name, user_company_id);
alter table crm_user_client_field add constraint unique_field_name unique(field_name);

-- 客户信息表
CREATE TABLE crm_user_client (
    id bigserial primary key,
    name varchar(60) DEFAULT '', -- 客户姓名
    mobile_phone varchar(25) DEFAULT '', -- 手机
    contact_phone varchar(25) DEFAULT '', -- 联系电话
    backend_user_id int DEFAULT 0, -- 创建人   
    created timestamp without time zone NOT NULL, -- 添加时间
    comment text DEFAULT '', -- 备注
    address text DEFAULT '', -- 地址
    dial_state varchar(255) default '',
    state smallint DEFAULT 0, -- 0：无意向，1：中等意向，2：有意向
    feature int DEFAULT 0, -- 客户特征
    complaint text DEFAULT '', -- 是否投诉
    latest_communicated timestamp without time zone, -- 最近沟通时间
    clue_from smallint DEFAULT 0, -- 线索来源。0：百度推广，1：QQ咨询，2：其他
    email varchar(255) DEFAULT '', -- 邮箱
    belong_backend_user_id int default 0, -- 归属人，0表示放入公共池
    user_company_id int default 0, -- 因为要添加公共池，所以就要加这个
    updated timestamp without time zone NOT NULL, -- 更新时间，排序根据这个排的
    column1 text DEFAULT '',
    column2 text DEFAULT '',
    column3 text DEFAULT '',
    column4 text DEFAULT '',
    column5 text DEFAULT '',
    column6 text DEFAULT '',
    column7 text DEFAULT '',
    column8 text DEFAULT '',
    column9 text DEFAULT '',
    column10 text DEFAULT '',
    column11 text DEFAULT '',
    column12 text DEFAULT '',
    column13 text DEFAULT '',
    column14 text DEFAULT '',
    column15 text DEFAULT '',
    column16 text DEFAULT ''
);

COMMENT ON TABLE crm_user_client IS '客户信息表';
COMMENT ON COLUMN crm_user_client.name IS '客户姓名';
COMMENT ON COLUMN crm_user_client.mobile_phone IS '手机';
COMMENT ON COLUMN crm_user_client.contact_phone IS '联系电话';
COMMENT ON COLUMN crm_user_client.backend_user_id IS '归属人';
COMMENT ON COLUMN crm_user_client.created IS '添加时间';
COMMENT ON COLUMN crm_user_client.comment IS '备注';
COMMENT ON COLUMN crm_user_client.address IS '地址';
COMMENT ON COLUMN crm_user_client.state IS '1：无意向，2：中等意向，3：有意向';
COMMENT ON COLUMN crm_user_client.feature IS '客户特征';
COMMENT ON COLUMN crm_user_client.complaint IS '是否投诉';
COMMENT ON COLUMN crm_user_client.latest_communicated IS '最近沟通时间';
COMMENT ON COLUMN crm_user_client.clue_from IS '线索来源。1：百度推广，2：QQ咨询，3：其他';
COMMENT ON COLUMN crm_user_client.belong_backend_user_id IS '归属人，0表示放入公共池';
COMMENT ON COLUMN crm_user_client.user_company_id IS '因为要添加公共池，所以就要加这个';
COMMENT ON COLUMN crm_user_client.updated IS '更新时间，排序根据这个排的';


-------------------------------------------------------------------------------
-- backend_user表的递归查询
-- 
with RECURSIVE le (id, real_name, user_name, path, depath) as (
    select id, real_name, user_name, Array[id] as path, 1 as depath from crm_backend_user where id=6 
    UNION all
    select u1.id, u1.real_name, u1.user_name, u2.path||u1.id, u2.depath+1 from crm_backend_user u1, le u2 where u1.leader_id=u2.id
)
select * from le 


-------------------------------------------------------------------------------
-- crm_task 2019/11/1 任务
-- 
CREATE TABLE crm_task (
    id serial primary key,
    name varchar(255) default '',
    state smallint default 0, -- 0：新建，1：进行中，2：暂停，3：已完结
    created timestamp without time zone NOT NULL,
    backend_user_id int default 0,
    user_company_id int default 0,
    desc text default ''
);

COMMENT ON TABLE crm_task IS '任务表';
COMMENT ON COLUMN crm_task.state IS '0：新建，1：进行中，2：暂停，3：已完结';

--
-- 任务细节 预计完成时间
--
CREATE TABLE crm_task_detail (
    id bigserial primary key,
    task_id int NOT NULL default 0,
    user_client_id int NOT NULL default 0,
    call_state text default '', -- 未拨打，呼不通，无人接听，拒接，挂机，已接通，未接通
    belong_user_id int default 0
);

COMMENT ON TABLE crm_task_detail IS '任务细节表';
COMMENT ON COLUMN crm_task_detail.call_state IS '备注呼叫情况';


-------------------------------------------------------------------------------
-- crm_gateway_phone 每个网关里的每个号码
-- 
CREATE TABLE crm_gateway_phone (
    id serial primary key,
    phone varchar(32) default '',
    gateway_id int default 0,
    agent_id int default 0,
    created timestamp without time zone NOT NULL,
);

COMMENT ON TABLE crm_gateway_phone IS '网关号码表';
COMMENT ON COLUMN crm_gateway_phone.created IS '创建日期';


-------------------------------------------------------------------------------
-- 发送邮件设置
-------------------------------------------------------------------------------

--
-- 邮件服务器设置
--
CREATE TABLE crm_email_server (
    id serial primary key,
    send_host varchar(255) DEFAULT '', 
    send_port int default 25,               -- ssl: 465
    send_ssl boolean default false,
    receive_host varchar(32) DEFAULT '',
    receive_port int default 110,           -- ssl: 995
    receive_ssl boolean default false,
    created_user_id int  NOT NULL REFERENCES crm_backend_user(id),
    created timestamp(6) without time zone NOT NULL,
    updated timestamp(6) without time zone,
    user_company_id int default 0
);

COMMENT ON TABLE crm_email_server IS '邮件服务器设置';
COMMENT ON COLUMN crm_email_server.send_ssl IS '发送是否使用ssl';
COMMENT ON COLUMN crm_email_server.receive_ssl IS '接受是否使用ssl';
COMMENT ON COLUMN crm_email_server.created_user_id IS '创建人';

--
-- 邮件账号设置
--
CREATE TABLE crm_email_account (
    id serial primary key,
    email_server_id int NOT NULL REFERENCES crm_email_server(id),
    user_name varchar(255) DEFAULT '', 
    password char(128) DEFAULT '',
    from_name varchar(255) DEFAULT '', -- 显示名称
    backend_user_id int  NOT NULL REFERENCES crm_backend_user(id),
    created timestamp(6) without time zone NOT NULL,
    updated timestamp(6) without time zone,
    user_company_id int default 0
);

COMMENT ON TABLE crm_email_account IS '邮件账户设置';
COMMENT ON COLUMN crm_email_account.password IS 'sha256';
COMMENT ON COLUMN crm_email_account.from_name IS '显示名称';


-- 邮件模板设置
CREATE TABLE crm_email_template (
    id serial primary key,
    title varchar(255) DEFAULT '',
    content text DEFAULT '',
    classify smallint DEFAULT 0, -- 模板类型：0：通知邮件、1：营销邮件。
    backend_user_id int NOT NULL REFERENCES crm_backend_user(id),
    created timestamp(6) without time zone NOT NULL,
    updated timestamp(6) without time zone,
    state int default 0, -- 0:新建，1：审核通过
    user_company_id int default 0
);

COMMENT ON TABLE crm_email_template IS '邮件模板设置';
COMMENT ON COLUMN crm_email_template.title IS '标题';
COMMENT ON COLUMN crm_email_template.content IS '内容';
COMMENT ON COLUMN crm_email_template.classify IS '模板类型：0：通知邮件、1：营销邮件。';
COMMENT ON COLUMN crm_email_template.state IS '状态:0:新建，1：审核通过';


-- 邮件记录
CREATE TABLE crm_email_record (
    id bigserial primary key,
    sender text NOT NULL,       -- 发件人
    recevier text NOT NULL,     -- 收件人
    carbon_copy text NOT NULL,  -- 抄送
    content text NOT NULL,      -- 发送内容
    attach_name text NOT NULL,  -- 附件
    result smallint NOT NULL,   -- 发送结果: 0: 失败， 1：成功。
    send_time timestamp(6) without time zone NOT NULL, -- 发送时间
    classify smallint DEFAULT 0, -- 类型：0：通知邮件、1：营销邮件。
    backend_user_id int NOT NULL REFERENCES crm_backend_user(id) -- 发送人
);

COMMENT ON TABLE crm_email_record IS '邮件记录';
COMMENT ON COLUMN crm_email_record.sender IS '发件人,一般为一个人';
COMMENT ON COLUMN crm_email_record.recevier IS '收件人,逗号分隔的多人';
COMMENT ON COLUMN crm_email_record.carbon_copy IS '抄送,逗号分隔的多人';
COMMENT ON COLUMN crm_email_record.content IS '发送内容';
COMMENT ON COLUMN crm_email_record.attach_name IS '附件文件名,逗号分隔的多个附件文件名';
COMMENT ON COLUMN crm_email_record.result IS '发送结果: 0: 失败， 1：成功。';
COMMENT ON COLUMN crm_email_record.send_time IS '发送时间';
COMMENT ON COLUMN crm_email_record.classify IS '类型：0：通知邮件、1：营销邮件。';
COMMENT ON COLUMN crm_email_record.backend_user_id IS '发送人';


-------------------------------------------------------------------------------
-- 发送邮件设置
-------------------------------------------------------------------------------