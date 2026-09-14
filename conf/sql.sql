create table t_order
(
    id        bigint auto_increment primary key,
    order_id  varchar(255) not null,
    ctime     timestamp default CURRENT_TIMESTAMP not null,
    utime     timestamp default CURRENT_TIMESTAMP not null,
    sku_id    bigint not null,
    num       int not null,
    price     int not null comment '',
    uid       bigint not null,
    constraint order_pk2 unique (order_id)
);

create table t_sku
(
    id    bigint auto_increment primary key,
    name  varchar(10) not null,
    price int null comment '分为单位',
    ctime timestamp default CURRENT_TIMESTAMP not null,
    utime timestamp default CURRENT_TIMESTAMP not null,
    num   int null
);
create table t_user
(
    id    bigint auto_increment primary key,
    name  varchar(20) not null,
    ctime timestamp default CURRENT_TIMESTAMP not null,
    utime timestamp default CURRENT_TIMESTAMP not null
);

-- 监控告警配置表
create table t_deploy_info
(
    id              bigint auto_increment primary key,
    app             varchar(10) default '' not null,
    hosts           varchar(55) null comment '部署的主机ip,多个以逗号隔开',
    port            int null comment '服务端口',
    live_probe      varchar(255) null comment '监控检查接口路径',
    phone_webhook   varchar(255) null,
    dingding_webhook varchar(255) null comment '钉钉webhook',
    phone varchar(255) null comment '手机号'
);