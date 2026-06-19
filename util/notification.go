package util

// util.SendNotification(notification.NotificationMessage{
// 	CustomerNo: ctx.Param("customer_no"),
// 	Channel:    "dispatch",
// 	Category:   "critical",
// 	Title:      "Alert Over Speed",
// 	Content:    "HD001 Over Speed (70 km/h)",
// 	Event: "",
// })

// func SendNotification(notificationMessage notification.NotificationMessage) error {
// 	notificationMessage.Timestamp = time.Now().Format("2006-01-02 15:04:05")
// 	serverURL := config.Get().WebSocketServer

// 	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
// 	if err != nil {
// 		log.Println("WebSocket connection error:", err)
// 		return fmt.Errorf("failed to connect to WebSocket server: %w", err)
// 	}
// 	defer func() {
// 		if err := conn.Close(); err != nil {
// 			log.Println("Error closing WebSocket connection:", err)
// 		}
// 	}()

// 	notificationMessage.Channel = fmt.Sprintf("%s:%s:%s", notificationMessage.CustomerNo, notificationMessage.Site, notificationMessage.Channel)

// 	message, err := json.Marshal(notificationMessage)
// 	if err != nil {
// 		log.Println("Error marshalling notification message:", err)
// 		return fmt.Errorf("failed to marshal notification: %w", err)
// 	}

// 	err = conn.WriteMessage(websocket.TextMessage, message)
// 	if err != nil {
// 		log.Println("Error sending message:", err)
// 		return fmt.Errorf("failed to send notification: %w", err)
// 	}

// 	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
// 	_, response, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Println("Error reading acknowledgment response:", err)
// 		return fmt.Errorf("failed to read acknowledgment: %w", err)
// 	}

// 	log.Println("Response from server:", string(response))
// 	return nil
// }
