package handlers

import (
	"encoding/json"
	"log"
	"orderworker/database"
	"orderworker/models"
	"strconv"
)

func ProcessEvent(body []byte) {
	log.Printf("Event received: %s", body)

	var payload models.EventPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error unmarshalling event: %s", err)
		return
	}

	var order models.Order
	if err := database.PgDB.Preload("Client").Preload("Seller").Preload("Items.Product").First(&order, payload.OrderID).Error; err != nil {
		log.Printf("Error getting order: %s into Postgres: %s", payload.OrderID, err)
		return
	}

	if payload.EventType == "ORDER_CREATED" {
		itemsForScylla := make([]map[string]string, len(order.Items))
		for i, item := range order.Items {
			itemsForScylla[i] = map[string]string{
				"product_id":   item.ProductID.String(),
				"product_name": item.Product.Name,
				"quantity":     strconv.Itoa(item.Quantity),
				"price":        strconv.FormatInt(item.PriceAtPurchase, 10),
			}
		}
		if err := database.ScyllaSession.Query(`
			INSERT INTO order_details (order_id, order_number, order_date, client_id, client_name, seller_id, seller_name, total_amount, status, items)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			order.ID.String(), order.OrderNumber, order.CreatedAt, order.ClientID.String(), order.Client.Name, order.SellerID.String(), order.Seller.Name, order.TotalAmount, order.Status, itemsForScylla,
		).Exec(); err != nil {
			log.Printf("Erro ao inserir detalhes do pedido %s no ScyllaDB: %s", order.ID, err)
		} else {
			log.Printf("Pedido %s inserido com sucesso na tabela de leitura do ScyllaDB.", order.ID)
		}
	}
	if payload.EventType == "ORDER_STATUS_UPDATED" {
		if err := database.ScyllaSession.Query(`
			UPDATE order_details SET status = ? WHERE order_id = ?`,
			payload.NewStatus, payload.OrderID.String(),
		).Exec(); err != nil {
			log.Printf("Erro ao atualizar o estado do pedido %s no ScyllaDB: %s", payload.OrderID, err)
		} else {
			log.Printf("Estado do pedido %s atualizado para %s no ScyllaDB.", payload.OrderID, payload.NewStatus)
		}

		if err := database.ScyllaSession.Query(`
			INSERT INTO order_status_history (order_id, status, event_timestamp) VALUES (?, ?, ?)`,
			payload.OrderID.String(), payload.NewStatus, payload.Timestamp,
		).Exec(); err != nil {
			log.Printf("Erro ao inserir histórico de estado para o pedido %s no ScyllaDB: %s", payload.OrderID, err)
		} else {
			log.Printf("Histórico de estado para o pedido %s inserido no ScyllaDB.", payload.OrderID)
		}
	}
}
