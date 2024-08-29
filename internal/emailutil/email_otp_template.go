package emailutil

import (
	"fmt"

	"github.com/dagota12/Loan-Tracker/bootstrap"
)

// OTPEmailTemplate generates an HTML email template with the provided OTP for user verification.
func OTPEmailTemplate(otp string, env *bootstrap.Env) string {
	return fmt.Sprintf(
		`<html>
            <head>
                    <style>
                            .container {
                                    max-width: 600px;
                                    margin: 0 auto;
                                    padding: 20px;
                                    font-family: Arial, sans-serif;
                            }

                            .header {
                                    text-align: center;
                                    margin-bottom: 20px;
                            }

                            .header h1 {
                                    font-size: 24px;
                                    font-weight: bold;
                            }

                            .content {
                                    margin-bottom: 20px;
                            }

                            .otp-code {
                                    font-size: 20px;
                                    font-weight: bold;
                                    text-align: center;
                                    margin-bottom: 10px;
                            }

                            .footer {
                                    text-align: center;
                                    font-size: 12px;
                                    color: #999;
                            }
                    </style>
            </head>
            <body>
                    <div class="container">
                            <div class="header">
                                    <h1>Reset Your Password</h1>
                            </div>
                            <div class="content">
                                    <p>To reset your password, please use the following OTP:</p>
                                    <div class="otp-code">%s</div>
                                    <p>This OTP is valid for a %v minutes. Please do not share it with anyone.</p>
                                    <p>Thank you!</p>
                            </div>
                            <div class="footer">
                                    <p>&copy; 2024 Your Company. All rights reserved.</p>
                            </div>
                    </div>
            </body>
            </html>`,
		otp,
		env.PassResetCodeExpirationMin,
	)
}
