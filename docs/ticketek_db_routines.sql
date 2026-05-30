-- MySQL dump 10.13  Distrib 8.0.43, for Win64 (x86_64)
--
-- Host: localhost    Database: ticketek_db
-- ------------------------------------------------------
-- Server version	9.4.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Temporary view structure for view `v_reporte_ventas_evento`
--

DROP TABLE IF EXISTS `v_reporte_ventas_evento`;
/*!50001 DROP VIEW IF EXISTS `v_reporte_ventas_evento`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_reporte_ventas_evento` AS SELECT 
 1 AS `evento_id`,
 1 AS `evento`,
 1 AS `fecha`,
 1 AS `estado_evento`,
 1 AS `total_entradas_vendidas`,
 1 AS `capacidad_total`,
 1 AS `porcentaje_ocupacion_total`,
 1 AS `ingresos_totales`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_ocupacion_sectores`
--

DROP TABLE IF EXISTS `v_ocupacion_sectores`;
/*!50001 DROP VIEW IF EXISTS `v_ocupacion_sectores`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_ocupacion_sectores` AS SELECT 
 1 AS `evento_id`,
 1 AS `evento`,
 1 AS `sector_id`,
 1 AS `sector`,
 1 AS `capacidad_maxima`,
 1 AS `capacidad_disponible`,
 1 AS `vendidas`,
 1 AS `porcentaje_ocupacion`,
 1 AS `precio`*/;
SET character_set_client = @saved_cs_client;

--
-- Final view structure for view `v_reporte_ventas_evento`
--

/*!50001 DROP VIEW IF EXISTS `v_reporte_ventas_evento`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `v_reporte_ventas_evento` AS select `e`.`id` AS `evento_id`,`e`.`titulo` AS `evento`,`e`.`fecha` AS `fecha`,`e`.`estado` AS `estado_evento`,count(`en`.`id`) AS `total_entradas_vendidas`,`e`.`capacidad_total` AS `capacidad_total`,round(((count(`en`.`id`) / `e`.`capacidad_total`) * 100),1) AS `porcentaje_ocupacion_total`,sum(`en`.`precio_pagado`) AS `ingresos_totales` from (`eventos` `e` left join `entradas` `en` on(((`en`.`evento_id` = `e`.`id`) and (`en`.`estado` = 'activa')))) group by `e`.`id`,`e`.`titulo`,`e`.`fecha`,`e`.`estado`,`e`.`capacidad_total` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_ocupacion_sectores`
--

/*!50001 DROP VIEW IF EXISTS `v_ocupacion_sectores`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `v_ocupacion_sectores` AS select `e`.`id` AS `evento_id`,`e`.`titulo` AS `evento`,`s`.`id` AS `sector_id`,`s`.`nombre` AS `sector`,`s`.`capacidad_maxima` AS `capacidad_maxima`,`s`.`capacidad_disponible` AS `capacidad_disponible`,(`s`.`capacidad_maxima` - `s`.`capacidad_disponible`) AS `vendidas`,round((((`s`.`capacidad_maxima` - `s`.`capacidad_disponible`) / `s`.`capacidad_maxima`) * 100),1) AS `porcentaje_ocupacion`,`s`.`precio` AS `precio` from (`sectores` `s` join `eventos` `e` on((`e`.`id` = `s`.`evento_id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-05-30 16:38:00
