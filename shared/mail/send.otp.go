package mail

import "fmt"

func GetSenderParamEmailRegist(toEmail, otpCode string) (string, string, []string) {
	subject := "OTP Verification for Amikom Pedia"

	// Use fmt.Sprintf to dynamically insert values into the content string
	content := fmt.Sprintf(`
        <h1>Hello %s,</h1>

        <p>We're excited to have you on board with Amikom Pedia! As part of our security measures, please use the following One-Time Password (OTP) to verify your account:</p>

        <h2>%s</h2>

        <p>This OTP is valid for a single use and will expire in 1 minute. Please do not share it with anyone for security reasons.</p>

        <p>If you did not attempt to create an account with Amikom Pedia, please disregard this email. Your account security is important to us.</p>

        <p>Thank you for choosing Amikom Pedia!</p>
    `, toEmail, otpCode)

	to := []string{toEmail}

	return subject, content, to

}

func GetSenderParamEmailForgotPass(toEmail, otpCode string) (string, string, []string) {
	subject := "Amikom Pedia: OTP for Forgot Password Request"

	// Use fmt.Sprintf to dynamically insert values into the content string
	content := fmt.Sprintf(`
        <h1>Hello %s,</h1>

        <p>We've received a request to reset the password for your Amikom Pedia account. To proceed, please use the following One-Time Password (OTP) to verify your identity:</p>

        <h2>%s</h2>

        <p>This OTP is valid for a single use and will expire in 1 minute. Please enter it on the password reset page to continue the process. If you didn't initiate this password reset, please contact us immediately.</p>

        <p>Thank you for choosing Amikom Pedia!</p>
    `, toEmail, otpCode)

	to := []string{toEmail}

	return subject, content, to
}
