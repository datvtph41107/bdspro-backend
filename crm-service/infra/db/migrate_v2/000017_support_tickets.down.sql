-- +migrate Down
DROP TABLE IF EXISTS support_ticket_events;
DROP TABLE IF EXISTS support_ticket_notes;
DROP TABLE IF EXISTS support_tickets;
