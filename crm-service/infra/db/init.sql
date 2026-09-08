CREATE UNIQUE INDEX uniq_sharings_contact_receiver 
ON sharing_access (contact_id, receiver_id, receiver_type);