package main

const viewerPasswordTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Required - {{.SiteTitle}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #fafafa;
            color: #262626;
            line-height: 1.6;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            max-width: 400px;
            width: 100%;
            background-color: #ffffff;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            padding: 40px 32px;
        }
        h1 {
            font-size: 28px;
            font-weight: 600;
            color: #262626;
            margin-bottom: 8px;
            text-align: center;
        }
        .subtitle {
            font-size: 14px;
            color: #8e8e8e;
            text-align: center;
            margin-bottom: 32px;
        }
        .message {
            padding: 12px 16px;
            margin-bottom: 20px;
            border-radius: 8px;
            font-size: 14px;
            line-height: 18px;
        }
        .error {
            background-color: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            font-size: 14px;
            color: #262626;
        }
        input[type="password"] {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            font-size: 14px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #fafafa;
            transition: border-color 0.2s, background-color 0.2s;
        }
        input[type="password"]:hover {
            border-color: #a8a8a8;
        }
        input[type="password"]:focus {
            outline: none;
            border-color: #0095f6;
            background-color: #ffffff;
        }
        button {
            background-color: #000000;
            color: #ffffff;
            padding: 12px 24px;
            border: none;
            cursor: pointer;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 600;
            width: 100%;
            transition: transform 0.2s;
        }
        button:hover {
            transform: translateY(-1px);
        }
        @media (max-width: 768px) {
            body {
                padding: 0;
                align-items: flex-start;
            }
            .container {
                border: none;
                border-radius: 0;
                min-height: 100vh;
                padding: 40px 24px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.SiteTitle}}</h1>
        <div class="subtitle">This blog is password protected</div>

        {{if .Error}}
        <div class="message error">{{.Error}}</div>
        {{end}}

        <form method="POST" action="/viewer-auth">
            {{if .Redirect}}
            <input type="hidden" name="redirect" value="{{.Redirect}}">
            {{end}}
            <div class="form-group">
                <label for="password">Enter Password:</label>
                <input type="password" id="password" name="password" required autofocus>
            </div>

            <button type="submit">Access Blog</button>
        </form>
    </div>
</body>
</html>`
