ALTER TABLE log_login
  ADD COLUMN request_id varchar(64) NOT NULL DEFAULT '' AFTER platform,
  ADD INDEX idx_request_id (request_id);

ALTER TABLE log_trace
  ADD COLUMN request_id varchar(64) NOT NULL DEFAULT '' AFTER reason,
  ADD INDEX rid (request_id);

ALTER TABLE log_msg
  ADD COLUMN request_id varchar(64) NOT NULL DEFAULT '' AFTER msg,
  ADD INDEX idx_request_id (request_id);
