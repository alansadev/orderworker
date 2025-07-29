# Table one: order_details_by_id

CREATE TABLE order_details_by_id (
order_id uuid PRIMARY KEY,
order_date timestamp,
order_status text,
total_amount bigint,
user_id uuid,
user_name text,
user_email text,
items list<frozen<map<text, text>>> // Lista de itens, com dados do produto duplicados
);

# Table two: order_status_history_by_order_id
CREATE TABLE order_status_history_by_order_id (
order_id uuid,
status_timestamp timestamp,
status text,
PRIMARY KEY (order_id, status_timestamp)
) WITH CLUSTERING ORDER BY (status_timestamp DESC); // Ordena do mais recente para o mais antigo