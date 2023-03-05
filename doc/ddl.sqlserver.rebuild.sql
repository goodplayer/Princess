use
    princess
go

IF OBJECT_ID(N'dbo.princess_user', N'U') IS NOT NULL
    DROP TABLE dbo.princess_user
go

create table dbo.princess_user
(
    user_id     nvarchar(40)  not null
        constraint princess_user_pk
            primary key,
    user_name   nvarchar(40)  not null
        constraint princess_user_name_uniq
            unique,
    alias_name  nvarchar(40)  not null,
    password    nvarchar(200) not null,
    user_status int           not null,
    created_at  bigint        not null,
    updated_at  bigint,
    deleted_at  bigint
)
go

