package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type OrderVisaRequest struct {
	AccessKey int64 `json:"access_key"`
	OrderVisa int   `json:"visa"`
}

type Request struct {
	OrderId int64 `json:"order_id"`
}

func (h *Handler) issueOrderVisa(c *gin.Context) {
	var input Request
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusOK)

	go func() {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		sendOrderVisaRequest(input)
	}()
}

func sendOrderVisaRequest(request Request) {
	var visa = 0
	if rand.Intn(10)%10 >= 2 {
		visa = 1 + rand.Intn(2)
		fmt.Printf("Приказ №%d обработан\n", request.OrderId)
		fmt.Println("Статус визы: ", visa)
	} else {
		fmt.Printf("Приказ №%d не получилось обработать\n", request.OrderId)
	}

	answer := OrderVisaRequest{
		AccessKey: 123,
		OrderVisa: visa,
	}

	client := &http.Client{}

	jsonAnswer, _ := json.Marshal(answer)
	bodyReader := bytes.NewReader(jsonAnswer)

	requestURL := fmt.Sprintf("http://127.0.0.1:8000/api/orders/%d/update_visa/", request.OrderId)

	req, _ := http.NewRequest(http.MethodPut, requestURL, bodyReader)

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending PUT request:", err)
		return
	}

	defer response.Body.Close()
}
