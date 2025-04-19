package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"math/big"
	"member_service/internal/database"
	"member_service/internal/email"
	"member_service/internal/validator"
	"os"

	mlog "github.com/mike504110403/goutils/log"

	basicauth "gitlab.com/gogogo2712128/common_moduals/auth/basicAuth"
)

// 驗證碼請求
func Validator(req ValidatorReq) error {
	mlog.Info("驗證碼請求")
	verificationCode, err := generateVerificationCode(6)
	if err != nil {
		return err
	}

	templateText := `<!DOCTYPE html>
	<html lang="en">
	<head>
    	<meta charset="UTF-8">
    	<title>Email Verification</title>
	</head>
	<body>
    	<p>亲爱的会员，</p>

    	<p>感谢您注册 Ur娱乐城！</p>

    	<p>为了确保您的账户安全并完成注册，我们需要验证您的电子邮件地址。请输入下方的验证码来完成验证：</p>

    	<p><strong>{{.VerificationCode}}</strong></p>

    	<p>如果您未曾注册 Ur娱乐城，请忽略此邮件或联系我们的客服团队。</p>

    	<p>感谢您的支持，期待您在 Ur娱乐城的愉快体验！</p>

	    <p>此致，<br>Ur娱乐城团队</p>

    	<p>---<br>（此邮件由系统自动发出，请勿回复。）</p>
		<img src="cid:image1" alt="Ur Logo" />
	</body>
	</html>
	`
	type EmailData struct {
		VerificationCode string
	}

	t, err := template.New("email").Parse(templateText)
	if err != nil {
		return err
	}

	// 填充模板數據
	data := EmailData{VerificationCode: verificationCode}
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	// 讀取要嵌入的圖檔
	imgData, err := os.ReadFile("static/img/logo.jpeg") // 修改成圖片的路徑
	if err != nil {
		log.Fatalf("無法讀取圖片: %v", err)
	}
	// 將圖檔轉換為 Base64 編碼的字符串
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	if err = email.SendEmail(email.EmailRequest{
		To:      req.ValidValue,
		Subject: "Ur娱乐城验证码",
		Body:    body.String(),
		Image:   imgBase64,
	}); err != nil {
		return err
	} else {
		validator.CacheVerificationCode(req.Username, req.ValidValue, verificationCode)
	}
	return nil
}

// 生成驗證碼
func generateVerificationCode(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}

	return string(code), nil
}

// 重設密碼
func ResetPassword(password string, mid int) error {
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}
	psHash, err := basicauth.HashPassword(password)
	if err != nil {
		return err
	}
	updateStr := `
		UPDATE Member
		SET password = ?
		WHERE id = ?
	`
	_, err = db.Exec(updateStr, psHash, mid)
	if err != nil {
		return err
	}
	return nil
}
