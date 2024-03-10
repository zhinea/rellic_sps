CREATE TABLE `request_logs` (
                                `id` int PRIMARY KEY AUTO_INCREMENT,
                                `container_id` int,
                                `ip_addr` varchar(50),
                                `domain` varchar(100),
                                `type` varchar(20),
                                `created_at` timestamp
);

CREATE TABLE `containers` (
                              `id` int PRIMARY KEY AUTO_INCREMENT,
                              `config` text,
                              `is_active` tinyint(1) DEFAULT 1
);

CREATE TABLE `domains` (
                           `id` char(26) PRIMARY KEY,
                           `container_id` int,
                           `domain` varchar(100)
);

CREATE INDEX `request_logs_index_0` ON `request_logs` (`created_at`);
CREATE INDEX `request_logs_index_1` ON `request_logs` (`container_id`);

CREATE INDEX `containers_index_2` ON `containers` (`id`);

CREATE INDEX `domains_index_4` ON `domains` (`id`);

CREATE INDEX `domains_index_5` ON `domains` (`container_id`);