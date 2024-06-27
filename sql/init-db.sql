CREATE TABLE IF NOT EXISTS `request_logs` (
                                `id` int PRIMARY KEY AUTO_INCREMENT,
                                `container_id` int,
                                `ip_addr` varchar(50),
                                `domain` varchar(100),
                                `type` varchar(20),
                                `created_at` timestamp
);

CREATE TABLE IF NOT EXISTS `containers` (
                              `id` int PRIMARY KEY AUTO_INCREMENT,
                              `config` text,
                              `is_active` tinyint(1) DEFAULT 1
);

CREATE TABLE IF NOT EXISTS `domains` (
                           `id` char(26) PRIMARY KEY,
                           `container_id` int,
                           `domain` varchar(100)
);


-- CREATE INDEX `request_logs_index_0` ON `request_logs` (`created_at`);
set @x := (select count(*) from information_schema.statistics where table_name = 'request_logs' and index_name = 'request_logs_index_0' and table_schema = database());
set @sql := if( @x > 0, 'select ''Index exists.''', 'Alter Table request_logs ADD Index request_logs_index_0 (created_at);');
PREPARE stmt FROM @sql;
EXECUTE stmt;


-- CREATE INDEX `request_logs_index_1` ON `request_logs` (`container_id`);
set @x1 := (select count(*) from information_schema.statistics where table_name = 'request_logs' and index_name = 'request_logs_index_1' and table_schema = database());
set @sql1 := if( @x1 > 0, 'select ''Index exists.''', 'Alter Table request_logs ADD Index request_logs_index_1 (container_id);');
PREPARE stmt1 FROM @sql1;
EXECUTE stmt1;


-- CREATE INDEX `containers_index_2` ON `containers` (`id`);
set @x2 := (select count(*) from information_schema.statistics where table_name = 'containers' and index_name = 'containers_index_2' and table_schema = database());
set @sql2 := if( @x2 > 0, 'select ''Index exists.''', 'Alter Table containers ADD Index containers_index_2 (id);');
PREPARE stmt2 FROM @sql2;
EXECUTE stmt2;


-- CREATE INDEX `domains_index_4` ON `domains` (`id`);
set @x3 := (select count(*) from information_schema.statistics where table_name = 'domains' and index_name = 'domains_index_4' and table_schema = database());
set @sql3 := if( @x2 > 0, 'select ''Index exists.''', 'Alter Table domains ADD Index domains_index_4 (id);');
PREPARE stmt3 FROM @sql3;
EXECUTE stmt3;


-- CREATE INDEX `domains_index_5` ON `domains` (`container_id`);
set @x4 := (select count(*) from information_schema.statistics where table_name = 'domains' and index_name = 'domains_index_5' and table_schema = database());
set @sql4 := if( @x2 > 0, 'select ''Index exists.''', 'Alter Table domains ADD Index domains_index_5 (container_id);');
PREPARE stmt4 FROM @sql4;
EXECUTE stmt4;

