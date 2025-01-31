-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Jan 31, 2025 at 07:50 AM
-- Server version: 8.0.34-cll-lve
-- PHP Version: 8.1.16

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

-- --------------------------------------------------------

--
-- Table structure for table `ada`
--

CREATE TABLE `ada` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `che`
--

CREATE TABLE `che` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `eng`
--

CREATE TABLE `eng` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `geo`
--

CREATE TABLE `geo` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `law`
--

CREATE TABLE `law` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `mcs`
--

CREATE TABLE `mcs` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `omo`
--

CREATE TABLE `omo` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `phy`
--

CREATE TABLE `phy` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `rav`
--

CREATE TABLE `rav` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `spo`
--

CREATE TABLE `spo` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `zam`
--

CREATE TABLE `zam` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

-- --------------------------------------------------------

--
-- Table structure for table `zis`
--

CREATE TABLE `zis` (
  `id` int NOT NULL,
  `course_id` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `name` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `weight` int DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `gender` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `teacher` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `faculty` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time1` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time2` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `time3` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time4` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time5` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci,
  `time_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `date_exam` text CHARACTER SET utf8mb3 COLLATE utf8mb3_persian_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_persian_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `ada`
--
ALTER TABLE `ada`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `che`
--
ALTER TABLE `che`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `eng`
--
ALTER TABLE `eng`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `geo`
--
ALTER TABLE `geo`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `law`
--
ALTER TABLE `law`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `mcs`
--
ALTER TABLE `mcs`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `omo`
--
ALTER TABLE `omo`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `phy`
--
ALTER TABLE `phy`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `rav`
--
ALTER TABLE `rav`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `spo`
--
ALTER TABLE `spo`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `zam`
--
ALTER TABLE `zam`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `zis`
--
ALTER TABLE `zis`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `ada`
--
ALTER TABLE `ada`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `che`
--
ALTER TABLE `che`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `eng`
--
ALTER TABLE `eng`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `geo`
--
ALTER TABLE `geo`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `law`
--
ALTER TABLE `law`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mcs`
--
ALTER TABLE `mcs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `omo`
--
ALTER TABLE `omo`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `phy`
--
ALTER TABLE `phy`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `rav`
--
ALTER TABLE `rav`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `spo`
--
ALTER TABLE `spo`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `zam`
--
ALTER TABLE `zam`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `zis`
--
ALTER TABLE `zis`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
