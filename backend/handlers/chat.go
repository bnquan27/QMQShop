package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bnquan27/QMQShop/backend/middleware"
)

type chatMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type chatRequest struct {
	Message string        `json:"message"`
	History []chatMessage `json:"history"`
}

type geminiSystemInstruction struct {
	Parts []geminiPart `json:"parts"`
}

type geminiReq struct {
	SystemInstruction *geminiSystemInstruction `json:"system_instruction,omitempty"`
	Contents          []geminiContent          `json:"contents"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResp struct {
	Candidates []geminiCandidate `json:"candidates"`
}

type geminiCandidate struct {
	Content      geminiContent `json:"content"`
	FinishReason string        `json:"finishReason"`
}

func GeminiChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := middleware.ParseJSON(r, &req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Dữ liệu không hợp lệ"})
		return
	}

	if req.Message == "" {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Vui lòng nhập tin nhắn"})
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "API key chưa được cấu hình"})
		return
	}

	// Limit history to last 10 exchanges to avoid token overrun
	maxHistory := 20
	if len(req.History) > maxHistory {
		req.History = req.History[len(req.History)-maxHistory:]
	}

	contents := []geminiContent{}
	for _, m := range req.History {
		if m.Role == "assistant" {
			contents = append(contents, geminiContent{Role: "model", Parts: []geminiPart{{Text: m.Text}}})
		} else {
			contents = append(contents, geminiContent{Role: "user", Parts: []geminiPart{{Text: m.Text}}})
		}
	}
	contents = append(contents, geminiContent{Role: "user", Parts: []geminiPart{{Text: req.Message}}})

	systemPrompt := `Bạn là trợ lý bán hàng thân thiện và đáng tin cậy của QMQ Shop (Quản lý Máy Tính Quận) — cửa hàng chuyên bán máy tính, linh kiện PC, laptop, màn hình và phụ kiện công nghệ tại Việt Nam.

NHIỆM VỤ CỦA BẠN:
- Tư vấn mua hàng, gợi ý sản phẩm phù hợp với nhu cầu khách hàng
- Giới thiệu các sản phẩm nổi bật, đang giảm giá hoặc khuyến mãi của shop
- Giải đáp thắc mắc về sản phẩm, thông số kỹ thuật cơ bản
- Giới thiệu về shop (uy tín, chất lượng, bảo hành, đổi trả)
- Hỗ trợ khách hàng chọn cấu hình PC phù hợp với ngân sách và nhu cầu
- Tận tình giải đáp mọi thắc mắc của khách hàng

PHONG CÁCH:
- Thân thiện, nhiệt tình, chu đáo như một nhân viên bán hàng thực thụ
- Trả lời bằng tiếng Việt
- Sử dụng ngôn ngữ dễ hiểu, tránh quá kỹ thuật nếu không cần thiết
- Luôn đặt lợi ích của khách hàng lên hàng đầu
- Khi khách hỏi về sản phẩm cụ thể, hãy gợi ý khách xem danh mục sản phẩm trên website shop
- Nếu khách cần tư vấn chi tiết về cấu hình, hãy hỏi rõ nhu cầu (làm việc văn phòng, gaming, render đồ họa, lập trình...) và ngân sách để tư vấn phù hợp`

	payload := geminiReq{
		SystemInstruction: &geminiSystemInstruction{
			Parts: []geminiPart{{Text: strings.TrimSpace(systemPrompt)}},
		},
		Contents: contents,
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-3-flash-preview:generateContent?key="+apiKey,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Lỗi kết nối đến Gemini"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Gemini API error (HTTP %d): %s", resp.StatusCode, string(respBody))
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Gemini API trả về lỗi"})
		return
	}

	var result geminiResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Lỗi đọc phản hồi từ Gemini"})
		return
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		middleware.JSON(w, http.StatusOK, map[string]string{"reply": "Xin lỗi, tôi không thể trả lời câu hỏi này."})
		return
	}

	reply := result.Candidates[0].Content.Parts[0].Text
	middleware.JSON(w, http.StatusOK, map[string]string{"reply": reply})
}
