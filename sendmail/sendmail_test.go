package sendmail

import (
	"testing"
)

func Test_run(t *testing.T) {
	host := "smtp.mxhichina.com"
	port := 587
	user := "davis@xxx.com"
	passwd := "xxx"

	client := NewGMail(host, port, user, passwd)

	err := client.Mail(user, []string{"davis@qq.com"}, nil, "测试发送", "这事邮件内容<h1></h1>", ContentTypeHtml)
	if err != nil {
		t.Fatalf("发送邮件失败：%s", err.Error())
	}
}
