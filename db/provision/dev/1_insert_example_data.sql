INSERT INTO DOCUMENTS (NAME, VERSION, CREATED_AT, UPDATED_AT, DELETED_AT, DESCRIPTION, AUTHOR, STATUS) VALUES
  ('Java', '1.8', now(), NULL, NULL, 'JDK 1.8', 'Joel Whittaker-Smith', 'APPROVED'),
  ('Java', '1.8.1', now(), NULL, NULL, 'JDK 1.8.1', 'Joel Whittaker-Smith', 'APPROVED');


INSERT INTO SNIPPETS (CREATED_AT, UPDATED_AT, DELETED_AT, TEXT, TEST_CASE, DOCUMENT_NAME,DOCUMENT_VERSION) VALUES
  (now(), NULL, NULL, 'withJava(''1.8''){}', 'test', 'Java', '1.8' ),
  (now(), NULL, NULL, 'withJava(''1.8.1''){}', 'test', 'Java', '1.8.1');

