CREATE TABLE `bank` (
  `id` bigint(10) PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `ifsc` varchar(15) NOT NULL UNIQUE,
  `branch` varchar(255) NOT NULL);