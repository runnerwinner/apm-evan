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
