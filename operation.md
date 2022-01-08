# Usage

## DDL 语法

创建 range partition table, 并支持自动创建新分区。

```sql
CREATE TABLE t0
(
  name       VARCHAR(20),
  begin TIMESTAMP NOT NULL
)
PARTITION BY RANGE( UNIX_TIMESTAMP(begin) ) INTERVAL 1 year -- auto create partition when insert new data meet no suitable partition error.
(
  PARTITION p0 VALUES LESS THAN( UNIX_TIMESTAMP('2022-01-01 00:00:00') )
  -- ... auto create partition auto_px by TiDB.
);
```

设置搬移冷数据条件语句：

```sql
ALTER TABLE t0 MOVE PARTITIONS VALUES LESS THAN 2021 TO ENGINE AWS_S3;  -- 把 2021 年之前的数据搬到 s3
ALTER TABLE t0 MOVE PARTITIONS VALUES LESS THAN UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 YEAR)) TO ENGINE AWS_S3;  -- 总是把当前时间的前一年的数据搬到 s3 , 有一个 bug: https://github.com/pingcap/tidb/issues/6981, 所以这里的 expr 得加一个 cast
ALTER TABLE t0 MOVE PARTITIONS VALUES LESS THAN CAST(UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 YEAR)) as SIGNED) TO ENGINE AWS_S3;
ALTER TABLE t0 MOVE PARTITIONS VALUES LESS THAN 0 TO ENGINE AWS_S3;
```

## TiDB 侧 meta 管理

### Auto move partition table worker

1. check the table interval partition, find the table need to move the partition.
2. move the table partition.
   a. make this partition is only readable.
   b. move data.
   c. make moved partition raed by aws_s3.
   d. drop old partition.


### Auto create partition worker



### demo sql

```sql
CREATE TABLE trade_log (
  time       TIMESTAMP,
  id         bigint,
  info       VARCHAR(20)
) PARTITION BY RANGE( UNIX_TIMESTAMP(time) ) INTERVAL 1 month (
  PARTITION p2021_11 VALUES LESS THAN( UNIX_TIMESTAMP('2021-11-30 23:59:59') ),
  PARTITION p2021_12 VALUES LESS THAN( UNIX_TIMESTAMP('2021-12-31 23:59:59') )
);

INSERT INTO trade_log VALUES ('2021-11-10 10:00:00', 1, 'info: ...');
INSERT INTO trade_log VALUES ('2021-11-15 11:00:00', 2, 'info: ...');

INSERT INTO trade_log VALUES ('2021-12-01 12:00:00', 3, 'info: ...');
INSERT INTO trade_log VALUES ('2021-12-02 12:00:00', 4, 'info: ...');

INSERT INTO trade_log VALUES ('2022-01-01 12:00:00', 5, 'info: ...');

show create table trade_log\G


ALTER TABLE trade_log MOVE PARTITIONS VALUES LESS THAN CAST(UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 MONTH)) as SIGNED) TO ENGINE AWS_S3;

 select * from mysql.interval_partition_jobs;

  select * from mysql.gc_delete_range_done order by job_id desc limit 3;


ALTER TABLE trade_log DELETE PARTITIONS VALUES LESS THAN CAST(UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 MONTH)) as SIGNED);

```



```sql
-- tpch h1
drop table if exists LINEITEM;
 CREATE TABLE LINEITEM (
	L_ORDERKEY BIGINT NOT NULL, 
	L_PARTKEY BIGINT NOT NULL,
	L_SUPPKEY BIGINT NOT NULL, 
	L_LINENUMBER INTEGER,
	L_QUANTITY  DOUBLE,
	L_EXTENDEDPRICE DOUBLE,
	L_DISCOUNT  DOUBLE,
	L_TAX	 DOUBLE,
	L_RETURNFLAG CHAR(1),
	L_LINESTATUS CHAR(1),
	L_SHIPDATE DATE,
	L_COMMITDATE DATE,
	L_RECEIPTDATE DATE,
	L_SHIPINSTRUCT CHAR(25),
	L_SHIPMODE CHAR(10),
	L_COMMENT VARCHAR(44),
	PRIMARY KEY (L_ORDERKEY, L_LINENUMBER)
)PARTITION BY RANGE(L_ORDERKEY)
(
  PARTITION p0 values less than (10),
    PARTITION p1 values less than (100000000))
;

insert into LINEITEM values(1, 1, 1, 1, 1, 1, 1, 1, '1', '1', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(2, 2, 2, 2, 2, 2, 2, 2, '1', '1', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(3, 3, 3, 3, 3, 3, 3, 3, '1', '1', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(4, 4, 4, 4, 4, 4, 4, 4, '1', '1', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(5, 5, 5, 5, 5, 5, 5, 5, '1', '1', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(11, 11, 11, 11, 11, 11, 11, 11, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(12, 12, 12, 12, 12, 12, 12, 12, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(13, 13, 13, 13, 13, 13, 13, 13, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(14, 14, 14, 14, 14, 14, 14, 14, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(15, 15, 15, 15, 15, 15, 15, 15, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');

insert into LINEITEM values(100, 100, 100, 100, 100, 100, 100, 100, '1', '1', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(200, 200, 200, 200, 200, 200, 200, 200, '1', '1', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(211, 211, 211, 211, 211, 211, 211, 211, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(212, 212, 212, 212, 212, 212, 212, 212, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(213, 213, 213, 213, 213, 213, 213, 213, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(214, 214, 214, 214, 214, 214, 214, 214, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
insert into LINEITEM values(215, 215, 215, 215, 215, 215, 215, 215, '2', '2', '1998-08-01', '1998-08-01', '1998-08-01', '11111', '111', '111');
ALTER TABLE LINEITEM MOVE PARTITIONS VALUES LESS THAN CAST(100 as SIGNED)  TO ENGINE AWS_S3;

select  /*+ AGG_TO_COP() */
    l_returnflag, 
    l_linestatus, 
    sum(l_quantity) as sum_qty, 
    sum(l_extendedprice) as sum_base_price, 
    sum(l_extendedprice * (1 - l_discount)) as sum_disc_price, 
    sum(l_extendedprice * (1 - l_discount) * (1 + l_tax)) as sum_charge, 
    avg(l_quantity) as avg_qty, 
    avg(l_extendedprice) as avg_price, 
    avg(l_discount) as avg_disc, 
    count(*) as count_order 
from 
    lineitem
where 
    l_shipdate <= date'1998-12-01' - interval '90' day 
group by 
    l_returnflag, 
    l_linestatus
order by 
    l_returnflag, 
    l_linestatus;
```


```sql
CREATE TABLE t1
(
  id         bigint,
  name       VARCHAR(20),
  begin TIMESTAMP NOT NULL,
  a float,
  b double,
  c decimal(30,10),
  d year,
  e tinyint,
  f SMALLINT,
  g int
)
PARTITION BY RANGE( UNIX_TIMESTAMP(begin) ) INTERVAL 1 year 
(
  PARTITION p0 VALUES LESS THAN( UNIX_TIMESTAMP('2019-12-31 23:59:59') ),
  PARTITION p1 VALUES LESS THAN( UNIX_TIMESTAMP('2020-12-31 23:59:59') )
);

INSERT INTO t1 (id, name, begin) VALUES (1 ,'a', '2019-01-01 10:00:00');
INSERT INTO t1 (id, name, begin) VALUES (2 ,'b', '2019-01-01 11:00:00');
INSERT INTO t1 (id, name, begin) VALUES (3 ,'c', '2020-01-01 12:00:00');
INSERT INTO t1 (id, name, begin) VALUES (4 ,'e', '2020-01-01 12:00:00');
INSERT INTO t1 (id, name, begin) VALUES (5 ,'f', '2020-01-01 12:00:00');
INSERT INTO t1 (id, name, begin) VALUES (6 ,'f', '2021-01-01 12:00:00');
update t1 set a= id, b=id, c=id*1.1, d = year(begin),e = id%128, f=id%1024, g =id%65536;


ALTER TABLE t1 MOVE PARTITIONS VALUES LESS THAN CAST(UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 YEAR)) as SIGNED) TO ENGINE AWS_S3;
ALTER TABLE t1 DELETE PARTITIONS VALUES LESS THAN CAST(UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 YEAR)) as SIGNED);

select * from mysql.gc_delete_range_done;

```
