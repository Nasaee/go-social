/*
 คือการ เพิ่ม foreign key constraint ให้กับตาราง posts
หรือพูดให้เข้าใจง่าย ๆ — คือ “บังคับให้ user_id ในตาราง posts ต้องอ้างถึง id ที่มีอยู่จริงในตาราง users เท่านั้น”
*/

ALTER TABLE
  posts
ADD
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);