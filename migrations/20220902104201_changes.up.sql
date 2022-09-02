-- create "pets" table
CREATE TABLE `pets` (`id` bigint NOT NULL AUTO_INCREMENT, `name` varchar(255) NOT NULL, `type` tinyint NOT NULL, `create_at` timestamp NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
