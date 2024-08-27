package emailutil

import "fmt"

func EmailTemplate(url string) string {
	return fmt.Sprintf(`
    <html>
    <head>
        <style>
            body {
                font-family: Arial, sans-serif;
                margin: 0;
                padding: 0;
            }

            .container {
                max-width: 600px;
                margin: 0 auto;
                padding: 20px;
                background-color: #f5f5f5;
            }

            .header {
                text-align: center;
                margin-bottom: 20px;
            }

            h1 {
                font-size: 24px;
                font-weight: bold;
            }

            .content {
                margin-bottom: 20px;
            }

            p {
                font-size: 16px;
                line-height: 1.5;
            }

            a {
                display: block;
                text-align: center;
                padding: 10px 20px;
                background-color: #007bff;
                color: #fff;
                text-decoration: none;
                border-radius: 5px;
            }

            .footer {
                text-align: center;
                font-size: 14px;
                color: #888;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Verify Your Account</h1>
            </div>
            <div class="content">
                <p>Please click the button below to verify your account:</p>
                <a href='%v'>Verify Email</a>
                <p>Thank you for registering!</p>
            </div>
            <div class="footer">
                <p>&copy; 2024 Your Company. All rights reserved.</p>
            </div>
        </div>
    </body>
    </html>`,
		url,
	)
}
