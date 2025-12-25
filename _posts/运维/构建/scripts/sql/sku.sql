drop table if exists sku;
create table sku(
  id text not null,
  shop_id text not null,
  platform text not null,
  plat_sku_id text not null,
  plat_sku_name text default '',
  plat_sku_url text default '',
  plat_sku_img text default '',
  create_time timestamp with time zone,
  update_time timestamp with time zone,
  plat_sku_cid text default '',
  plat_spec jsonb,
  props jsonb,
  custom_props jsonb,
  system_props jsonb,
  relation_props jsonb,
  platform_extra text default '',
  shop_supply text default '',
  created text default '',
  modified text default '',
  tags text[],
  status bigint,
  cid int default 0,
  plat_cid text default '',
  presale int,
  presale_end_time jsonb,
  presale_delivery_time jsonb,
  extra jsonb,
  plat_extra text default '',
  goods_id text default '',
  quantity int,
  price text default '',
  outer_id text default '',
  PRIMARY KEY (id,shop_id)
)PARTITION BY HASH(shop_id);

-- 增加商品热度(最近30天咨询量)字段
alter table sku add column goods_hot int default 0;

-- 增加sku对应的spu_id 字段
alter table sku add column plat_spu_id text default '';

-- 函数索引所用到的一些函数

CREATE OR REPLACE FUNCTION prop_to_text(input_props JSONB) RETURNS TEXT AS $$
DECLARE
    prop JSONB;
    new_props TEXT := '';
    value_arr JSONB;
    element TEXT;
BEGIN
    FOR prop IN SELECT jsonb_array_elements(input_props)
    LOOP
        -- 获取 value 字段的所有元素
        value_arr := prop->'value';
        FOR i IN 0 .. jsonb_array_length(value_arr) - 1
        LOOP
            element := value_arr->>i;
            -- 如果 new_props 不为空，则添加分隔符
            IF new_props != '' THEN
                new_props := new_props || '$';
            END IF;
            -- 添加当前元素到 new_props
            new_props := new_props || element;
        END LOOP;
    END LOOP;
    RETURN lower(new_props);
END;
$$ LANGUAGE plpgsql IMMUTABLE;


-- 创建分区表的函数
CREATE OR REPLACE FUNCTION create_partitions(table_name TEXT, total_partitions INTEGER) RETURNS VOID AS $$
DECLARE
    partition_index INTEGER;
BEGIN
    FOR partition_index IN 0..total_partitions - 1 LOOP
        EXECUTE format('CREATE TABLE %1$s_p%2$s PARTITION OF %1$s FOR VALUES WITH (MODULUS %3$s, REMAINDER %2$s);', 
                       table_name, 
                       partition_index, 
                       total_partitions);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- 创建分区
SELECT create_partitions('sku', 8);

-- 创建插件
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE EXTENSION IF NOT EXISTS pg_bigm;


-- 创建索引
create unique index CONCURRENTLY   IF NOT EXISTS sku_shopid_plat_sku_id on sku using btree (shop_id,plat_sku_id);
create index CONCURRENTLY   IF NOT EXISTS shopid_platskucid on sku using btree (shop_id,plat_sku_cid);
create index CONCURRENTLY   IF NOT EXISTS platform_platskucid on sku using btree (platform,plat_sku_cid);
create index CONCURRENTLY   IF NOT EXISTS sku_shopid_createtime on sku using btree (shop_id,create_time desc);
create index CONCURRENTLY   IF NOT EXISTS sku_shopid_goodshot on sku using btree (shop_id,goods_hot desc);
create index CONCURRENTLY   IF NOT EXISTS sku_shopid_updatetime on sku using btree (shop_id,update_time desc);
create index CONCURRENTLY   IF NOT EXISTS sku_all on sku using gin(shop_id, status,platform,plat_sku_id,plat_sku_cid,create_time,update_time,plat_sku_name gin_bigm_ops, custom_props jsonb_path_ops, props jsonb_path_ops, system_props jsonb_path_ops, relation_props jsonb_path_ops);
create index CONCURRENTLY  IF NOT EXISTS sku_tags on sku using gin (shop_id, platform, tags array_ops) where  tags !='{}'::text[];
create index CONCURRENTLY   IF NOT EXISTS sku_shopid_plat_spu_id on sku using btree (shop_id,plat_spu_id);

CREATE INDEX CONCURRENTLY  IF NOT EXISTS sku_plat_sku_name_gin_bigm ON sku USING gin (
        shop_id,
        platform,
        lower(plat_sku_name) gin_bigm_ops
);

CREATE INDEX CONCURRENTLY  IF NOT EXISTS sku_props_gin ON sku USING gin (
        shop_id,
        platform,
        prop_to_text(props) gin_bigm_ops
);

CREATE INDEX CONCURRENTLY  IF NOT EXISTS sku_custom_props_gin ON sku USING gin (
        shop_id,
        platform,
        prop_to_text(custom_props) gin_bigm_ops
);

