#Queries to be fired

UPDATE tsmdb.channels SET `tenant_id`='7ca2664b-f379-43db-bdf9-7fdd40707219';
UPDATE tsmdb.discussions SET `tenant_id`='7ca2664b-f379-43db-bdf9-7fdd40707219';
UPDATE tsmdb.institutes SET `tenant_id`='7ca2664b-f379-43db-bdf9-7fdd40707219';
UPDATE tsmdb.likes SET `tenant_id`='7ca2664b-f379-43db-bdf9-7fdd40707219';
UPDATE tsmdb.notification_types SET `tenant_id`='7ca2664b-f379-43db-bdf9-7fdd40707219';
UPDATE tsmdb.notifications SET `tenant_id`='7ca2664b-f379-43db-bdf9-7fdd40707219';
UPDATE tsmdb.replies SET `tenant_id`='7ca2664b-f379-43db-bdf9-7fdd40707219';

# Find all related tables 
USE INFORMATION_SCHEMA;

SELECT *
FROM
  KEY_COLUMN_USAGE
WHERE
  REFERENCED_TABLE_NAME = 'programming_assignments'
  AND REFERENCED_COLUMN_NAME = 'id'
  AND TABLE_SCHEMA = 'tsmdb';

# NOTES BELOW
SELECT @i:=0;
UPDATE tsmdb.talent_academics SET `id` = @i:=@i+1;

#Innodb
SELECT  CONCAT('ALTER TABLE `', table_name, '` ENGINE=InnoDB;') AS sql_statements
FROM    information_schema.tables AS tb
WHERE   table_schema = 'tsmdb'
AND     `ENGINE` != 'InnoDB'
AND     `TABLE_TYPE` = 'BASE TABLE'
ORDER BY table_name DESC;

#Uncommon rows in two tables
SELECT * FROM T2
WHERE NOT EXISTS 
(SELECT * FROM T1 WHERE T1.id = T2.t1_id);

# Set of delete commans for uncommon rows
SELECT CONCAT("DELETE FROM tsmdb.talent_enquiries_technologies WHERE `enquiry_id` = '",T2.enquiry_id,"';")
as queries FROM tsmdb.talent_enquiries_technologies T2
WHERE NOT EXISTS 
(SELECT * FROM tsmdb.talent_enquiries T1 WHERE T1.ID = T2.enquiry_id)
GROUP BY T2.enquiry_id;



# Niranjan Local
DROP TABLE IF EXISTS tsam.likes;
DROP TABLE IF EXISTS tsam.notifications;
DROP TABLE IF EXISTS tsam.replies;
DROP TABLE IF EXISTS tsam.discussions;
DROP TABLE IF EXISTS tsam.channels;

SELECT COUNT(*) AS TOTAL, 
SUM(CASE WHEN `created_by`="cbdb77da-2a8b-43f7-bcb6-6a67388d951d" then 1 else 0 END) AS total_calling_count, 
SUM(CASE WHEN `updated_by`="cbdb77da-2a8b-43f7-bcb6-6a67388d951d" then 1 else 0 END) AS total_talent_count 
FROM tsmdb.talent_call_records;

. = any char except newline
\. = the actual dot character
.? = .{0,1} = match any char except newline zero or one times
.* = .{0,} = match any char except newline zero or more times
.+ = .{1,} = match any char except newline one or more times


WITH new_table AS
(
   SELECT talent_id,SUM(TIMESTAMPDIFF(month, from_date, IFNULL(to_date, CURDATE()))) AS 
  total_experience
   FROM talent_experiences
   group by talent_id
)
SELECT first_name,last_name FROM talents
 INNER JOIN new_table ON new_table.talent_id = id
 where new_table.total_experience = 93;


#  ================
 WITH RECURSIVE subordinate AS (
    SELECT  employee_credential_id,
            supervisor_credential_id,
            0 AS level
    FROM tsmdb.employee_supervisors
    WHERE employee_credential_id = 4529
 
    UNION ALL
 
    SELECT  e.employee_id, 
            e.first_name,
            e.last_name,
            e.manager_id,
            level + 1
    FROM employee e 
JOIN subordinate subo 
ON e.manager_id = subo.employee_id)
 
SELECT 
subo.employee_id,
    subo.first_name AS subordinate_first_name,
    subo.last_name AS subordinate_last_name,
    manager.employee_id AS direct_superior_id,
    manager.first_name AS direct_superior_first_name,
    manager.last_name AS direct_superior_last_name,
    subo.level
FROM subordinate subo 
JOIN employee manager 
ON subo.manager_id = manager.employee_id
ORDER BY level;

=======================================================
BTA :
SELECT * from (
    select
        id,
        `faculty_remarks`,
        `batch_topic_assignment_id`,
        `talent_id`,
        row_number() over 
        (partition by `talent_assignment_submissions`.`talent_id`,`talent_assignment_submissions`.`batch_topic_assignment_id`
        order by `submitted_on` desc) as rn
    from
        tsmdb.talent_assignment_submissions
) as t
where t.rn = 1