CREATE TABLE `users` (
                         `id` varchar(36) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                         `avatar` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
                         `nickname` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
                         `username` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
                         `email` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
                         `password` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `status` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `sex` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `created_at` timestamp NULL DEFAULT NULL,
                         `updated_at` timestamp NULL DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `idx_username` (`username`),
                         UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;