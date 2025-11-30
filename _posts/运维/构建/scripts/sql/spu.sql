drop table if exists spu;
create table spu(
    id text not null,
    shop_id text not null,
    platform text not null,
    plat_goods_id text not null,
    plat_goods_name text default '',
    plat_goods_url text default '',
    plat_goods_img text default '',
    create_time timestamp with time zone,
    update_time timestamp with time zone,
    article_number text default '',
    sold_quantity int default 0,
    platform_extra text default '',
    shop_supply text default '',
    price text default '',
    props jsonb,
    sku_props jsonb,
    custom_props jsonb,
    system_props jsonb,
    relation_props jsonb,
    cid int default 0,
    plat_cid text default '',
    seller_cids text [],
    list_time text,
    created_time timestamp with time zone,
    pf_modified text,
    status bigint,
    is_onsale boolean,
    tags text [],
    profit jsonb,
    profit_stat jsonb,
    presale int,
    presale_end_time jsonb,
    presale_delivery_time jsonb,
    is_answer_goods boolean,
    load_from_page boolean,
    create_answer_goods_time timestamp with time zone,
    extra jsonb,
    plat_user_id text,
    plat_extra jsonb,
    PRIMARY KEY (id,shop_id)
)   PARTITION BY HASH(shop_id);
-- 增加商品热度字段
alter table spu
add column goods_hot int default 0;
-- 增加外部商品编码字段
alter table spu
add column outer_id text default '';

-- 函数索引所用到的一些函数
-- 处理属性   优化了sku_props的处理
CREATE OR REPLACE FUNCTION prop_to_text_opt(input_props JSONB, input_sku_props JSONB) 
RETURNS TEXT AS $$
DECLARE
    prop JSONB;
    sku_prop JSONB;
    new_props TEXT := '';
    value_arr JSONB;
    element TEXT;
    exists BOOLEAN;
BEGIN
    FOR prop IN SELECT jsonb_array_elements(input_props)
    LOOP
        -- 检查 prop 是否在 input_sku_props 中
        exists := FALSE;
        FOR sku_prop IN SELECT jsonb_array_elements(input_sku_props)
        LOOP
            IF sku_prop->>'prop_name' = prop->>'prop_name' THEN
                exists := TRUE;
                EXIT;
            END IF;
        END LOOP;
        
        IF NOT exists THEN
            IF jsonb_typeof(value_arr) = 'array' THEN
                -- 获取 value 字段的所有元素
                value_arr := prop->'value';
                FOR i IN 0 .. jsonb_array_length(value_arr) - 1
                LOOP
                    element := value_arr->>i;
                    -- 如果 new_props 不为空，则在添加新元素前添加分隔符
                    IF new_props != '' THEN
                        new_props := new_props || '$';
                    END IF;
                    -- 添加当前元素到 new_props
                    new_props := new_props || element;
                END LOOP;
            END IF;
        END IF;
    END LOOP;
    RETURN lower(new_props);
END;
$$ LANGUAGE plpgsql IMMUTABLE;



-- 属性处理 通用
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

SELECT create_partitions('spu', 8);



-- 创建插件
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE EXTENSION IF NOT EXISTS pg_bigm;


-- 索引
create unique index CONCURRENTLY  IF NOT EXISTS shopid_plat_goods_id on spu using btree (shop_id, plat_goods_id);
create index CONCURRENTLY  IF NOT EXISTS shopid_createtime on spu using btree (shop_id, create_time desc);
create index CONCURRENTLY  IF NOT EXISTS shopid_goodshot on spu using btree (shop_id, goods_hot desc);
create index CONCURRENTLY  IF NOT EXISTS shopid_platcid on spu using btree (shop_id, plat_cid);
create index CONCURRENTLY  IF NOT EXISTS platform_platcid on spu using btree (platform, plat_cid);
create index CONCURRENTLY  IF NOT EXISTS shopid_updatetime on spu using btree (shop_id, update_time desc);
create index CONCURRENTLY  IF NOT EXISTS shopid_createdtime on spu using btree (shop_id, created_time desc);
create index CONCURRENTLY  IF NOT EXISTS shopid_listtime on spu using btree (shop_id, list_time desc);
create index CONCURRENTLY  IF NOT EXISTS spu_all on spu using gin(
    shop_id,
    status,
    platform,
    status,
    plat_cid,
    plat_goods_id,
    create_time,
    update_time,
    custom_props jsonb_path_ops,
    props jsonb_path_ops,
    system_props jsonb_path_ops,
    relation_props jsonb_path_ops
);


create index CONCURRENTLY IF NOT EXISTS spu_tags on spu using gin (shop_id, platform, tags array_ops)
where tags != '{}'::text [];

CREATE INDEX CONCURRENTLY IF NOT EXISTS spu_props_gin ON spu USING gin (
        shop_id,
        platform,
        prop_to_text_opt(props, sku_props) gin_bigm_ops
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS spu_custom_props_gin ON spu USING gin (
        shop_id,
        platform,
        prop_to_text(custom_props) gin_bigm_ops
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS spu_plat_goods_name_gin_bigm ON spu USING gin (
        shop_id,
        platform,
        lower(plat_goods_name) gin_bigm_ops
);